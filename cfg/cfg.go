package cfg

import (
	"errors"

	"github.com/fanyang1988/force-go/config"
	"github.com/fanyang1988/force-go/types"

	"github.com/codexnetwork/trunk/logger"
)

// Relayer transfer, watcher and checker
type Relayer struct {
	SideAccount  types.PermissionLevel
	RelayAccount types.PermissionLevel
}

// RelayCfg cfg for relay
type RelayCfg struct {
	Chain         string `json:"chain"`
	RelayContract string `json:"relaycontract"`
}

var relayCfg RelayCfg

// ChainCfgs cfg for each chain
var chainCfgs map[string]*config.ConfigData
var chainP2PCfgs map[string][]string
var chainTyp map[string]types.ClientType

var transfers []Relayer
var watchers []Relayer

// GetChainCfg get chain cfg
func GetChainCfg(name string) (*config.ConfigData, []string) {
	c, ok := chainCfgs[name]
	if !ok || c == nil {
		panic(errors.New("no find chain cfg "))
	}
	p, ok := chainP2PCfgs[name]
	if !ok {
		panic(errors.New("no find chain p2p cfg "))
	}
	return c, p
}

// GetChainTyp get chain cfg type
func GetChainTyp(name string) types.ClientType {
	p, ok := chainTyp[name]
	if !ok {
		panic(errors.New("no find chain type cfg "))
	}
	return p
}

// GetTransfers get transfers
func GetTransfers() []Relayer {
	return transfers
}

// GetRelayCfg get cfg for relay
func GetRelayCfg() RelayCfg {
	return relayCfg
}

// GetWatchers get watchers
func GetWatchers() []Relayer {
	return watchers
}

// LoadCfgs load cfg for force-relay
func LoadCfgs(path string) error {
	cfgInFile := struct {
		Chains []struct {
			Name string            `json:"name"`
			Type string            `json:"typ"`
			P2P  []string          `json:"p2p"`
			Cfg  config.ConfigData `json:"cfg"`
		} `json:"chains"`
		Transfer []struct {
			SideAcc  string `json:"sideacc"`
			RelayAcc string `json:"relayacc"`
		} `json:"transfer"`
		Watcher []struct {
			SideAcc  string `json:"sideacc"`
			RelayAcc string `json:"relayacc"`
		} `json:"watcher"`
		Relay RelayCfg `json:"relay"`
	}{}

	err := config.LoadJSONFile(path, &cfgInFile)
	if err != nil {
		return err
	}

	chainCfgs = make(map[string]*config.ConfigData)
	for idx, c := range cfgInFile.Chains {
		chainCfgs[c.Name] = &cfgInFile.Chains[idx].Cfg
	}

	logger.Debugf("chain cfgs %v", chainCfgs)

	chainP2PCfgs = make(map[string][]string)
	for _, c := range cfgInFile.Chains {
		logger.Debugf("load p2p cfg %s -> %v", c.Name, c.P2P)
		chainP2PCfgs[c.Name] = c.P2P
	}

	logger.Debugf("chain cfgs %v", chainP2PCfgs)

	chainTyp = make(map[string]types.ClientType)
	for _, c := range cfgInFile.Chains {
		logger.Debugf("load chain cfg %s -> %v", c.Name, c.Type)
		chainTyp[c.Name] = types.String2ClientType(c.Type)
	}

	logger.Debugf("chain cfgs %v", chainTyp)

	for _, t := range cfgInFile.Transfer {
		transfers = append(transfers, Relayer{
			SideAccount: types.PermissionLevel{
				Actor:      t.SideAcc,
				Permission: "active",
			},
			RelayAccount: types.PermissionLevel{
				Actor:      t.RelayAcc,
				Permission: "active",
			},
		})
	}

	for _, t := range cfgInFile.Watcher {
		watchers = append(watchers, Relayer{
			SideAccount: types.PermissionLevel{
				Actor:      t.SideAcc,
				Permission: "active",
			},
			RelayAccount: types.PermissionLevel{
				Actor:      t.RelayAcc,
				Permission: "active",
			},
		})
	}

	relayCfg = cfgInFile.Relay

	return nil
}
