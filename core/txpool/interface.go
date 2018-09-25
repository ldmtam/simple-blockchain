package txpool

import (
	"github.com/simpleblockchain/common"
	"github.com/simpleblockchain/core/transaction"
)

// TxPool ...
type TxPool interface {
	AddTx(tx *transaction.TxImpl) bool
	GetTx(txHash common.Hash) *transaction.TxImpl
}
