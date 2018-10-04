package txpool

import (
	"sync"

	"github.com/ldmtam/tam-chain/common/sorted"

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

type sortedTx struct {
	txsByFee *sorted.Slice
	mu       sync.RWMutex
}

func feeCmp(a, b interface{}) int {
	txa := a.(abstraction.Transaction)
	txb := b.(abstraction.Transaction)
	if txa.Fee().Cmp(txb.Fee()) == 0 {
		if txa.Nonce() <= txb.Nonce() {
			return 1
		}
		return -1
	}
	return txa.Fee().Cmp(txb.Fee())
}

// NewSortedTx returns new instance of sorted tx
func newSortedTx() *sortedTx {
	return &sortedTx{
		txsByFee: sorted.NewSlice(feeCmp),
	}
}

func (s *sortedTx) Push(tx abstraction.Transaction) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.txsByFee.Push(tx)
}

func (s *sortedTx) PopRight() abstraction.Transaction {
	s.mu.Lock()
	defer s.mu.Unlock()

	txInterface := s.txsByFee.PopRight()
	return txInterface.(abstraction.Transaction)
}

func (s *sortedTx) PopLeft() abstraction.Transaction {
	s.mu.Lock()
	defer s.mu.Unlock()

	txInterface := s.txsByFee.PopLeft()
	return txInterface.(abstraction.Transaction)
}

func (s *sortedTx) Delete(tx abstraction.Transaction) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.txsByFee.Del(tx)
}

func (s *sortedTx) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.txsByFee.Len()
}
