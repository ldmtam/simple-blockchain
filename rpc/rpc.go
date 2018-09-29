package rpc

import (
	"context"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/ldmtam/simpleblockchain/abstraction"
	log "github.com/sirupsen/logrus"
)

// JSONServer json based api rpc server.
type JSONServer struct {
	port     string
	endPoint string
	srv      *http.Server
}

// NewJSONServer returns new instance of JsonServer
func NewJSONServer(endPoint, port string) *JSONServer {
	if !strings.HasPrefix(endPoint, ":") {
		endPoint = "localhost:" + endPoint
	}

	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	return &JSONServer{
		port:     port,
		endPoint: endPoint,
	}
}

// Start the server
func (j *JSONServer) Start(txPool abstraction.TxPool) {
	go func() {
		r := mux.NewRouter()

		r.HandleFunc("/", homeHandler).Methods("GET")

		r.HandleFunc("/generatekeypair", generateKeypairHandler).Methods("POST")

		r.HandleFunc("/createrawtx", createRawTxHandler).Methods("POST")

		r.HandleFunc("/signrawtx", signRawTxHandler).Methods("POST")

		r.HandleFunc("/sendrawtx", func(w http.ResponseWriter, r *http.Request) {
			sendRawTxHandler(w, r, txPool)
		}).Methods("POST")

		j.srv = &http.Server{
			Addr:    j.port,
			Handler: r,
		}

		j.srv.ListenAndServe()
	}()
	log.Info("JSON RPC server started.")
}

// Stop the server
func (j *JSONServer) Stop() {
	err := j.srv.Shutdown(context.Background())
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("JSON RPC shutdown failed.")
	}
	log.Info("JSON RPC server stopped")
}
