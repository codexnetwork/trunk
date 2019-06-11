package blockdb

import (
	"go.uber.org/zap"

	"github.com/codexnetwork/codex-go/types"
	"github.com/codexnetwork/trunk/logger"
)

const defaultForkBlocksInOneNumSize int = 32

type blockItem struct {
	preBlockIdx int
	block       *types.BlockGeneralInfo
}

type blockInNum struct {
	blockNum uint32
	blocks   []blockItem
}

func (b *blockInNum) find(block *types.BlockGeneralInfo) int {
	for idx, bk := range b.blocks {
		if IsBlockEq(bk.block, block) {
			return idx
		}
	}
	return -1
}

func (b *blockInNum) findByID(blockID types.Checksum256) int {
	for idx, bk := range b.blocks {
		if IsChecksum256Eq(blockID, bk.block.ID) {
			return idx
		}
	}
	return -1
}

func (b *blockInNum) append(bi blockItem) int {
	idx := b.find(bi.block)
	if idx >= 0 {
		return idx
	}

	b.blocks = append(b.blocks, bi)
	return len(b.blocks) - 1
}

type PeerBlockState struct {
	peer   string
	blocks []blockInNum
}

func (p *PeerBlockState) init() {
	p.blocks = make([]blockInNum, 0, 512)
}

func (p *PeerBlockState) newBlockNum() uint32 {
	if len(p.blocks) == 0 {
		return 0
	}
	return p.blocks[len(p.blocks)-1].blockNum
}

func (p *PeerBlockState) appendBlock(block *types.BlockGeneralInfo) error {
	// TODO just a simple imp
	blockNum := block.BlockNum
	blockID := block.ID

	//logger.Logger().Debug("block append",
	//	zap.Uint32("num", blockNum), zap.String("id", blockID.String()))

	if len(p.blocks) == 0 {
		// no init
		p.blocks = append(p.blocks, blockInNum{
			blockNum: blockNum,
			blocks:   make([]blockItem, 0, defaultForkBlocksInOneNumSize),
		})
	}

	newBlockNum := p.newBlockNum()
	firstBlockNum := p.blocks[0].blockNum

	if blockNum < firstBlockNum {
		// no need process
		logger.Logger().Debug("no need process",
			zap.Uint32("num", blockNum),
			zap.String("id", blockID.String()),
			zap.String("reason", "blockNum < firstBlockNum"))
		return nil
	}

	if blockNum > (newBlockNum + 1) {
		// no need process
		logger.Logger().Debug("no need process",
			zap.Uint32("num", blockNum),
			zap.String("id", blockID.String()),
			zap.String("reason", "blockNum > (newBlockNum + 1)"))
		return nil
	}

	if blockNum == (newBlockNum + 1) {
		for i := 0; i < int(blockNum-newBlockNum); i++ {
			p.blocks = append(p.blocks, blockInNum{
				blockNum: newBlockNum + uint32(i) + 1,
				blocks:   make([]blockItem, 0, defaultForkBlocksInOneNumSize),
			})
		}
	}

	var perBlockIdx int = -1
	if blockNum > firstBlockNum {
		perBlockItemIdx := blockNum - firstBlockNum - 1
		perBlockIdx = p.blocks[perBlockItemIdx].findByID(block.Previous)
		if perBlockIdx < 0 {
			// no process no pre block id
			logger.Logger().Debug("no need process",
				zap.Uint32("num", blockNum),
				zap.String("id", blockID.String()),
				zap.String("reason", "no process no pre block id"))
			return nil
		}
	}

	p.blocks[blockNum-firstBlockNum].append(blockItem{
		perBlockIdx,
		block,
	})

	//p.debugLogStat()

	return nil
}

func (p *PeerBlockState) debugLogStat() {
	logger.Logger().Sugar().Infof("block stat %s [%d] form %d to %d",
		p.peer, len(p.blocks), p.blocks[0].blockNum, p.newBlockNum())
	/*
		for _, b := range p.blocks {
			for _, bb := range b.blocks {
				logger.Logger().Sugar().Infof("block %d %s (%d) with %d trx",
					bb.blockNum, bb.blockID.String(), bb.preBlockIdx, len(bb.block.Transactions))
			}
		}
	*/
}

// GetFirstBlock
func (p *PeerBlockState) GetBlock(num uint32) *blockItem {
	if len(p.blocks) == 0 {
		return nil
	}

	if num < p.blocks[0].blockNum {
		return nil
	}

	idx := int(num - p.blocks[0].blockNum)
	if idx >= len(p.blocks) {
		return nil
	}

	if len(p.blocks[idx].blocks) >= 1 {
		return &p.blocks[idx].blocks[0]
	}

	// TODO if no sync
	return nil
}

// DelFirstBlock
func (p *PeerBlockState) DelBlockBefore(num uint32) {
	if len(p.blocks) == 0 {
		return
	}

	if num < p.blocks[0].blockNum {
		return
	}

	idx := int(num - p.blocks[0].blockNum)

	if idx >= len(p.blocks) {
		return
	}

	num2del := idx + 1
	if len(p.blocks) > num2del {
		for i := 0; i < (len(p.blocks) - num2del); i++ {
			p.blocks[i] = p.blocks[i+num2del]
		}
		p.blocks = p.blocks[:len(p.blocks)-num2del]
	}

	if len(p.blocks) == num2del {
		p.blocks = p.blocks[0:0]
	}

	//p.debugLogStat()

	return
}

func (p *PeerBlockState) BlockLen() int {
	return len(p.blocks)
}
