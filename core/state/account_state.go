package state

import (
	"errors"
	"math/big"

	"github.com/gogo/protobuf/proto"
	"github.com/ldmtam/simpleblockchain/common"
	"github.com/ldmtam/simpleblockchain/proto"
)

// Errors
var (
	ErrBalanceInsufficient = errors.New("balance is sufficient")
	ErrAccountNotFound     = errors.New("cannot find account in storage")

	errInvalidAccountToProto = errors.New("account cannot be converted to protobuf")
	errInvalidProtoToAccount = errors.New("protobuf message cannot be converted to account")
)

type account struct {
	address common.Hash
	balance *big.Int
	nonce   uint64
}

// Marshal encode account struct with protobuf
func (acc *account) Marshal() ([]byte, error) {
	accAddress := acc.address.CloneBytes()
	accBalance := acc.balance.Bytes()

	pbAcc := &corepb.Account{
		Address: accAddress,
		Balance: accBalance,
		Nonce:   acc.nonce,
	}

	serializedData, err := proto.Marshal(pbAcc)
	if err != nil {
		return nil, errInvalidAccountToProto
	}
	return serializedData, nil
}

// Unmarshal decode account struct with bytes
func (acc *account) Unmarshal(data []byte) error {
	pbAcc := &corepb.Account{}
	err := proto.Unmarshal(data, pbAcc)
	if err != nil {
		return errInvalidProtoToAccount
	}
	acc.address.SetBytes(pbAcc.Address)
	acc.balance = new(big.Int)
	acc.balance.SetBytes(pbAcc.Balance)
	acc.nonce = pbAcc.Nonce
	return nil
}

// Balance get account's balance
func (acc *account) Balance() *big.Int {
	return acc.balance
}

// Address get account's address
func (acc *account) Address() common.Hash {
	return acc.address
}

// Nonce get account's nonce
func (acc *account) Nonce() uint64 {
	return acc.nonce
}

// IncreaseNonce increase nonce by 1
func (acc *account) IncreaseNonce() {
	acc.nonce++
}

func (acc *account) AddToBalance(value *big.Int) error {
	newBalance := new(big.Int)
	newBalance.Add(acc.balance, value)
	acc.balance = newBalance
	return nil
}

func (acc *account) SubFromBalance(value *big.Int) error {
	if acc.balance.Cmp(value) == -1 {
		return ErrBalanceInsufficient
	}
	newBalance := new(big.Int)
	newBalance.Sub(acc.balance, value)
	acc.balance = newBalance
	return nil
}
