package rpc

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/ed25519"

	log "github.com/inconshreveable/log15"
	"github.com/ldmtam/tam-chain/abstraction"
	"github.com/ldmtam/tam-chain/account"
	"github.com/ldmtam/tam-chain/common"
	"github.com/ldmtam/tam-chain/core/transaction"
)

func renderErrorMessage(err error, w http.ResponseWriter) {
	d := map[string]string{"error": err.Error()}
	json.NewEncoder(w).Encode(d)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello, I am simple blockchain. Nice to meet you ;)")
}

func generateKeypairHandler(w http.ResponseWriter, r *http.Request) {
	type keypair struct {
		PrivateKey string `json:"private_key"`
		PublicKey  string `json:"public_key"`
	}

	pubKeyBytes, privKeyBytes, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		log.Error("Cannot generate ed25519 keypair", "error", err)

		renderErrorMessage(err, w)
		return
	}
	privKey := hex.EncodeToString(privKeyBytes)
	pubKey := hex.EncodeToString(pubKeyBytes)

	kp := keypair{
		PrivateKey: privKey,
		PublicKey:  pubKey,
	}

	json.NewEncoder(w).Encode(kp)
}

func createRawTxHandler(w http.ResponseWriter, r *http.Request) {
	type createRawTx struct {
		ChainID string `json:"chainid"`
		From    string `json:"from"`
		To      string `json:"to"`
		Value   string `json:"value"`
		Fee     string `json:"fee"`
		Nonce   string `json:"nonce"`
	}

	data := new(createRawTx)
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &data)

	if data.From == data.To {
		errString := `from and to address must not be the same`
		log.Error("addresses are the same", "error", errString)

		renderErrorMessage(errors.New(errString), w)
		return
	}

	from := &account.KeyPairImpl{}
	to := &account.KeyPairImpl{}

	err := from.DecodePublicKey(data.From)
	if err != nil {
		log.Error("can not decode `from` field", "error", err)

		renderErrorMessage(err, w)
		return
	}

	err = to.DecodePublicKey(data.To)
	if err != nil {
		log.Error("cannot decode `to` field", "error", err)

		renderErrorMessage(err, w)
		return
	}

	txChainID, err := strconv.Atoi(data.ChainID)
	if err != nil {
		log.Error("cannot convert `chainID` to int", "error", err)

		renderErrorMessage(err, w)
		return
	}

	txValue, err := strconv.Atoi(data.Value)
	if err != nil {
		log.Error("cannot convert `value` to int", "error", err)

		renderErrorMessage(err, w)
		return
	}

	txFee, err := strconv.Atoi(data.Fee)
	if err != nil {
		log.Error("cannot convert `fee` to int", "error", err)

		renderErrorMessage(err, w)
		return
	}

	txNonce, err := strconv.Atoi(data.Nonce)
	if err != nil {
		log.Error("cannot convert `nonce` to int", "error", err)

		renderErrorMessage(err, w)
		return
	}

	var txFrom common.Address
	copy(txFrom[:], from.PublicKey)
	var txTo common.Address
	copy(txTo[:], to.PublicKey)

	tx, err := transaction.NewTransaction(
		uint32(txChainID),
		txFrom,
		txTo,
		big.NewInt(int64(txValue)),
		big.NewInt(int64(txFee)),
		uint64(txNonce),
		time.Now().Unix(),
	)
	if err != nil {
		log.Error("cannot create new raw transaction", "error", err)

		renderErrorMessage(err, w)
		return
	}

	pbMess, err := tx.Marshal()
	if err != nil {
		log.Error("cannot encode transaction with protobuf", "error", err)

		renderErrorMessage(err, w)
		return
	}

	d := map[string]string{"unsigned_raw_tx": hex.EncodeToString(pbMess)}
	json.NewEncoder(w).Encode(d)
}

func signRawTxHandler(w http.ResponseWriter, r *http.Request) {
	type signRawTx struct {
		PrivateKey    string `json:"private_key"`
		UnsignedRawTx string `json:"unsigned_raw_tx"`
	}

	data := new(signRawTx)
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &data)

	tx := new(transaction.TxImpl)
	rawTxBytes, err := hex.DecodeString(data.UnsignedRawTx)
	if err != nil {
		log.Error("cannot convert tx hex string to bytes", "error", err)

		renderErrorMessage(err, w)
		return
	}

	err = tx.Unmarshal(rawTxBytes)
	if err != nil {
		log.Error("cannot decode transaction with protobuf", "error", err)

		renderErrorMessage(err, w)
		return
	}

	kp := &account.KeyPairImpl{}

	err = kp.DecodePrivateKey(data.PrivateKey)
	if err != nil {
		log.Error("cannot decode private key hex string to bytes", "error", err)

		renderErrorMessage(err, w)
		return
	}

	tx.Sign(kp)

	txBytes, err := tx.Marshal()
	if err != nil {
		log.Error("cannot encode tx with protobuf", "error", err)

		renderErrorMessage(err, w)
		return
	}

	d := map[string]string{"raw_tx": hex.EncodeToString(txBytes)}
	json.NewEncoder(w).Encode(d)
}

func sendRawTxHandler(w http.ResponseWriter, r *http.Request, txPool abstraction.TxPool) {
	type sendRawTx struct {
		RawTx string `json:"raw_tx"`
	}

	data := new(sendRawTx)
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &data)

	txBytes, err := hex.DecodeString(data.RawTx)
	if err != nil {
		renderErrorMessage(err, w)
		return
	}
	tx := &transaction.TxImpl{}
	tx.Unmarshal(txBytes)

	err = txPool.AddTx(tx, true)
	if err != nil {
		log.Error("cannot add tx to tx pool", "error", err)

		renderErrorMessage(err, w)
		return
	}

	d := map[string]string{"result": "success"}
	json.NewEncoder(w).Encode(d)
}
