package chaindata

import (
	eos "github.com/eosforce/goforceio"
	forceio "github.com/fanyang1988/force-go/forceio"
	"github.com/fanyang1988/force-go/types"
	"github.com/pkg/errors"
)

// MapTokenStat map token stat in relay.token
type MapTokenStat struct {
	Supply                   eos.Asset       `json:"supply"`
	MaxSupply                eos.Asset       `json:"max_supply"`
	Issuer                   eos.AccountName `json:"issuer"`
	Chain                    eos.Name        `json:"chain"`
	SideAccount              eos.AccountName `json:"side_account"`
	SideAction               eos.ActionName  `json:"side_action"`
	RewardsPool              eos.Asset       `json:"rewards_pool"`
	TotalMineage             eos.Int128      `json:"total_mineage"`
	TotalMineageUpdateHeight uint32          `json:"total_mineage_update_height"`
	TotalPendingMineage      eos.Int64       `json:"total_pending_mineage"`
}

// GetTokenMapDatas get token map stat data from table
func GetTokenMapDatas(chain string, client types.ClientInterface) ([]MapTokenStat, error) {
	// just relay chain can get this data
	cli, ok := client.(*forceio.API)
	if !ok {
		return nil, errors.New("only relay chain can GetTokenMapData")
	}

	req := eos.GetTableRowsRequest{
		Code:  relayTokenAcc,
		Scope: chain,
		Table: "stat",
	}

	res, err := cli.GetTableRows(req)
	if err != nil {
		return nil, errors.Wrapf(err, "get table")
	}

	rspData := make([]MapTokenStat, 0, 32)
	err = res.BinaryToStructs(&rspData)
	if err != nil {
		return nil, errors.Wrapf(err, "to struct")
	}

	return rspData, nil
}

// GetTokenMapData get token map stat data from table
func GetTokenMapData(chain string, symbol string, client types.ClientInterface) (MapTokenStat, error) {
	// just relay chain can get this data
	cli, ok := client.(*forceio.API)
	if !ok {
		return MapTokenStat{}, errors.New("only relay chain can GetTokenMapData")
	}

	req := eos.GetTableRowsRequest{
		Code:  relayTokenAcc,
		Scope: chain,
		Table: "stat",
	}

	res, err := cli.GetTableRows(req)
	if err != nil {
		return MapTokenStat{}, errors.Wrapf(err, "get table")
	}

	rspData := make([]MapTokenStat, 0, 32)
	err = res.BinaryToStructs(&rspData)
	if err != nil {
		return MapTokenStat{}, errors.Wrapf(err, "to struct")
	}

	for _, d := range rspData {
		if (d.MaxSupply.Symbol.Symbol == symbol) && (string(d.Chain) == chain) {
			return d, nil
		}
	}

	return MapTokenStat{}, errors.New("no found map token stat")
}
