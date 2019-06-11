package main

import (
	"github.com/codexnetwork/codex-go/types"
	"github.com/codexnetwork/trunk/blockdb"
)

type handlerImp struct {
	verifier *blockdb.FastBlockVerifier
}

func (h *handlerImp) OnBlock(peer string, msg *types.BlockGeneralInfo) error {
	//logger.Debugf("on b %s", msg.BlockNum)
	return h.verifier.OnBlock(peer, msg)
}

func (h *handlerImp) OnGoAway(peer string, reason uint8, nodeID types.Checksum256) error {
	return nil
}
