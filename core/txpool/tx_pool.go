package txpool

import (
	log "github.com/inconshreveable/log15"
	"github.com/ldmtam/tam-chain/abstraction"
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

// verifyTx verifies tx before adding it to tx pool.
//
// [DONE] step 1: check whether the signature belongs to `from` address or not.
// [DONE] step 2: recalculate the tx hash and check if it matches with the tx hash sent by user.
// [TODO] step 3: check whether tx nonce = `from` nonce + 1 or not
// [TODO] step 4: check whether `from` balance is greater than or equal to tx value or not.
func (pool *TxPImpl) verifyTx(tx abstraction.Transaction) error {
	// step 1 & 2.
	if err := tx.VerifyIntegrity(); err != nil {
		return err
	}

	return nil
}

// AddTx add transaction to tx pool.
func (pool *TxPImpl) AddTx(tx abstraction.Transaction) error {
	err := pool.verifyTx(tx)
	if err != nil {
		return err
	}

	return nil
}
