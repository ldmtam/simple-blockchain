package common

import (
	"bytes"
	"encoding/hex"
)

const (
	// HashLength length of hash
	HashLength = 32

	// AddressLength of address
	AddressLength = 32
)

// Hash is hash of data
type Hash [HashLength]byte

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

// Address ..
type Address [AddressLength]byte

// CloneBytes ...
func (a Address) CloneBytes() []byte {
	b := make([]byte, AddressLength)
	copy(b, a[:])
	return b
}

// SetBytes ...
func (a *Address) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-AddressLength:]
	}
	copy(a[AddressLength-len(b):], b)
}

func (a Address) String() string {
	return hex.EncodeToString(a[:])
}

// Equals ...
func (a Address) Equals(b Address) bool {
	a1 := a.CloneBytes()
	a2 := b.CloneBytes()
	return Equal(a1, a2)
}
