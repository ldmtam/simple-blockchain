package txpool

import (
	"sync"

	log "github.com/inconshreveable/log15"
	"github.com/ldmtam/tam-chain/abstraction"
	"github.com/ldmtam/tam-chain/common"
)

// TxPImpl ...
type TxPImpl struct {
	all    *txLookup // All transaction to look up
	fee    *sortedTx // All transaction sorted by fee
	locals map[common.Hash]abstraction.Transaction

	mu     sync.RWMutex
	quitCh chan struct{}
}

// NewTxPImpl returns a new TxPImpl instance.
func NewTxPImpl() *TxPImpl {
	return &TxPImpl{
		all:    newTxLookup(),
		fee:    newSortedTx(),
		locals: make(map[common.Hash]abstraction.Transaction),
		quitCh: make(chan struct{}),
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
	for {
		select {
		case <-pool.quitCh:
			return
		}
	}
}

// verifyTx verifies tx before adding it to tx pool.
//
// [DONE] step 1: check whether the signature belongs to `from` address or not.
// [DONE] step 2: recalculate the tx hash and check if it matches with the tx hash sent by user.
// [TODO] step 3: check whether tx nonce = `from` nonce + 1 or not
// [TODO] step 4: check whether `from` balance is greater than or equal to (tx value + tx fee) or not.
func (pool *TxPImpl) verifyTx(tx abstraction.Transaction) error {
	// step 1 & 2.
	if err := tx.VerifyIntegrity(); err != nil {
		return err
	}

	return nil
}

// AddTx add transaction to tx pool.
func (pool *TxPImpl) AddTx(tx abstraction.Transaction, local bool) error {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	err := pool.verifyTx(tx)
	if err != nil {
		return err
	}

	pool.all.Add(tx)
	pool.fee.Push(tx)
	if local == true {
		pool.locals[tx.Hash()] = tx
	}

	return nil
}

// DelTx remove a tx from the tx pool.
func (pool *TxPImpl) DelTx(hash common.Hash) {
	pool.all.Remove(hash)
	pool.fee.Delete(pool.all.Get(hash))
}
