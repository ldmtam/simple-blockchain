package account

import (
	"crypto/rand"
	"encoding/hex"
	"errors"

	"golang.org/x/crypto/ed25519"
)

var (
	errGenerateKeyFailed = errors.New("error: generate key failed")
	errSignFailed        = errors.New("error: sign failed")
	errDecodeFailed      = errors.New("error: decode failed")
)

// KeyPairImpl ...
type KeyPairImpl struct {
	PrivateKey ed25519.PrivateKey
	PublicKey  ed25519.PublicKey
}

// NewKeyPair returns new ed25519 key pair
func NewKeyPair() (*KeyPairImpl, error) {
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, errGenerateKeyFailed
	}
	return &KeyPairImpl{
		PrivateKey: privKey,
		PublicKey:  pubKey,
	}, nil
}

// Sign use private key to sign on data
func (kp *KeyPairImpl) Sign(message []byte) []byte {
	return ed25519.Sign(kp.PrivateKey, message)
}

// Verify uses public key to verify a signature
func (kp *KeyPairImpl) Verify(sig, message []byte) bool {
	return ed25519.Verify(kp.PublicKey, message, sig)
}

// Encode encodes keypair from []byte to hex string
func (kp *KeyPairImpl) Encode() (string, string) {
	return hex.EncodeToString(kp.PrivateKey), hex.EncodeToString(kp.PublicKey)
}

// Decode decodes keypair from hex string to []byte
func (kp *KeyPairImpl) Decode(privKey, pubKey string) error {
	privKeyBytes, err := hex.DecodeString(privKey)
	if err != nil {
		return errDecodeFailed
	}

	pubKeyBytes, err := hex.DecodeString(pubKey)
	if err != nil {
		return errDecodeFailed
	}

	kp.PrivateKey = ed25519.PrivateKey(privKeyBytes)
	kp.PublicKey = ed25519.PublicKey(pubKeyBytes)

	return nil
}
