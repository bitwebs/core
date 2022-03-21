package market

import (
	"testing"

	core "github.com/bitwebs/iq-core/types"
	"github.com/bitwebs/iq-core/x/market/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	uSDRAmt    = sdk.NewInt(1005 * core.MicroUnit)
	stakingAmt = sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction)

	randomPrice = sdk.NewDec(1700)
)

func setup(t *testing.T) (keeper.TestInput, sdk.Handler) {
	input := keeper.CreateTestInput(t)

	params := input.MarketKeeper.GetParams(input.Ctx)
	input.MarketKeeper.SetParams(input.Ctx, params)
	input.OracleKeeper.SetBiqExchangeRate(input.Ctx, core.MicroBSDRDenom, randomPrice)
	input.OracleKeeper.SetBiqExchangeRate(input.Ctx, core.MicroBKRWDenom, randomPrice)
	h := NewHandler(input.MarketKeeper)

	return input, h
}
