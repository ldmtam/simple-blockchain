package txpool

import (
	"sync"

	avl "github.com/emirpasic/gods/trees/avltree"
	"github.com/simpleblockchain/common"
	"github.com/simpleblockchain/core/transaction"
)

type sortedTxMap struct {
	tree  *avl.Tree
	txMap map[string]*transaction.TxImpl
	rw    *sync.RWMutex
}

func compareTx(a, b interface{}) int {
	txA := a.(*transaction.TxImpl)
	txB := b.(*transaction.TxImpl)

	return int(txB.Timestamp() - txA.Timestamp())
}

func newSortedTxMap() *sortedTxMap {
	return &sortedTxMap{
		tree:  avl.NewWith(compareTx),
		txMap: make(map[string]*transaction.TxImpl),
		rw:    new(sync.RWMutex),
	}
}

// Get get the tx
func (st *sortedTxMap) Get(hash common.Hash) *transaction.TxImpl {
	st.rw.Lock()
	defer st.rw.Unlock()
	return st.txMap[hash.String()]
}

// Add add the tx
func (st *sortedTxMap) Add(tx *transaction.TxImpl) {
	st.rw.Lock()
	defer st.rw.Unlock()
	st.tree.Put(tx, true)

	h := tx.Hash()
	st.txMap[h.String()] = tx
}

// Del del the tx
func (st *sortedTxMap) Del(hash common.Hash) {
	st.rw.Lock()
	defer st.rw.Unlock()

	tx := st.txMap[hash.String()]
	if tx == nil {
		return
	}
	st.tree.Remove(tx)
	delete(st.txMap, hash.String())
}

// Size size of the txMap
func (st *sortedTxMap) Size() int {
	st.rw.Lock()
	defer st.rw.Unlock()

	return len(st.txMap)
}
