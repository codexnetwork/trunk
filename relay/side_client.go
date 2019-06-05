package relay

import (
	"time"

	"github.com/fanyang1988/force-go/types"

	"github.com/codexnetwork/trunk/logger"
	force "github.com/fanyang1988/force-go"
	"github.com/fanyang1988/force-go/config"
	"go.uber.org/zap"
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
		client, err = force.NewClient(typ, cfg)
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
