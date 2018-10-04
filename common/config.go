package common

// P2PConfig is the config of p2p network.
type P2PConfig struct {
	ListenAddr string
	SeedNodes  []string
	Version    uint16
	ChainID    uint32
	DataPath   string
}
