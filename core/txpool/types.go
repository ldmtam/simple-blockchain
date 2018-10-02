package txpool

import (
	"sync"

	"github.com/ldmtam/tam-chain/abstraction"

	"github.com/ldmtam/tam-chain/common"
)

type txLookup struct {
	all  map[common.Hash]abstraction.Transaction
	lock sync.RWMutex
}

func newTxLookup() *txLookup {
	return &txLookup{
		all: make(map[common.Hash]abstraction.Transaction),
	}
}

func (t *txLookup) Get(hash common.Hash) abstraction.Transaction {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return t.all[hash]
}

func (t *txLookup) Count() int {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return len(t.all)
}

func (t *txLookup) Add(tx abstraction.Transaction) {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.all[tx.Hash()] = tx
}

func (t *txLookup) Remove(hash common.Hash) {
	t.lock.Lock()
	defer t.lock.Unlock()

	delete(t.all, hash)
}
