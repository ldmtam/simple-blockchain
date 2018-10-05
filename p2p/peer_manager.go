package p2p

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/multiformats/go-multiaddr"

	log "github.com/inconshreveable/log15"
	"github.com/ldmtam/tam-chain/common"

	"github.com/libp2p/go-libp2p-host"
	kbucket "github.com/libp2p/go-libp2p-kbucket"
	peer "github.com/libp2p/go-libp2p-peer"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	"github.com/uber-go/atomic"
)

var (
	dumpRoutingTableInterval = 2 * time.Minute
	syncRoutingTableInterval = 30 * time.Second
)

const (
	maxNeighborCount = 32
	bucketSize       = 20

	routingTableFile = "routing.table"
)

// PeerManager manages all peers we connect
//
// keep list of Peers
// discovering other peer and maintain routing table.
type PeerManager struct {
	neighbors     *sync.Map // map[peer.ID]*Peer
	neighborCount int

	subs   *sync.Map // map[MessageType]map[string]chan IncomingMessage
	quitCh chan struct{}

	host           host.Host
	config         *common.P2PConfig
	routingTable   *kbucket.RoutingTable
	peerStore      peerstore.Peerstore
	lastUpdateTime atomic.Int64

	wg *sync.WaitGroup
}

// NewPeerManager returns a new instance of PeerManager struct.
func NewPeerManager(host host.Host, config *common.P2PConfig) *PeerManager {
	routingTable := kbucket.NewRoutingTable(bucketSize, kbucket.ConvertPeerID(host.ID()), time.Second, host.Peerstore())
	return &PeerManager{
		neighbors:    new(sync.Map),
		subs:         new(sync.Map),
		quitCh:       make(chan struct{}),
		routingTable: routingTable,
		host:         host,
		config:       config,
		peerStore:    host.Peerstore(),
		wg:           new(sync.WaitGroup),
	}
}

// Start starts peer manager's jobs.
func (pm *PeerManager) Start() {
	pm.parseSeeds()
	pm.loadRoutingTable()

	go pm.dumpRoutingTableLoop()
}

// Stop stops peer manager's jobs.
func (pm *PeerManager) Stop() {
	close(pm.quitCh)
	pm.wg.Wait()
}

func (pm *PeerManager) parseSeeds() {
	for _, seed := range pm.config.SeedNodes {
		peerID, addr, err := parseMultiaddr(seed)
		if err != nil {
			log.Error("Parse seed nodes error", "err", err)
			continue
		}
		pm.storePeer(peerID, []multiaddr.Multiaddr{addr})
	}
}

func (pm *PeerManager) loadRoutingTable() {
	file, err := os.Open(filepath.Join(pm.config.DataPath, routingTableFile))
	if err != nil {
		log.Error("Reading routing table file failed.", "err", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		pid, addr, err := parseMultiaddr(line)
		if err != nil {
			log.Error("Parsing multiaddress failed.", "err", err)
			continue
		}
		pm.storePeer(pid, []multiaddr.Multiaddr{addr})
	}
}

func (pm *PeerManager) dumpRoutingTableLoop() {
	pm.wg.Add(1)
	var lastSaveTime int64

	for {
		select {
		case <-pm.quitCh:
			pm.wg.Done()
			return
		case <-time.After(dumpRoutingTableInterval):
			if lastSaveTime < pm.lastUpdateTime.Load() {
				pm.dumpRoutingTable()
				lastSaveTime = time.Now().Unix()
			}
		}
	}
}

func (pm *PeerManager) dumpRoutingTable() {
	file, err := os.Create(filepath.Join(pm.config.DataPath, routingTableFile))
	if err != nil {
		log.Error("Creating routing table file failed.", "err", err)
		return
	}
	defer file.Close()
	file.WriteString(fmt.Sprintf("# %s\n", time.Now().String()))
	for _, pid := range pm.routingTable.ListPeers() {
		for _, addr := range pm.peerStore.Addrs(pid) {
			line := fmt.Sprintf("%s/ipfs/%s\n", addr.String(), pid.Pretty())
			file.WriteString(line)
		}
	}
}

func (pm *PeerManager) storePeer(peerID peer.ID, addrs []multiaddr.Multiaddr) {
	pm.peerStore.AddAddrs(peerID, addrs, peerstore.PermanentAddrTTL)
	pm.routingTable.Update(peerID)
	pm.lastUpdateTime.Store(time.Now().Unix())
}
