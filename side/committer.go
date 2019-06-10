package side

import (
	"strings"
	"time"

	gocodex "github.com/codexnetwork/codex-go"
	"github.com/codexnetwork/codex-go/config"
	"github.com/codexnetwork/codex-go/types"
	"github.com/codexnetwork/trunk/cfg"
	"github.com/codexnetwork/trunk/chainhandler"
	"github.com/codexnetwork/trunk/logger"
)

const (
	sideBlockNumPerSec  = 2
	relayBlockNumPerSec = 2
	maxBlockPerTrx      = 1
	timeBetweenTrx      = relayBlockNumPerSec*3 + 1
	retryTimes          = 100
)

type commitWorker struct {
	committer     types.PermissionLevel
	works         chan commitParam
	client        types.ClientInterface
	ActionToRelay *ActionsToRelay
}

type commitWorkers struct {
	cws []*commitWorker
}

func newCommitWorkers(clientCfg *config.ConfigData, committers []cfg.Relayer, sideChainType types.ClientType) *commitWorkers {
	res := &commitWorkers{
		cws: make([]*commitWorker, 0, len(committers)),
	}
	for _, c := range committers {
		cw := &commitWorker{
			committer: c.RelayAccount,
		}
		cw.Start(clientCfg, sideChainType)
		res.cws = append(res.cws, cw)
	}

	return res
}

var commitWorkerMng *commitWorkers

func InitCommitWorker(clientCfg *config.ConfigData, committers []cfg.Relayer, sideChainType types.ClientType) {
	commitWorkerMng = newCommitWorkers(clientCfg, committers, sideChainType)
}

func (c *commitWorkers) OnBlock(block *chainhandler.Block, actions []chainhandler.Action) {
	for _, cw := range c.cws {
		cw.OnBlock(block, actions)
	}
}

func (c *commitWorker) Start(cfg *config.ConfigData, sideChainType types.ClientType) {
	c.works = make(chan commitParam, 4096)

	for {
		client, err := gocodex.NewClient(types.FORCEIO, cfg)
		if err != nil {
			logger.LogError("create client error, need retry", err)
			time.Sleep(1 * time.Second)
		} else {
			c.client = client

			c.ActionToRelay, err = GetRelayActions(sideChainType)
			if err != nil {
				logger.LogError("get actions to relay err ", err)
				time.Sleep(1 * time.Second)
			} else {
				break
			}
		}
	}

	logger.Infof("start worker loop")

	go func(cc *commitWorker) {
		cc.Loop()
	}(c)
}

func (c *commitWorker) OnBlock(block *chainhandler.Block, actions []chainhandler.Action) {
	cc := commitParam{
		Name:     c.client.Name(cfg.GetRelayCfg().Chain),
		Transfer: c.client.Name(c.committer.Actor),
	}
	cc.FromGeneral(c.ActionToRelay,
		types.NewSwitcherInterface(types.FORCEIO),
		block, actions)

	if cc.IsNeedCommit() {
		c.works <- cc
	}
}

func (c *commitWorker) Loop() {
	ticker := time.NewTicker(timeBetweenTrx * time.Second)
	defer ticker.Stop()
	works2do := make([]commitParam, 0, 4096)
	for {
		select {
		case w := <-c.works:
			works2do = append(works2do, w)
			if len(works2do) >= maxBlockPerTrx {
				c.CommitTrx(works2do)
				works2do = works2do[:0]
				time.Sleep(50 * time.Millisecond)
			}
		default:
		}

		select {
		case <-ticker.C:
			if len(works2do) > 0 {
				c.CommitTrx(works2do)
				works2do = works2do[:0]
			}
		default:
		}
	}
}

func (c *commitWorker) CommitTrx(cps []commitParam) {
	actions := make([]*types.Action, 0, len(cps))

	for _, cp := range cps {
		actions = append(actions, &types.Action{
			Account: "codex.relay",
			Name:    "commit",
			Authorization: []types.PermissionLevel{
				c.committer,
			},
			Data: cp,
		})
	}

	logger.Debugf("commit %s blocks num : %d -> %d",
		string(c.committer.Actor), cps[0].Block.Num, cps[len(cps)-1].Block.Num)

	for idx, act := range actions {
		logger.Debugf("commit %d by %v", cps[idx].Block.Num, act.Data)
	}

	for i := 0; ; i++ {
		if i > 1 {
			time.Sleep(10 * time.Millisecond)
			logger.Warnf("commit err re commit times %d", i)
		}

		pushRes, err := c.client.PushActions(actions...)

		if err != nil {
			c.processCommitErr(err)
		} else {
			err = c.waitCommitComplate(cps, pushRes)
			if err == nil {
				break
			}
		}
	}
}

func (c *commitWorker) processCommitErr(err error) {
	logger.LogError("commit action err", err)
	if strings.Contains(err.Error(), "Transaction took too long") {
		logger.Warnf("need wait chain err by took too long")
		time.Sleep(1 * time.Second)
	}

	if strings.Contains(err.Error(), "RAM") {
		logger.Warnf("need wait other chain err by RAM")
		time.Sleep(1 * time.Second)
	}
}

func (c *commitWorker) waitCommitComplate(cps []commitParam, pushRes *types.PushTransactionFullResp) error {
	logger.Infof("commit to relay %s %d %s, trx id %s",
		pushRes.StatusCode, pushRes.BlockNum, pushRes.BlockID,
		pushRes.TransactionID)

	return nil
}
