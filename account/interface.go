package account

// KeyPair interface
type KeyPair interface {
	Sign([]byte) []byte
	Verify([]byte, []byte) bool
	Encode() (string, string)
	Decode(string, string) error
}
