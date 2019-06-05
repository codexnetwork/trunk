package main

import (
	"github.com/fanyang1988/force-block-ev/blockdb"
	"github.com/fanyang1988/force-go/types"
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
