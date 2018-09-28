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

	"github.com/simpleblockchain/abstraction"
	"github.com/simpleblockchain/account"
	"github.com/simpleblockchain/core/transaction"
	log "github.com/sirupsen/logrus"
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
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Cannot generate ed25519 keypair")

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
		From  string `json:"from"`
		To    string `json:"to"`
		Value string `json:"value"`
		Nonce string `json:"nonce"`
	}

	data := new(createRawTx)
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &data)

	if data.From == data.To {
		errString := `from and to address must not be the same`
		log.WithFields(log.Fields{
			"error": errString,
		}).Error("addresses are the same")

		renderErrorMessage(errors.New(errString), w)
		return
	}

	from := &account.KeyPairImpl{}
	to := &account.KeyPairImpl{}

	err := from.DecodePublicKey(data.From)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("can not decode `from` field")

		renderErrorMessage(err, w)
		return
	}

	err = to.DecodePublicKey(data.To)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("cannot decode `to` field")

		renderErrorMessage(err, w)
		return
	}

	txValue, err := strconv.Atoi(data.Value)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("cannot convert `value` to int")

		renderErrorMessage(err, w)
		return
	}

	txNonce, err := strconv.Atoi(data.Nonce)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("cannot convert `nonce` to int")

		renderErrorMessage(err, w)
		return
	}

	txFrom := make([]byte, ed25519.PublicKeySize)
	copy(txFrom, from.PublicKey[:])
	txTo := make([]byte, ed25519.PublicKeySize)
	copy(txTo, to.PublicKey[:])

	tx, err := transaction.NewTransaction(
		txFrom,
		txTo,
		big.NewInt(int64(txValue)),
		uint64(txNonce),
		time.Now().Unix(),
	)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("cannot create new raw transaction")

		renderErrorMessage(err, w)
		return
	}

	txHashInHash := tx.Hash()
	txHashInBytes := txHashInHash.CloneBytes()

	log.WithFields(log.Fields{
		"from":  data.From,
		"to":    data.To,
		"nonce": txNonce,
		"value": data.Value,
		"hash":  hex.EncodeToString(txHashInBytes),
	}).Info("Transaction data")

	pbMess, err := tx.Marshal()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("cannot encode transaction with protobuf")

		renderErrorMessage(err, w)
		return
	}

	d := map[string]string{"raw_tx": hex.EncodeToString(pbMess)}
	json.NewEncoder(w).Encode(d)
}

func signRawTxHandler(w http.ResponseWriter, r *http.Request) {
	type signRawTx struct {
		PrivateKey string `json:"private_key"`
		RawTx      string `json:"raw_tx"`
	}

	data := new(signRawTx)
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &data)

	log.WithFields(log.Fields{
		"Private key": data.PrivateKey,
		"Raw Tx":      data.RawTx,
	}).Info("Value received from sign raw tx")

	tx := new(transaction.TxImpl)
	rawTxBytes, err := hex.DecodeString(data.RawTx)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("cannot convert tx hex string to bytes")

		renderErrorMessage(err, w)
		return
	}

	err = tx.Unmarshal(rawTxBytes)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("cannot decode transaction with protobuf")

		renderErrorMessage(err, w)
		return
	}

	kp := &account.KeyPairImpl{}

	err = kp.DecodePrivateKey(data.PrivateKey)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("cannot decode private key hex string to bytes")

		renderErrorMessage(err, w)
		return
	}

	tx.Sign(kp)

	log.WithFields(log.Fields{
		"Signature": hex.EncodeToString(tx.Signature()),
	}).Info("Signature of the message")

	txBytes, err := tx.Marshal()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("cannot encode tx with protobuf")

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

	txHashInHash := tx.Hash()
	txHashInBytes := txHashInHash.CloneBytes()

	log.WithFields(log.Fields{
		"Hash":      hex.EncodeToString(txHashInBytes),
		"From":      hex.EncodeToString(tx.From()),
		"To":        hex.EncodeToString(tx.To()),
		"Value":     tx.Value(),
		"Nonce":     tx.Nonce(),
		"Timestamp": tx.Timestamp(),
		"Signature": hex.EncodeToString(tx.Signature()),
	}).Info("Info of raw tx")

	err = txPool.AddTx(tx)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("cannot add tx to tx pool")

		renderErrorMessage(err, w)
		return
	}

	d := map[string]string{"result": "success"}
	json.NewEncoder(w).Encode(d)
}
