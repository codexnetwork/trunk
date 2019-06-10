package relay

import (
	"time"

	"go.uber.org/zap"

	"github.com/codexnetwork/codex-go/types"
	gocodex "github.com/codexnetwork/codex-go"
	"github.com/codexnetwork/codex-go/config"
	"github.com/codexnetwork/trunk/logger"
)

// client client to force relay chain
var client types.ClientInterface

// CreateSideClient create client to force side chain
func CreateSideClient(typ types.ClientType, cfg *config.ConfigData) {
	for {
		var err error
		logger.Logger().Info("create client cfg",
			zap.String("url", cfg.URL),
			zap.String("chainID", cfg.ChainID))
		client, err = gocodex.NewClient(typ, cfg)
		if err != nil {
			logger.LogError("create client error, need retry", err)
			time.Sleep(1 * time.Second)
		} else {
			return
		}
	}
}

func Client() types.ClientInterface {
	return client
}
