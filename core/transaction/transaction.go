package transaction

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/simpleblockchain/account"
	"github.com/simpleblockchain/common"
	"github.com/simpleblockchain/crypto/sha3"
)

var (
	errTxInvalidArgument = errors.New("invalid argument when creating tx")
)

// TxImpl struct of a transaction
type TxImpl struct {
	hash      common.Hash
	from      []byte
	to        []byte
	value     *big.Int
	nonce     uint64
	timestamp int64

	signature []byte
}

// NewTransaction returns new transaction
func NewTransaction(from, to []byte, value *big.Int, nonce uint64, timestamp int64) (*TxImpl, error) {
	if from == nil || to == nil || value == nil {
		return nil, errTxInvalidArgument
	}

	txImpl := &TxImpl{
		from:      from,
		to:        to,
		value:     value,
		nonce:     nonce,
		timestamp: timestamp,
	}
	hash, err := txImpl.calcHash()
	if err != nil {
		return nil, err
	}
	txImpl.hash = hash
	return txImpl, nil
}

// From returns `from` address.
func (tx *TxImpl) From() []byte {
	return tx.from
}

// To returns `to` address.
func (tx *TxImpl) To() []byte {
	return tx.to
}

// Value returns tx value
func (tx *TxImpl) Value() *big.Int {
	return tx.value
}

// Nonce returns tx nonce
func (tx *TxImpl) Nonce() uint64 {
	return tx.nonce
}

// Timestamp returns time at which tx is created.
func (tx *TxImpl) Timestamp() int64 {
	return tx.timestamp
}

// Hash returns hash of transaction.
func (tx *TxImpl) Hash() common.Hash {
	return tx.hash
}

/**
 * For now, we will use json as encoding and decoding algorithm for our transaction type. Will
 * use more efficient algorithm as protobuf later.
 */

// Marshal encode tx
func (tx *TxImpl) Marshal() ([]byte, error) {
	bytes, err := json.Marshal(tx)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// Unmarshal decode tx
func (tx *TxImpl) Unmarshal(data []byte) (*TxImpl, error) {
	err := json.Unmarshal(data, &tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (tx *TxImpl) String() string {
	return fmt.Sprintf(`{"hash":"%s", "from":"%s", "to":"%s", "nonce":"%s", "value":"%s", "timestamp":"%s"}`,
		tx.hash.String(),
		tx.from,
		tx.to,
		string(tx.nonce),
		tx.value,
		string(tx.timestamp),
	)
}

// Sign signs the tx
func (tx *TxImpl) Sign(kp account.KeyPair) ([]byte, error) {
	return kp.Sign(tx.hash.CloneBytes()), nil
}

// Verify verifies signature of tx
func (tx *TxImpl) Verify(kp account.KeyPair) bool {
	return kp.Verify(tx.signature, tx.hash.CloneBytes())
}

// calcHash calculate hash of the transaction.
func (tx *TxImpl) calcHash() (common.Hash, error) {
	hasher := sha3.New256()

	value := tx.value.String()

	hasher.Write(tx.from)
	hasher.Write(tx.to)
	hasher.Write([]byte(value))
	hasher.Write(common.FromUint64(tx.nonce))
	hasher.Write(common.FromInt64(tx.timestamp))

	var h common.Hash
	h.SetBytes(hasher.Sum(nil))

	return h, nil
}
