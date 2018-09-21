package rpc

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/ed25519"

	"github.com/simpleblockchain/core/transaction"
	log "github.com/sirupsen/logrus"
)

type keypair struct {
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
}

type createRawTx struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value string `json:"value"`
	Nonce string `json:"nonce"`
}

type signRawTx struct {
	PrivateKey string `json:"private_key"`
	RawTx      string `json:"raw_tx"`
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello, I am simple blockchain. Nice to meet you ;)")
}

func generateKeypairHandler(w http.ResponseWriter, r *http.Request) {
	pubKeyBytes, privKeyBytes, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Cannot generate ed25519 keypair")
	}
	privKey := hex.EncodeToString(privKeyBytes)
	pubKey := hex.EncodeToString(pubKeyBytes)

	kp := keypair{
		PrivateKey: privKey,
		PublicKey:  pubKey,
	}

	json.NewEncoder(w).Encode(kp)
}

func sendRawTxHandler(w http.ResponseWriter, r *http.Request) {

}

func createRawTxHandler(w http.ResponseWriter, r *http.Request) {
	data := new(createRawTx)
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &data)

	// log.WithFields(log.Fields{
	// 	"from":  data.From,
	// 	"to":    data.To,
	// 	"value": data.Value,
	// 	"nonce": data.Nonce,
	// }).Info("Values received from createRawTxHandler")

	txFrom, err := hex.DecodeString(data.From)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("can not decode `from` field")
	}
	txTo, err := hex.DecodeString(data.To)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("cannot decode `to` field")
	}
	txValue, err := strconv.Atoi(data.Value)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("cannot convert `value` to int")
	}
	txNonce, err := strconv.Atoi(data.Nonce)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("cannot convert `nonce` to int")
	}

	tx, err := transaction.NewTransaction(
		txFrom,
		txTo,
		big.NewInt(int64(txValue)),
		uint64(txNonce),
		time.Now().Unix(),
	)
	if err != nil {
		log.Error("cannot create new raw transaction")
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
		}).Error("cannot encode transaction")
	}

	d := map[string]string{"tx_hash": hex.EncodeToString([]byte(pbMess))}
	json.NewEncoder(w).Encode(d)
}

func signRawTxHandler(w http.ResponseWriter, r *http.Request) {

}
