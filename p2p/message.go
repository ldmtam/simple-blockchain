package p2p

import (
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"

	"github.com/golang/snappy"
)

/*
P2PMessage protocol:

 0               1               2               3              (bytes)
 0 1 2 3 4 5 6 7 0 1 2 3 4 5 6 7 0 1 2 3 4 5 6 7 0 1 2 3 4 5 6 7
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                         Chain ID                              |
+-------------------------------+-------------------------------+
|          Message Type         |            Version            |
+-------------------------------+-------------------------------+
|                         Data Length                           |
+---------------------------------------------------------------+
|                         Data Checksum							|
+---------------------------------------------------------------+
|                                                               |
.                             Data								.
|                                                               |
+---------------------------------------------------------------+
*/

// MessageType represents the message type.
type MessageType uint16

// consts
const (
	Ping MessageType = iota + 1
	Pong
	RoutingTableQuery
	RoutingTableResponse
	PublishTx
)

func (m MessageType) String() string {
	switch m {
	case Ping:
		return "Ping"
	case RoutingTableQuery:
		return "RoutingTableQuery"
	case RoutingTableResponse:
		return "RoutingTableResponse"
	case PublishTx:
		return "PublishTx"
	default:
		return fmt.Sprintf("unknown message type: %d \n", m)
	}
}

type p2pMessage []byte

const (
	chainIDBegin, chainIDEnd           = 0, 4
	messageTypeBegin, messageTypeEnd   = 4, 6
	versionBegin, versionEnd           = 6, 8
	dataLengthBegin, dataLengthEnd     = 8, 12
	dataChecksumBegin, dataChecksumEnd = 12, 16
	dataBegin                          = 16
)

var (
	errInvalidChecksum   = errors.New("invalid data checksum")
	errUnmatchDataLength = errors.New("unmatch data length")
	errMessageTooShort   = errors.New("message too short")
)

func (m *p2pMessage) content() []byte {
	return []byte(*m)
}

func (m *p2pMessage) chainID() uint32 {
	return binary.BigEndian.Uint32(m.content()[chainIDBegin:chainIDEnd])
}

func (m *p2pMessage) messageType() MessageType {
	return MessageType(binary.BigEndian.Uint16(m.content()[messageTypeBegin:messageTypeEnd]))
}

func (m *p2pMessage) version() uint16 {
	return binary.BigEndian.Uint16(m.content()[versionBegin:versionEnd])
}

func (m *p2pMessage) dataLength() uint32 {
	return binary.BigEndian.Uint32(m.content()[dataLengthBegin:dataLengthEnd])
}

func (m *p2pMessage) checksum() uint32 {
	return binary.BigEndian.Uint32(m.content()[dataChecksumBegin:dataChecksumEnd])
}

func (m *p2pMessage) rawData() []byte {
	return m.content()[dataBegin:]
}

func (m *p2pMessage) data() ([]byte, error) {
	// uncompress data
	data, err := snappy.Decode(nil, m.rawData())
	if err != nil {
		return nil, err
	}

	return data, nil
}

func newP2PMessage(chainID uint32, messageType MessageType, version uint16, data []byte) *p2pMessage {
	var compressedData []byte
	if len(data) > 0 {
		compressedData = snappy.Encode(nil, data)
	}

	m := make([]byte, dataBegin+len(compressedData))

	binary.BigEndian.PutUint32(m[chainIDBegin:chainIDEnd], chainID)
	binary.BigEndian.PutUint16(m[messageTypeBegin:messageTypeEnd], uint16(messageType))
	binary.BigEndian.PutUint16(m[versionBegin:versionEnd], version)
	binary.BigEndian.PutUint32(m[dataLengthBegin:dataLengthEnd], uint32(len(compressedData)))
	binary.BigEndian.PutUint32(m[dataChecksumBegin:dataChecksumEnd], crc32.ChecksumIEEE(compressedData))
	copy(m[dataBegin:], compressedData)

	var message p2pMessage = m
	return &message
}

func parseP2PMessage(data []byte) (*p2pMessage, error) {
	// check data length
	if len(data) < dataBegin {
		return nil, errMessageTooShort
	}

	m := p2pMessage(data)

	// check message checksum
	if crc32.ChecksumIEEE(m.rawData()) != m.checksum() {
		return nil, errInvalidChecksum
	}

	// check message length
	if len(m.rawData()) != int(m.dataLength()) {
		return nil, errUnmatchDataLength
	}

	return &m, nil
}

// IncomingMessage is the struct of message sent via the stream
type IncomingMessage struct {
	from PeerID
	data []byte
	typ  MessageType
}

// From returns the peerID who sends the message.
func (m *IncomingMessage) From() PeerID {
	return m.from
}

// Data returns the bytes.
func (m *IncomingMessage) Data() []byte {
	return m.data
}

// Type returns the message type.
func (m *IncomingMessage) Type() MessageType {
	return m.typ
}
