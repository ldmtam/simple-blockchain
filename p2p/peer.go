package p2p

import (
	"context"
	"encoding/binary"
	"errors"
	"io"
	"sync"
	"time"

	"github.com/willf/bloom"

	log "github.com/inconshreveable/log15"
	libnet "github.com/libp2p/go-libp2p-net"
	"github.com/libp2p/go-libp2p-peer"
	"github.com/multiformats/go-multiaddr"
)

// errors
var (
	ErrStreamCountExceed  = errors.New("stream count exceed")
	ErrMessageChannelFull = errors.New("message channel is full")
)

const (
	bloomMaxItemCount = 100000
	bloomErrRate      = 0.001

	msgChanSize    = 1024
	maxStreamCount = 4
)

// Peer represents a neighbor which we connect directly
//
// Peer's jobs are:
//  * managing streams which are responsible for sending and reading messages.
//  * recording messages we have sent and received so as to reduce redundant message in network.
//  * maintaining a priority queue of message to be sending.
type Peer struct {
	id          peer.ID
	addr        multiaddr.Multiaddr
	conn        libnet.Conn
	peerManager *PeerManager

	// stream is a channel where we get and send data through.
	streams     chan libnet.Stream
	streamCount int
	streamMutex sync.Mutex

	recentMsg      *bloom.BloomFilter
	bloomMutex     sync.Mutex
	bloomItemCount int

	urgentMsgCh chan *p2pMessage
	normalMsgCh chan *p2pMessage

	quitWriteCh chan struct{}
}

// NewPeer returns a new instance of Peer struct
func NewPeer(stream libnet.Stream, pm *PeerManager) *Peer {
	peer := &Peer{
		id:          stream.Conn().RemotePeer(),
		addr:        stream.Conn().RemoteMultiaddr(),
		conn:        stream.Conn(),
		peerManager: pm,
		streams:     make(chan libnet.Stream, maxStreamCount),
		recentMsg:   bloom.NewWithEstimates(bloomMaxItemCount, bloomErrRate),
		urgentMsgCh: make(chan *p2pMessage, msgChanSize),
		normalMsgCh: make(chan *p2pMessage, msgChanSize),
		quitWriteCh: make(chan struct{}),
	}
	peer.AddStream(stream)
	return peer
}

// Start  start peer's loop
func (p *Peer) Start() {
	log.Info("Peer is started.", "id", p.id.Pretty())

	go p.writeLoop()
}

// AddStream tries to add a Stream in stream pool.
func (p *Peer) AddStream(stream libnet.Stream) error {
	p.streamMutex.Lock()
	defer p.streamMutex.Unlock()

	if p.streamCount > maxStreamCount {
		return ErrStreamCountExceed
	}
	p.streams <- stream
	p.streamCount++
	go p.readLoop(stream)
	return nil
}

// CloseStream closes a stream and decrease the stream count.CloseStream
// only closes for writing, reading will works
func (p *Peer) CloseStream(stream libnet.Stream) {
	p.streamMutex.Lock()
	defer p.streamMutex.Unlock()

	stream.Close()
	p.streamCount--
}

func (p *Peer) newStream() (libnet.Stream, error) {
	p.streamMutex.Lock()
	defer p.streamMutex.Unlock()
	if p.streamCount > maxStreamCount {
		return nil, ErrStreamCountExceed
	}
	stream, err := p.peerManager.host.NewStream(context.Background(), p.id, protocolID)
	if err != nil {
		log.Error("Creating stream failed.", "pid", p.id.Pretty(), "err", err)
		return nil, err
	}
	p.streamCount++
	go p.readLoop(stream)
	return stream, nil
}

func (p *Peer) getStream() (libnet.Stream, error) {
	select {
	case stream := <-p.streams:
		return stream, nil
	default:
		stream, err := p.newStream()
		if err == ErrStreamCountExceed {
			break
		}
		return stream, err
	}
	return <-p.streams, nil
}

func (p *Peer) write(m *p2pMessage) error {
	stream, err := p.getStream()
	// if getStream fails, the TCP connection may be broken and we should stop the peer.
	if err != nil {
		log.Error("Get stream fails.", "err", err)
		return err
	}

	// 5kB/s
	deadline := time.Now().Add(time.Duration(len(m.content())/1024/5+1) * time.Second)
	if err = stream.SetWriteDeadline(deadline); err != nil {
		log.Warn("Write message failed.", "err", err)
		p.CloseStream(stream)
		return err
	}

	_, err = stream.Write(m.content())
	if err != nil {
		log.Warn("Write message failed.", "err", err)
		p.CloseStream(stream)
		return err
	}

	p.streams <- stream
	return nil
}

func (p *Peer) writeLoop() {
	for {
		select {
		case <-p.quitWriteCh:
			log.Info("Peer is stopped", "pid", p.id.Pretty(), "addr", p.addr)
			return
		case um := <-p.urgentMsgCh:
			go p.write(um)
		case nm := <-p.normalMsgCh:
			for done := false; !done; {
				select {
				case <-p.quitWriteCh:
					log.Info("Peer is stopped", "pid", p.id.Pretty(), "addr", p.addr)
					return
				case um := <-p.urgentMsgCh:
					go p.write(um)
				default:
					done = true
				}
			}
			go p.write(nm)
		}
	}
}

func (p *Peer) readLoop(stream libnet.Stream) {
	header := make([]byte, dataBegin)
	for {
		_, err := io.ReadFull(stream, header)
		if err != nil {
			log.Warn("Read header failed", "err", err)
			return
		}
		chainID := binary.BigEndian.Uint32(header[chainIDBegin:chainIDEnd])
		if chainID != p.peerManager.config.ChainID {
			log.Warn("Mismatched chainID.", "chainID", chainID)
			return
		}
		version := binary.BigEndian.Uint16(header[versionBegin:versionEnd])
		if version != p.peerManager.config.Version {
			log.Warn("Mismatched version.", "version", version)
			return
		}
		length := binary.BigEndian.Uint32(header[dataLengthBegin:dataLengthEnd])
		data := make([]byte, dataBegin+length)
		_, err = io.ReadFull(stream, data[dataBegin:])
		if err != nil {
			log.Warn("Read message failed", "err", err)
			return
		}
		copy(data[0:dataBegin], header)
		msg, err := parseP2PMessage(data)
		if err != nil {
			log.Error("Parse p2pmessage failed", "err", err)
			return
		}

		p.handleMessage(msg)
	}
}

// SendMessage puts message into corresponding channel.
func (p *Peer) SendMessage(msg *p2pMessage, mp MessagePriority, deduplicate bool) error {

	ch := p.urgentMsgCh
	if mp == NormalMessage {
		ch = p.normalMsgCh
	}
	select {
	case ch <- msg:
	default:
		log.Error("Sending message failed. Channel is full.", "messagePriority", mp)
		return ErrMessageChannelFull
	}
	return nil
}

func (p *Peer) handleMessage(msg *p2pMessage) error {
	p.peerManager.HandleMessage(msg, p.id)
	return nil
}
