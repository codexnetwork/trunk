package chainhandler

import (
	"sync"

	"github.com/codexnetwork/codex-go/types"
	"github.com/codexnetwork/trunk/logger"
)

type HandlerFunc func(block *Block, actions []Action)

type ChainHandler struct {
	handler    HandlerFunc
	blockQueue chan blockQueueItem
	wg         sync.WaitGroup
	switcher   types.SwitcherInterface
}

func NewChainHandler(h HandlerFunc, chainTyp types.ClientType) *ChainHandler {
	res := &ChainHandler{
		handler:    h,
		blockQueue: make(chan blockQueueItem, 4096),
		switcher:   types.NewSwitcherInterface(chainTyp),
	}

	res.wg.Add(1)
	go func(ch *ChainHandler) {
		defer ch.wg.Done()
		logger.Logger().Info("start chain handler")
		for {
			bi, ok := <-ch.blockQueue
			if !ok {
				logger.Logger().Error("handler chan close")
				return
			}
			ch.handler(&bi.block, bi.actions)
		}
	}(res)

	return res
}

func (c *ChainHandler) OnBlock(block *types.BlockGeneralInfo) error {
	//logger.Logger().Debug("on block",
	//	zap.Uint32("num", block.BlockNum),
	//	zap.String("id", block.ID.String()))
	var bqi blockQueueItem
	bqi.block = Block{
		Producer:         block.Producer,
		Num:              block.BlockNum,
		ID:               block.ID,
		Previous:         block.Previous,
		Confirmed:        uint16(block.Confirmed),
		TransactionMRoot: block.TransactionMRoot,
		ActionMRoot:      block.ActionMRoot,
	}

	actions := make([]Action, 0, 1024)
	for _, trx := range block.Transactions {
		if trx.Status != types.TransactionStatusExecuted {
			continue
		}

		for _, act := range trx.Transaction.Actions {
			auth := make([]PermissionLevel, 0, 8)
			for _, authToTrx := range act.Authorization {
				auth = append(auth, PermissionLevel{
					Actor:      authToTrx.Actor,
					Permission: authToTrx.Permission,
				})
			}
			actions = append(actions, Action{
				Account:       act.Account,
				Name:          act.Name,
				Authorization: auth,
				Data:          act.HexData,
			})
		}
	}
	bqi.actions = actions[:]
	c.blockQueue <- bqi
	return nil
}

func (c *ChainHandler) Close() {
	close(c.blockQueue)
	c.wg.Wait()
}
