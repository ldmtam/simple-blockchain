package common

import (
	"crypto/sha256"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	arrBytes       = []byte{223, 32, 185, 230, 212, 229, 214, 213, 41, 94, 211, 253, 196, 69, 199, 193, 248, 193, 32, 182, 139, 154, 111, 14, 56, 167, 137, 236, 160, 182, 49, 173}
	hexString      = "df20b9e6d4e5d6d5295ed3fdc445c7c1f8c120b68b9a6f0e38a789eca0b631ad"
	originalString = "simpleblockchain"
)

func TestHash(t *testing.T) {
	hashStr := sha256.Sum256([]byte(originalString))

	isEqual := Equal(hashStr[:], arrBytes)
	assert.True(t, isEqual)

	var h Hash
	h.SetBytes(hashStr[:])
	assert.Equal(t, hexString, h.String())

	bytes := h.CloneBytes()
	assert.Equal(t, arrBytes, bytes)
}
