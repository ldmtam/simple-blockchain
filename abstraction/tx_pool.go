package abstraction

// TxPool interface
type TxPool interface {
	AddTx(Transaction) error
	Start()
	Stop()
	//GetTx(txHash common.Hash) transaction.Transaction
}
