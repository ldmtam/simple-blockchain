package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/inconshreveable/log15"
	"github.com/ldmtam/tam-chain/abstraction"
	"github.com/ldmtam/tam-chain/core/txpool"
	"github.com/ldmtam/tam-chain/rpc"
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
	log.Info("Server received interrupt, shutting down...", "signal", i)
}
