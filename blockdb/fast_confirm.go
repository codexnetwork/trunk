package blockdb

import (
	"go.uber.org/zap"

	"github.com/codexnetwork/codex-go/types"
	"github.com/codexnetwork/trunk/logger"
)

type VerifyHandler interface {
	OnBlock(block *types.BlockGeneralInfo) error
}

type FastBlockVerifier struct {
	db            *BlockDB
	verifyHandler VerifyHandler
	lastVerifyNum uint32
	startBlock    uint32
}

// NewFastBlockVerifier create
func NewFastBlockVerifier(peers []string, startBlock uint32, verifyHandler VerifyHandler) *FastBlockVerifier {
	blocks := &BlockDB{}
	blocks.Init(peers)
	res := &FastBlockVerifier{
		db:            blocks,
		verifyHandler: verifyHandler,
		startBlock:    startBlock,
	}
	if startBlock > 1 {
		res.lastVerifyNum = startBlock - 1
	}
	return res
}

// CheckBlocks Call when append block to db
func (f *FastBlockVerifier) checkBlocks() error {
	bi, ok := f.TryGetVerifyBlock()
	if ok {
		err := f.verifyHandler.OnBlock(bi.block)
		if err != nil {
			// no del block
			return err
		}
		f.lastVerifyNum = bi.block.BlockNum
		f.db.DelBlockBefore(f.lastVerifyNum - 3)

		return nil
	}

	return nil
}

// OnBlock
func (f *FastBlockVerifier) OnBlock(peer string, block *types.BlockGeneralInfo) error {
	if block.BlockNum > f.lastVerifyNum {
		err := f.db.OnBlock(peer, block)
		if err != nil {
			return err
		}
	}
	return f.checkBlocks()
}

func (f *FastBlockVerifier) TryGetVerifyBlock() (*blockItem, bool) {
	// TODO now is a simple imp
	num2Verify := f.lastVerifyNum + 1
	var bi *blockItem
	for peer, stat := range f.db.Peers {
		// 1. every peer blocks are no fork
		for _, blocks := range stat.blocks {
			if len(blocks.blocks) > 1 {
				logger.Logger().Info("peer block forking",
					zap.String("peer", peer),
					zap.Uint32("num", blocks.blockNum),
					zap.Int("fork", len(blocks.blocks)))
				return nil, false
			}
		}

		// 2. all peer are same
		first := stat.GetBlock(num2Verify)
		if bi == nil && first != nil {
			bi = first
		} else {
			if (first == nil) || (bi != nil && !IsBlockEq(bi.block, first.block)) {
				return nil, false
			}
		}
	}

	if bi == nil {
		return nil, false
	}

	return bi, true
}
