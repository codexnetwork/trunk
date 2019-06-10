package main

import (
	"github.com/fanyang1988/force-block-ev/blockdb"
	"github.com/codexnetwork/codex-go/p2p"
	"github.com/pkg/errors"

	"github.com/codexnetwork/trunk/cfg"
	"github.com/codexnetwork/trunk/chainhandler"
	"github.com/codexnetwork/trunk/logger"
	"github.com/codexnetwork/trunk/relay"
	"github.com/codexnetwork/trunk/side"
)

func startSideService() {
	// frome side need to commit block to relay
	chainCfgs, _ := cfg.GetChainCfg("relay")

	data, p2ps := cfg.GetChainCfg("side")
	chainTyp := cfg.GetChainTyp("side")

	side.InitCommitWorker(chainCfgs, cfg.GetTransfers(), chainTyp)

	// for p2p chain id
	info, err := relay.Client().GetInfoData()
	if err != nil {
		panic(errors.New("get info err"))
	}

	lastCommitted, err := side.GetLastCommittedBlock()
	if err != nil {
		logger.Errorf("err by %s", err.Error())
		panic(errors.New("GetLastCommittedBlock info err"))
	}

	logger.Debugf("get last committed block %v %d", lastCommitted, data.StartNum)

	lastNum := lastCommitted.Num
	if lastNum > 3 {
		lastNum -= 2
	}

	if lastNum == 0 {
		lastNum = data.StartNum
	}

	p2pPeers := p2p.NewP2PClient(chainTyp, p2p.P2PInitParams{
		Name:          "relay",
		ClientID:      info.ChainID.String(),
		StartBlockNum: lastNum,
		Peers:         p2ps,
		Logger:        logger.Logger(),
	})

	p2pPeers.RegHandler(&handlerImp{
		verifier: blockdb.NewFastBlockVerifier(p2ps, lastNum, chainhandler.NewChainHandler(
			func(block *chainhandler.Block, actions []chainhandler.Action) {
				side.HandSideBlock(block, actions)
			}, chainTyp)),
	})
	p2pPeers.Start()
}
