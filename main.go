package main

import (
	"flag"
	"fmt"
	"runtime"
	"time"

	"github.com/eosforce/goforceio/ecc"
	"github.com/eosforce/goforceio/p2p"

	"github.com/codexnetwork/trunk/cfg"
	"github.com/codexnetwork/trunk/logger"
	"github.com/codexnetwork/trunk/relay"
	"github.com/codexnetwork/trunk/side"
)

var configPath = flag.String("cfg", "./config.json", "config file path")
var isDebug = flag.Bool("d", false, "run in debug mode")

func init() {
	ecc.PublicKeyPrefixCompat = "CDX"
}

func main() {
	flag.Parse()
	logger.EnableLogging(*isDebug)

	if *isDebug {
		p2p.EnableP2PLogging()
	}

	defer func() {
		err := logger.Logger().Sync()
		if err != nil {
			fmt.Printf("logger sync err by %s", err.Error())
		}
	}()

	runtime.GOMAXPROCS(2)

	err := cfg.LoadCfgs(*configPath)
	if err != nil {
		logger.Sugar().Errorf("load cfg err by %s", err.Error())
		return
	}

	sideChainCfg, _ := cfg.GetChainCfg("side")
	sideChainTyp := cfg.GetChainTyp("side")
	relay.CreateSideClient(sideChainTyp, sideChainCfg)

	relayChainCfg, _ := cfg.GetChainCfg("relay")
	relayChainTyp := cfg.GetChainTyp("relay")
	side.CreateClient(relayChainTyp, relayChainCfg)

	go func() {
		if len(cfg.GetWatchers()) == 0 {
			logger.Sugar().Infof("no need start relay")
			return
		}
		logger.Sugar().Infof("start relay service")
		startRelayService()
	}()

	go func() {
		if len(cfg.GetTransfers()) == 0 {
			logger.Sugar().Infof("no need start side")
			return
		}
		logger.Sugar().Infof("start side service")
		startSideService()
	}()

	for {
		time.Sleep(1 * time.Second)
		// TODO check status
	}
}
