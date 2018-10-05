package abstraction

// P2PService interface of p2p service.
type P2PService interface {
	Start() error
	Stop()
}
