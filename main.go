package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ldmtam/simpleblockchain/abstraction"
	"github.com/ldmtam/simpleblockchain/core/txpool"
	"github.com/ldmtam/simpleblockchain/rpc"
	log "github.com/sirupsen/logrus"
)

func main() {
	var txp abstraction.TxPool
	txp = txpool.NewTxPImpl()
	txp.Start()

	rpc := rpc.NewJSONServer("0.0.0.0", "3000")
	rpc.Start(txp)

	waitExit()

	rpc.Stop()
	txp.Stop()
}

func waitExit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	i := <-c
	log.WithFields(log.Fields{
		"signal": i,
	}).Info("Server received interrupt, shutting down...")
}
