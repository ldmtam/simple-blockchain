package abstraction

import (
	"math/big"

	"github.com/ldmtam/tam-chain/common"
)

// Transaction interface
type Transaction interface {
	Sign(KeyPair)
	Verify([]byte) bool
	VerifyIntegrity() error

	Hash() common.Hash
	Timestamp() int64
	Nonce() uint64
	Fee() *big.Int
	ChainID() uint32
}
