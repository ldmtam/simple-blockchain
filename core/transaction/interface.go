package transaction

import "github.com/simpleblockchain/account"

// Transaction interface of transaction type
type Transaction interface {
	Sign(account.KeyPair)
	Verify([]byte) bool
}
