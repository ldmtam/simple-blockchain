package common

import (
	"bytes"
	"encoding/hex"
)

const (
	// HashLength length of hash
	HashLength = 32
)

// Hash is hash of data
type Hash [HashLength]byte

// FromHex convert string to bytes
func FromHex(s string) []byte {
	if len(s) > 1 {
		// if hex string contains "0x" or "0X", trim those characters
		if s[0:2] == "0x" || s[0:2] == "0X" {
			s = s[2:]
		}
	}
	// prepend 0 to odd hex string to make it even
	if len(s)%2 == 1 {
		s = "0" + s
	}

	h, _ := hex.DecodeString(s)
	return h
}

// CloneBytes returns a copy of the bytes which represent the hash as a byte
func (hash *Hash) CloneBytes() []byte {
	newHash := make([]byte, HashLength)
	copy(newHash, hash[:])

	return newHash
}

// SetBytes sets the byte which represent the hash.
func (hash *Hash) SetBytes(newHash []byte) {
	if len(newHash) != HashLength {
		return
	}
	copy(hash[:], newHash)
}

// String returns hash a as hex string
func (hash *Hash) String() string {
	return hex.EncodeToString(hash[:])
}

// Equals compare two Hash. True is equal, otherwise false/
func (hash *Hash) Equals(b *Hash) bool {
	h1 := hash.CloneBytes()
	h2 := b.CloneBytes()
	return bytes.Compare(h1, h2) == 0
}

// Equal checks whether byte slice a and b are equal.
func Equal(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
