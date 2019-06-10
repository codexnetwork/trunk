package side

import (
	"github.com/codexnetwork/codex-go/types"

	"github.com/codexnetwork/trunk/chainhandler"
	forceio "github.com/eosforce/goforceio"
)

type BlockToForceio struct {
	Producer         forceio.AccountName `json:"producer"`
	Num              uint32              `json:"num"`
	ID               forceio.Checksum256 `json:"id"`
	Previous         forceio.Checksum256 `json:"previous"`
	Confirmed        uint16              `json:"confirmed"`
	TransactionMRoot forceio.Checksum256 `json:"transaction_mroot"`
	ActionMRoot      forceio.Checksum256 `json:"action_mroot"`
}

func (b *BlockToForceio) FromGeneral(sw types.SwitcherInterface, bk *chainhandler.Block) {
	b.Producer = forceio.AN(bk.Producer)
	b.Num = bk.Num
	b.Confirmed = bk.Confirmed
	b.ID = forceio.Checksum256(bk.ID)
	b.Previous = forceio.Checksum256(bk.Previous)
	b.TransactionMRoot = forceio.Checksum256(bk.TransactionMRoot)
	b.ActionMRoot = forceio.Checksum256(bk.ActionMRoot)
}

type ActionToCommit struct {
	Account       interface{}               `json:"account"`
	Name          interface{}               `json:"name"`
	Authorization []PermissionLevelToCommit `json:"authorization"`
	Data          []byte                    `json:"data"`
}

type PermissionLevelToCommit struct {
	Actor      interface{} `json:"actor"`
	Permission interface{} `json:"permission"`
}

type commitParam struct {
	Name     interface{}      `json:"chain"`
	Transfer interface{}      `json:"transfer"`
	Block    BlockToForceio   `json:"block"`
	Actions  []ActionToCommit `json:"actions"`
}

func (c *commitParam) IsNeedCommit() bool {
	return len(c.Actions) > 0
}

func (c *commitParam) FromGeneral(acts *ActionsToRelay, sw types.SwitcherInterface, block *chainhandler.Block, actions []chainhandler.Action) {
	c.Block.FromGeneral(sw, block)
	c.Actions = make([]ActionToCommit, 0, len(actions))
	for _, act := range actions {
		if !acts.IsNeedCommit(act.Account, act.Name, act.Data) {
			continue
		}

		act2Commit := ActionToCommit{
			Account:       sw.NameFromCommon(act.Account),
			Name:          sw.NameFromCommon(act.Name),
			Authorization: make([]PermissionLevelToCommit, 0, len(act.Authorization)),
			Data:          act.Data,
		}
		for _, p := range act.Authorization {
			act2Commit.Authorization = append(act2Commit.Authorization, PermissionLevelToCommit{
				Actor:      sw.NameFromCommon(p.Actor),
				Permission: sw.NameFromCommon(p.Permission),
			})
		}
		c.Actions = append(c.Actions, act2Commit)
	}
}
