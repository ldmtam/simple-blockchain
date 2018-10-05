package p2p

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/libp2p/go-libp2p"

	"github.com/ldmtam/tam-chain/common"

	crypto "github.com/libp2p/go-libp2p-crypto"
	"github.com/libp2p/go-libp2p-host"

	log "github.com/inconshreveable/log15"
	"github.com/libp2p/go-libp2p-peer"
)

// PeerID is the alias of peer.ID
type PeerID = peer.ID

const (
	protocolID  = "tamchain/1.0"
	privKeyFile = "priv.key"
)

// errors
var (
	ErrPortUnavailable = errors.New("port is unavailable")
)

// NetService is the implementation of Service interface
type NetService struct {
	host        host.Host
	peerManager *PeerManager
	config      *common.P2PConfig
}

// NewNetService returns a NetService instance
func NewNetService(config *common.P2PConfig) (*NetService, error) {
	ns := &NetService{
		config: config,
	}

	if err := os.MkdirAll(config.DataPath, 0766); config.DataPath != "" && err != nil {
		log.Error("failed to create p2p datapath", "err", err, "path", config.DataPath)
		return nil, err
	}

	privKey, err := getOrCreateKey(filepath.Join(config.DataPath, privKeyFile))
	if err != nil {
		log.Error("failed to get private key.", "err", err, "path", config.DataPath)
		return nil, err
	}

	ipv4Address, err := ExternalIPv4()
	if err != nil {
		log.Error("failed to get external IP.", "err", err)
		return nil, err
	}

	listenAddress := fmt.Sprintf("%s:%s", ipv4Address, config.Port)
	host, err := ns.createHost(privKey, listenAddress)
	if err != nil {
		log.Error("failed to create a host.", "err", err, "listenAddr", listenAddress)
		return nil, err
	}
	ns.host = host

	ns.peerManager = NewPeerManager(host, config)

	return ns, nil
}

func (ns *NetService) createHost(pk crypto.PrivKey, listenAddr string) (host.Host, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		return nil, err
	}

	if !isPortAvailable(tcpAddr.Port) {
		return nil, ErrPortUnavailable
	}

	opts := []libp2p.Option{
		libp2p.Identity(pk),
		libp2p.NATPortMap(),
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/%s/tcp/%d", tcpAddr.IP, tcpAddr.Port)),
	}

	h, err := libp2p.New(context.Background(), opts...)
	if err != nil {
		return nil, err
	}
	h.SetStreamHandler(protocolID, nil)
	return h, nil
}

// Start starts the job.
func (ns *NetService) Start() error {
	go ns.peerManager.Start()
	log.Info("Net service started")
	for _, addr := range ns.host.Addrs() {
		log.Info("Listening on address", "address", fmt.Sprintf("%s/ipfs/%s", addr, ns.host.ID().Pretty()))
	}
	return nil
}

// Stop stops the job.
func (ns *NetService) Stop() {
	log.Info("Net service stopped")
}
