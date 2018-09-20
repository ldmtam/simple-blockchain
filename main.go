package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/simpleblockchain/rpc"
	log "github.com/sirupsen/logrus"
)

func main() {
	rpc := rpc.NewJSONServer("0.0.0.0", "3000")
	rpc.Start()

	waitExit()
	rpc.Stop()
}

func waitExit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	i := <-c
	log.WithFields(log.Fields{
		"signal": i,
	}).Info("Server received interrupt, shutting down...")
}
