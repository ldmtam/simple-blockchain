package transaction

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/simpleblockchain/account"
	"github.com/simpleblockchain/common"
	"github.com/simpleblockchain/core/pb"
	"github.com/simpleblockchain/crypto/sha3"
	"golang.org/x/crypto/ed25519"
)

var (
	errTxInvalidArgument         = errors.New("invalid argument when creating tx")
	errInvalidProtoToTransaction = errors.New("protobuf message cannot be converted into Transaction")
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

// Marshal encodes tx using protobuf
func (tx *TxImpl) Marshal() (string, error) {
	txHash := tx.hash.CloneBytes()
	txValue := tx.value.Bytes()

	pbTx := &pb.Transaction{
		Hash:      txHash,
		From:      tx.from,
		To:        tx.to,
		Value:     txValue,
		Nonce:     tx.nonce,
		Timestamp: tx.timestamp,
	}

	return pbTx.String(), nil
}

// Unmarshal decode tx using protobuf
func (tx *TxImpl) Unmarshal(data string) error {
	pbTx := &pb.Transaction{}
	err := pbTx.XXX_Unmarshal([]byte(data))
	if err != nil {
		return errInvalidProtoToTransaction
	}
	tx.hash.SetBytes(pbTx.Hash)
	tx.from = pbTx.From
	tx.to = pbTx.To
	tx.value.SetBytes(pbTx.Value)
	tx.nonce = pbTx.Nonce
	tx.timestamp = pbTx.Timestamp
	return nil
}

func (tx *TxImpl) String() string {
	return fmt.Sprintf(`{"hash":"%s", "from":"%s", "to":"%s", "nonce":"%v", "value":"%s", "timestamp":"%v"}`,
		tx.hash.String(),
		tx.from,
		tx.to,
		tx.nonce,
		tx.value,
		tx.timestamp,
	)
}

// Sign signs the tx
func (tx *TxImpl) Sign(kp account.KeyPair) {
	sig := kp.Sign(tx.hash.CloneBytes())
	tx.signature = sig
}

// Verify verifies signature of tx
func (tx *TxImpl) Verify(pubKey []byte) bool {
	kp := &account.KeyPairImpl{
		PublicKey: ed25519.PublicKey(pubKey),
	}
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
