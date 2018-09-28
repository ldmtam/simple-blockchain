package abstraction

// KeyPair interface
type KeyPair interface {
	Sign([]byte) []byte
	Verify([]byte, []byte) bool
	EncodePrivateKey() string
	EncodePublicKey() string
	DecodePrivateKey(string) error
	DecodePublicKey(string) error
}
