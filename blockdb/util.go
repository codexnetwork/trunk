package blockdb

import (
	"github.com/codexnetwork/codex-go/types"
)

func IsChecksum256Eq(l, r types.Checksum256) bool {
	if len(l) != len(r) {
		return false
	}

	if len(l) == 0 {
		return true
	}

	for i := 0; i < len(l); i++ {
		if l[i] != r[i] {
			return false
		}
	}

	return true
}

// IsBlockEq return if block l == r
func IsBlockEq(l, r *types.BlockGeneralInfo) bool {
	return IsChecksum256Eq(l.ID, l.ID) && IsChecksum256Eq(l.Previous, r.Previous)
}
