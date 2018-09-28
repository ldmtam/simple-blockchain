package abstraction

// Transaction interface
type Transaction interface {
	Sign(KeyPair)
	Verify([]byte) bool
	VerifyIntegrity() error
}
