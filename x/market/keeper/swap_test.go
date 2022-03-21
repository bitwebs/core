package keeper

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	core "github.com/bitwebs/iq-core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestApplySwapToPool(t *testing.T) {
	input := CreateTestInput(t)

	biqPriceInSDR := sdk.NewDecWithPrec(17, 1)
	input.OracleKeeper.SetBiqExchangeRate(input.Ctx, core.MicroBSDRDenom, biqPriceInSDR)

	offerCoin := sdk.NewCoin(core.MicroBiqDenom, sdk.NewInt(1000))
	askCoin := sdk.NewDecCoin(core.MicroBSDRDenom, sdk.NewInt(1700))
	oldSDRPoolDelta := input.MarketKeeper.GetIqPoolDelta(input.Ctx)
	input.MarketKeeper.ApplySwapToPool(input.Ctx, offerCoin, askCoin)
	newSDRPoolDelta := input.MarketKeeper.GetIqPoolDelta(input.Ctx)
	sdrDiff := newSDRPoolDelta.Sub(oldSDRPoolDelta)
	require.Equal(t, sdk.NewDec(-1700), sdrDiff)

	// reverse swap
	offerCoin = sdk.NewCoin(core.MicroBSDRDenom, sdk.NewInt(1700))
	askCoin = sdk.NewDecCoin(core.MicroBiqDenom, sdk.NewInt(1000))
	oldSDRPoolDelta = input.MarketKeeper.GetIqPoolDelta(input.Ctx)
	input.MarketKeeper.ApplySwapToPool(input.Ctx, offerCoin, askCoin)
	newSDRPoolDelta = input.MarketKeeper.GetIqPoolDelta(input.Ctx)
	sdrDiff = newSDRPoolDelta.Sub(oldSDRPoolDelta)
	require.Equal(t, sdk.NewDec(1700), sdrDiff)

	// IQ <> IQ, no pool changes are expected
	offerCoin = sdk.NewCoin(core.MicroBSDRDenom, sdk.NewInt(1700))
	askCoin = sdk.NewDecCoin(core.MicroBKRWDenom, sdk.NewInt(3400))
	oldSDRPoolDelta = input.MarketKeeper.GetIqPoolDelta(input.Ctx)
	input.MarketKeeper.ApplySwapToPool(input.Ctx, offerCoin, askCoin)
	newSDRPoolDelta = input.MarketKeeper.GetIqPoolDelta(input.Ctx)
	sdrDiff = newSDRPoolDelta.Sub(oldSDRPoolDelta)
	require.Equal(t, sdk.NewDec(0), sdrDiff)
}
func TestComputeSwap(t *testing.T) {
	input := CreateTestInput(t)

	// Set Oracle Price
	biqPriceInSDR := sdk.NewDecWithPrec(17, 1)
	input.OracleKeeper.SetBiqExchangeRate(input.Ctx, core.MicroBSDRDenom, biqPriceInSDR)

	for i := 0; i < 100; i++ {
		swapAmountInSDR := biqPriceInSDR.MulInt64(rand.Int63()%10000 + 2).TruncateInt()
		offerCoin := sdk.NewCoin(core.MicroBSDRDenom, swapAmountInSDR)
		retCoin, spread, err := input.MarketKeeper.ComputeSwap(input.Ctx, offerCoin, core.MicroBiqDenom)

		require.NoError(t, err)
		require.True(t, spread.GTE(input.MarketKeeper.MinStabilitySpread(input.Ctx)))
		require.Equal(t, sdk.NewDecFromInt(offerCoin.Amount).Quo(biqPriceInSDR), retCoin.Amount)
	}

	offerCoin := sdk.NewCoin(core.MicroBSDRDenom, biqPriceInSDR.QuoInt64(2).TruncateInt())
	_, _, err := input.MarketKeeper.ComputeSwap(input.Ctx, offerCoin, core.MicroBiqDenom)
	require.Error(t, err)
}

func TestComputeInternalSwap(t *testing.T) {
	input := CreateTestInput(t)

	// Set Oracle Price
	biqPriceInSDR := sdk.NewDecWithPrec(17, 1)
	input.OracleKeeper.SetBiqExchangeRate(input.Ctx, core.MicroBSDRDenom, biqPriceInSDR)

	for i := 0; i < 100; i++ {
		offerCoin := sdk.NewDecCoin(core.MicroBSDRDenom, biqPriceInSDR.MulInt64(rand.Int63()+1).TruncateInt())
		retCoin, err := input.MarketKeeper.ComputeInternalSwap(input.Ctx, offerCoin, core.MicroBiqDenom)
		require.NoError(t, err)
		require.Equal(t, offerCoin.Amount.Quo(biqPriceInSDR), retCoin.Amount)
	}

	offerCoin := sdk.NewDecCoin(core.MicroBSDRDenom, biqPriceInSDR.QuoInt64(2).TruncateInt())
	_, err := input.MarketKeeper.ComputeInternalSwap(input.Ctx, offerCoin, core.MicroBiqDenom)
	require.Error(t, err)
}

func TestIlliquidTobinTaxListParams(t *testing.T) {
	input := CreateTestInput(t)

	// Set Oracle Price
	biqPriceInSDR := sdk.NewDecWithPrec(17, 1)
	biqPriceInMNT := sdk.NewDecWithPrec(7652, 1)
	input.OracleKeeper.SetBiqExchangeRate(input.Ctx, core.MicroBSDRDenom, biqPriceInSDR)
	input.OracleKeeper.SetBiqExchangeRate(input.Ctx, core.MicroBMNTDenom, biqPriceInMNT)

	tobinTax := sdk.NewDecWithPrec(25, 4)
	params := input.MarketKeeper.GetParams(input.Ctx)
	input.MarketKeeper.SetParams(input.Ctx, params)

	illiquidFactor := sdk.NewDec(2)
	input.OracleKeeper.SetTobinTax(input.Ctx, core.MicroBSDRDenom, tobinTax)
	input.OracleKeeper.SetTobinTax(input.Ctx, core.MicroBMNTDenom, tobinTax.Mul(illiquidFactor))

	swapAmountInSDR := biqPriceInSDR.MulInt64(rand.Int63()%10000 + 2).TruncateInt()
	offerCoin := sdk.NewCoin(core.MicroBSDRDenom, swapAmountInSDR)
	_, spread, err := input.MarketKeeper.ComputeSwap(input.Ctx, offerCoin, core.MicroBMNTDenom)
	require.NoError(t, err)
	require.Equal(t, tobinTax.Mul(illiquidFactor), spread)
}
