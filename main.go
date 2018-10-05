package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/inconshreveable/log15"
	"github.com/ldmtam/tam-chain/abstraction"
	"github.com/ldmtam/tam-chain/common"
	"github.com/ldmtam/tam-chain/core/txpool"
	"github.com/ldmtam/tam-chain/p2p"
	"github.com/ldmtam/tam-chain/rpc"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "port",
			Usage: "node port",
		},
		cli.StringFlag{
			Name:  "datapath",
			Usage: "data path",
		},
		cli.StringFlag{
			Name:  "bootnode",
			Usage: "list of boot nodes",
		},
	}

	app.Action = func(c *cli.Context) error {
		p2pConfig := &common.P2PConfig{
			ChainID:   1,
			Version:   1,
			Port:      c.String("port"),
			SeedNodes: []string{c.String("bootnode")},
			DataPath:  c.String("datapath"),
		}

		var net abstraction.P2PService
		net, _ = p2p.NewNetService(p2pConfig)
		net.Start()

		var txp abstraction.TxPool
		txp = txpool.NewTxPImpl()
		txp.Start()

		rpc := rpc.NewJSONServer("0.0.0.0", "3000")
		rpc.Start(txp)

		waitExit()

		rpc.Stop()
		txp.Stop()
		net.Stop()

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Error("some errors happen", "err", err)
	}
}

func waitExit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	i := <-c
	log.Info("Server received interrupt, shutting down...", "signal", i)
}
