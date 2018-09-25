package txpool

import (
	"github.com/simpleblockchain/core/transaction"
	log "github.com/sirupsen/logrus"
)

// TxPImpl ...
type TxPImpl struct {
	pendingTx *sortedTxMap

	quitCh chan struct{}
}

// NewTxPImpl returns a new TxPImpl instance.
func NewTxPImpl() *TxPImpl {
	return &TxPImpl{
		pendingTx: newSortedTxMap(),
		quitCh:    make(chan struct{}),
	}
}

// Start starts the tx pool.
func (pool *TxPImpl) Start() {
	go pool.loop()
}

// Stop stops the tx pool.
func (pool *TxPImpl) Stop() {
	log.Info("Tx pool stop")
	close(pool.quitCh)
}

func (pool *TxPImpl) loop() {

}

func (pool *TxPImpl) verifyTx(tx *transaction.TxImpl) bool {
	return true
}

// AddTx add transaction to tx pool.
func (pool *TxPImpl) AddTx(tx *transaction.TxImpl) bool {
	pool.pendingTx.Add(tx)
	return true
}
