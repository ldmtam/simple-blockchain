package account

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	message = []byte("hello, i am simple chain")
)

func TestKeyPair(t *testing.T) {
	// generate key
	k, err := NewKeyPair()
	assert.Nil(t, err)
	assert.NotNil(t, k)

	// uses private key to sign on message
	sig := k.Sign(message)
	fmt.Println(hex.EncodeToString(sig))
	assert.NotNil(t, sig)

	// uses public key to verify signature
	verified := k.Verify(sig, message)
	assert.True(t, verified)
}
