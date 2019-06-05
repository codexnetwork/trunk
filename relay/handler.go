package relay

import (
	"fmt"

	"github.com/fanyang1988/force-go/types"
	"go.uber.org/zap"

	"github.com/codexnetwork/trunk/cfg"
	"github.com/codexnetwork/trunk/chaindata"
	"github.com/codexnetwork/trunk/chainhandler"
	"github.com/codexnetwork/trunk/logger"
	"github.com/codexnetwork/trunk/side"
	forceio "github.com/eosforce/goforceio"
)

// Destroy action data for relay.token::destroy action
type Destroy struct {
	Chain    forceio.Name        `json:"chain"`
	From     forceio.AccountName `json:"from"`
	Quantity forceio.Asset       `json:"quantity"`
	Memo     string              `json:"memo"`
}

func HandRelayBlock(block *chainhandler.Block, actions []chainhandler.Action) {
	if block.GetNum()%1000 == 0 {
		logger.Debugf("on block from relay %d", block.GetNum())
	}
	for idx, act := range actions {
		if act.Account != "relay.token" || act.Name != "destroy" {
			continue
		}

		var actData Destroy
		err := forceio.UnmarshalBinary(act.Data, &actData)
		if err != nil {
			logger.LogError("UnmarshalBinary act err", err)
			continue
		}

		onTokenReturnSideChain(block, idx, &actData)
	}
}

func onTokenReturnSideChain(block *chainhandler.Block, idx int, act *Destroy) {
	logger.Debugf("on return in block %d : %s %v by %v in %s",
		block.GetNum(), act.Chain, act.From, act.Quantity, act.Memo)

	for _, w := range cfg.GetWatchers() {
		err := commitOutAction(w, block.Num, idx, act)
		if err != nil {
			logger.Logger().Error("commit out action err", zap.Error(err))
		}
	}
}

// OutAction  capi_name committer, uint64_t num, capi_name to, name chain, name contract, const asset& quantity, const std::string& memo
type OutAction struct {
	Committer interface{} `json:"committer"`
	Num       uint64      `json:"num"`
	To        interface{} `json:"to"`
	Chain     interface{} `json:"chain"`
	Contract  interface{} `json:"contract"`
	Action    interface{} `json:"action"`
	Quantity  interface{} `json:"quantity"`
	Memo      string      `json:"memo"`
}

// just use a large num
const maxActionInBlock = 100000

func commitOutAction(committer cfg.Relayer, blockNum uint32, idx int, act *Destroy) error {
	mapTokenStat, err := chaindata.GetTokenMapData(string(act.Chain), act.Quantity.Symbol.Symbol, side.Client())
	if err != nil {
		return err
	}

	actToCommit := &types.Action{
		Account: cfg.GetRelayCfg().RelayContract,
		Name:    "out",
		Authorization: []types.PermissionLevel{
			committer.SideAccount,
		},
		Data: OutAction{
			Committer: client.Name(string(committer.RelayAccount.Actor)),
			Num:       uint64(blockNum)*maxActionInBlock + uint64(idx),
			To:        client.Name(string(act.Memo)),
			Chain:     client.Name(string(act.Chain)),
			Contract:  client.Name(string(mapTokenStat.SideAccount)),
			Action:    client.Name(string(mapTokenStat.SideAction)),
			Quantity: client.Asset(&types.Asset{
				Amount: int64(act.Quantity.Amount),
				Symbol: types.Symbol{
					Precision: act.Quantity.Precision,
					Symbol:    act.Quantity.Symbol.Symbol,
				},
			}),
			Memo: fmt.Sprintf("from codex %s to %s in %s block %d",
				string(act.From), act.Memo, string(act.Chain), blockNum),
		},
	}

	_, err = client.PushActions(actToCommit)
	if err != nil {
		logger.Logger().Error("push action error", zap.Error(err))
	}

	return err
}
