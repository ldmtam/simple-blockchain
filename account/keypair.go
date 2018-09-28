package account

import (
	"crypto/rand"
	"encoding/hex"
	"errors"

	"golang.org/x/crypto/ed25519"
)

var (
	errGenerateKeyFailed = errors.New("cannot generate key pair")
	errSignFailed        = errors.New("cannot sign message")

	errInvalidPrivateKeyHexString = errors.New("invalid private key hex string")
	errInvalidPublicKeyHexString  = errors.New("invalid public key hex string")
	errInvalidPrivateKeyLength    = errors.New("invalid private key length")
	errInvalidPublicKeyLength     = errors.New("invalid public key length")

	errInvalidPublicKey  = errors.New("invalid public key format")
	errInvalidPrivateKey = errors.New("invalid private key format")
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

// EncodePrivateKey encode private key to string
func (kp *KeyPairImpl) EncodePrivateKey() string {
	return hex.EncodeToString(kp.PrivateKey)
}

// EncodePublicKey encode public key to string
func (kp *KeyPairImpl) EncodePublicKey() string {
	return hex.EncodeToString(kp.PublicKey)
}

// DecodePrivateKey decode private key hex string
func (kp *KeyPairImpl) DecodePrivateKey(privKey string) error {
	privKeyBytes, err := hex.DecodeString(privKey)
	if err != nil {
		return errInvalidPrivateKeyHexString
	}

	if len(privKeyBytes) != ed25519.PrivateKeySize {
		return errInvalidPrivateKeyLength
	}

	kp.PrivateKey = ed25519.PrivateKey(privKeyBytes)

	return nil
}

// DecodePublicKey decode public key hex string
func (kp *KeyPairImpl) DecodePublicKey(pubKey string) error {
	pubKeyBytes, err := hex.DecodeString(pubKey)
	if err != nil {
		return errInvalidPublicKeyHexString
	}

	if len(pubKeyBytes) != ed25519.PublicKeySize {
		return errInvalidPublicKeyLength
	}

	kp.PublicKey = ed25519.PublicKey(pubKeyBytes)

	return nil
}
