package transaction

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/gogo/protobuf/proto"
	"github.com/ldmtam/tam-chain/abstraction"
	"github.com/ldmtam/tam-chain/account"
	"github.com/ldmtam/tam-chain/common"
	"github.com/ldmtam/tam-chain/crypto/sha3"
	"github.com/ldmtam/tam-chain/proto"
	"golang.org/x/crypto/ed25519"
)

var (
	errTxInvalidArgument           = errors.New("invalid argument when creating tx")
	errInvalidProtoToTransaction   = errors.New("protobuf message cannot be converted into Transaction")
	errInvalidTransactionToProto   = errors.New("transaction cannot be converted to protobuf message")
	errInvalidTransacionHash       = errors.New("invalid transaction hash")
	errInvalidTransactionSignature = errors.New("invalid transaction signature")
)

// TxImpl struct of a transaction
type TxImpl struct {
	hash      common.Hash
	from      common.Address
	to        common.Address
	value     *big.Int
	nonce     uint64
	timestamp int64

	signature []byte
}

// NewTransaction returns new transaction
func NewTransaction(from, to common.Address, value *big.Int, nonce uint64, timestamp int64) (*TxImpl, error) {
	if from.CloneBytes() == nil || to.CloneBytes() == nil || value == nil {
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
	return tx.from.CloneBytes()
}

// To returns `to` address.
func (tx *TxImpl) To() []byte {
	return tx.to.CloneBytes()
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

// Signature returns signature of the tx.
func (tx *TxImpl) Signature() []byte {
	return tx.signature
}

// Hash returns hash of transaction.
func (tx *TxImpl) Hash() common.Hash {
	return tx.hash
}

// Marshal encodes tx using protobuf
func (tx *TxImpl) Marshal() ([]byte, error) {
	txFrom := tx.from.CloneBytes()
	txTo := tx.to.CloneBytes()
	txHash := tx.hash.CloneBytes()
	txValue := tx.value.Bytes()

	pbTx := &corepb.Transaction{
		Hash:      txHash,
		From:      txFrom,
		To:        txTo,
		Value:     txValue,
		Nonce:     tx.nonce,
		Timestamp: tx.timestamp,
		Signature: tx.signature,
	}

	serializedData, err := proto.Marshal(pbTx)
	if err != nil {
		return nil, errInvalidTransactionToProto
	}
	return serializedData, nil
}

// Unmarshal decode tx using protobuf
func (tx *TxImpl) Unmarshal(data []byte) error {
	pbTx := &corepb.Transaction{}
	err := proto.Unmarshal(data, pbTx)
	if err != nil {
		return errInvalidProtoToTransaction
	}
	tx.hash.SetBytes(pbTx.Hash)
	tx.from.SetBytes(pbTx.From)
	tx.to.SetBytes(pbTx.To)
	// notice: we have to initialize tx.value before pointing to it.
	tx.value = new(big.Int)
	tx.value.SetBytes(pbTx.Value)
	tx.nonce = pbTx.Nonce
	tx.timestamp = pbTx.Timestamp
	tx.signature = pbTx.Signature
	return nil
}

func (tx *TxImpl) String() string {
	return fmt.Sprintf(`{"hash":"%s", "from":"%s", "to":"%s", "nonce":"%v", "value":"%s", "timestamp":"%v"}`,
		tx.hash.String(),
		hex.EncodeToString(tx.From()),
		hex.EncodeToString(tx.To()),
		tx.nonce,
		tx.value,
		tx.timestamp,
	)
}

// Sign signs the tx
func (tx *TxImpl) Sign(kp abstraction.KeyPair) {
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

	hasher.Write(tx.From())
	hasher.Write(tx.To())
	hasher.Write([]byte(value))
	hasher.Write(common.FromUint64(tx.nonce))
	hasher.Write(common.FromInt64(tx.timestamp))

	var h common.Hash
	h.SetBytes(hasher.Sum(nil))

	return h, nil
}

// VerifyIntegrity verifies transaction information
func (tx *TxImpl) VerifyIntegrity() error {
	// verify tx hash
	wantedHash, err := tx.calcHash()
	if err != nil {
		return err
	}
	if wantedHash.Equals(&tx.hash) == false {
		return errInvalidTransacionHash
	}

	// verify signature
	if isValidSignature := tx.Verify(tx.From()); isValidSignature == false {
		return errInvalidTransactionSignature
	}

	return nil
}
