package chainhandler

import (
	"encoding/binary"

	"github.com/fanyang1988/force-go/types"
)

type Block struct {
	Producer         string            `json:"producer"`
	Num              uint32            `json:"num"`
	ID               types.Checksum256 `json:"id"`
	Previous         types.Checksum256 `json:"previous"`
	Confirmed        uint16            `json:"confirmed"`
	TransactionMRoot types.Checksum256 `json:"transaction_mroot"`
	ActionMRoot      types.Checksum256 `json:"action_mroot"`
}

func (b *Block) GetNum() uint32 {
	return BlockID2Num(b.ID)
}

func BlockID2Num(blockID []byte) uint32 {
	if len(blockID) < 32 {
		return 0
	}
	return binary.BigEndian.Uint32(blockID[:32])
}

type Action struct {
	Account       string            `json:"account"`
	Name          string            `json:"name"`
	Authorization []PermissionLevel `json:"authorization"`
	Data          []byte            `json:"data"`
}

type PermissionLevel struct {
	Actor      string `json:"actor"`
	Permission string `json:"permission"`
}
