package abstraction

import (
	"math/big"

	"github.com/simpleblockchain/common"
)

// Account interface
type Account interface {
	Address() common.Hash
	Balance() *big.Int
	Nonce() uint64

	Marshal() ([]byte, error)
	Unmarshal([]byte) error

	IncreaseNonce()
	AddToBalance(*big.Int) error
	SubFromAccount(*big.Int) error
}
