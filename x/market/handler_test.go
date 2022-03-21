package market

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	core "github.com/bitwebs/iq-core/types"
	"github.com/bitwebs/iq-core/x/market/keeper"
	"github.com/bitwebs/iq-core/x/market/types"
)

func TestMarketFilters(t *testing.T) {
	input, h := setup(t)

	// Case 1: non-oracle message being sent fails
	bankMsg := banktypes.MsgSend{}
	_, err := h(input.Ctx, &bankMsg)
	require.Error(t, err)

	// Case 2: Normal MsgSwap submission goes through
	offerCoin := sdk.NewCoin(core.MicroBiqDenom, sdk.NewInt(10))
	prevoteMsg := types.NewMsgSwap(keeper.Addrs[0], offerCoin, core.MicroBSDRDenom)
	_, err = h(input.Ctx, prevoteMsg)
	require.NoError(t, err)
}

func TestSwapMsg(t *testing.T) {
	input, h := setup(t)

	params := input.MarketKeeper.GetParams(input.Ctx)
	params.MinStabilitySpread = sdk.ZeroDec()
	input.MarketKeeper.SetParams(input.Ctx, params)

	beforeIqPoolDelta := input.MarketKeeper.GetIqPoolDelta(input.Ctx)

	amt := sdk.NewInt(10)
	offerCoin := sdk.NewCoin(core.MicroBiqDenom, amt)
	swapMsg := types.NewMsgSwap(keeper.Addrs[0], offerCoin, core.MicroBSDRDenom)
	_, err := h(input.Ctx, swapMsg)
	require.NoError(t, err)

	afterIqPoolDelta := input.MarketKeeper.GetIqPoolDelta(input.Ctx)
	diff := beforeIqPoolDelta.Sub(afterIqPoolDelta)

	// calculate estimation
	basePool := input.MarketKeeper.GetParams(input.Ctx).BasePool
	price, _ := input.OracleKeeper.GetBiqExchangeRate(input.Ctx, core.MicroBSDRDenom)
	cp := basePool.Mul(basePool)

	iqPool := basePool.Add(beforeIqPoolDelta)
	biqPool := cp.Quo(iqPool)
	estmiatedDiff := iqPool.Sub(cp.Quo(biqPool.Add(price.MulInt(amt))))
	require.True(t, estmiatedDiff.Sub(diff.Abs()).LTE(sdk.NewDecWithPrec(1, 6)))

	// invalid recursive swap
	swapMsg = types.NewMsgSwap(keeper.Addrs[0], offerCoin, core.MicroBiqDenom)

	_, err = h(input.Ctx, swapMsg)
	require.Error(t, err)

	// valid zero tobin tax test
	input.OracleKeeper.SetTobinTax(input.Ctx, core.MicroBKWRDenom, sdk.ZeroDec())
	input.OracleKeeper.SetTobinTax(input.Ctx, core.MicroBSDRDenom, sdk.ZeroDec())
	offerCoin = sdk.NewCoin(core.MicroBSDRDenom, amt)
	swapMsg = types.NewMsgSwap(keeper.Addrs[0], offerCoin, core.MicroBKWRDenom)
	_, err = h(input.Ctx, swapMsg)
	require.NoError(t, err)
}

func TestSwapSendMsg(t *testing.T) {
	input, h := setup(t)

	amt := sdk.NewInt(10)
	offerCoin := sdk.NewCoin(core.MicroBiqDenom, amt)
	retCoin, spread, err := input.MarketKeeper.ComputeSwap(input.Ctx, offerCoin, core.MicroBSDRDenom)
	require.NoError(t, err)

	expectedAmt := retCoin.Amount.Mul(sdk.OneDec().Sub(spread)).TruncateInt()

	swapSendMsg := types.NewMsgSwapSend(keeper.Addrs[0], keeper.Addrs[1], offerCoin, core.MicroBSDRDenom)
	_, err = h(input.Ctx, swapSendMsg)
	require.NoError(t, err)

	balance := input.BankKeeper.GetBalance(input.Ctx, keeper.Addrs[1], core.MicroBSDRDenom)
	require.Equal(t, expectedAmt, balance.Amount)
}
