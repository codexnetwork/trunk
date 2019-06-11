package blockdb

import (
	"github.com/codexnetwork/codex-go/types"
	"github.com/pkg/errors"
)

// a Just simple imp

type BlockDB struct {
	Peers map[string]*PeerBlockState
}

func (b *BlockDB) Init(peers []string) {
	b.Peers = make(map[string]*PeerBlockState, len(peers))
	for _, peer := range peers {
		n := &PeerBlockState{
			peer: peer,
		}
		n.init()
		b.Peers[peer] = n
	}
}

func (b *BlockDB) OnBlock(peer string, block *types.BlockGeneralInfo) error {
	ps, ok := b.Peers[peer]
	if !ok || ps == nil {
		return errors.New("no peer")
	}

	return ps.appendBlock(block)
}

func (b *BlockDB) DelBlockBefore(num uint32) {
	if num <= 0 {
		return
	}
	for _, stat := range b.Peers {
		stat.DelBlockBefore(num)
	}
}
