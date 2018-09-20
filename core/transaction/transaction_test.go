package transaction

import (
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	"github.com/simpleblockchain/account"
	"github.com/stretchr/testify/assert"
)

const (
	fromPrivKey = "272e1da3327afa205dffe5acf8794b1ce79278b3ef44245b49754bdf0bbe8a77e05f3e24f866e33929120458507ee42c3cc2b3ab86dffe147ed0cbb9ba3d06c2"
	fromPubKey  = "e05f3e24f866e33929120458507ee42c3cc2b3ab86dffe147ed0cbb9ba3d06c2"

	toPrivKey = "c711d7757a2bb7027488c8f3db6b52d352b6b1d590b3e50c08034858166ae86d2af9c075359c199ec85c69d5f737149e4048c6c7b3f69a8ed917903192c73a1d"
	toPubKey  = "2af9c075359c199ec85c69d5f737149e4048c6c7b3f69a8ed917903192c73a1d"
)

func TestCreateTransaction(t *testing.T) {
	from, err := hex.DecodeString(fromPubKey)
	assert.Nil(t, err)
	assert.NotNil(t, from)

	to, err := hex.DecodeString(toPubKey)
	assert.Nil(t, err)
	assert.NotNil(t, to)

	value := big.NewInt(int64(20))
	assert.NotNil(t, value)

	nonce := uint64(1)
	timestamp := time.Now().Unix()

	tx, _ := NewTransaction(from, to, value, nonce, timestamp)
	assert.Nil(t, err)
	assert.NotNil(t, tx)
}

func createTx() *TxImpl {
	from, _ := hex.DecodeString(fromPubKey)

	to, _ := hex.DecodeString(toPubKey)

	value := big.NewInt(int64(20))

	nonce := uint64(1)

	timestamp := time.Now().Unix()

	tx, _ := NewTransaction(from, to, value, nonce, timestamp)
	return tx
}

func TestGetFrom(t *testing.T) {
	tx := createTx()

	assert.Equal(t, tx.from, tx.From())
}

func TestGetTo(t *testing.T) {
	tx := createTx()

	assert.Equal(t, tx.to, tx.To())
}

func TestGetValue(t *testing.T) {
	tx := createTx()

	assert.Equal(t, tx.value, tx.Value())
}

func TestGetNonce(t *testing.T) {
	tx := createTx()

	assert.Equal(t, tx.nonce, tx.Nonce())
}

func TestGetTimestamp(t *testing.T) {
	tx := createTx()

	assert.Equal(t, tx.timestamp, tx.Timestamp())
}

func TestGetHash(t *testing.T) {
	tx := createTx()

	assert.Equal(t, tx.hash, tx.Hash())
}

func TestSignTx(t *testing.T) {
	tx := createTx()

	fromKp := &account.KeyPairImpl{}
	err := fromKp.Decode(fromPrivKey, fromPubKey)
	assert.Nil(t, err)

	tx.Sign(fromKp)
	assert.Nil(t, err)
	assert.NotNil(t, tx.signature)

	toKp := &account.KeyPairImpl{}
	err = toKp.Decode(toPrivKey, toPubKey)
	assert.Nil(t, err)

}
