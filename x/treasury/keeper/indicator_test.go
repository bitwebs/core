package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"
	core "github.com/bitwebs/iq-core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestFeeRewardsForEpoch(t *testing.T) {
	input, _ := setupValidators(t)

	taxAmount := sdk.NewInt(1000).MulRaw(core.MicroUnit)

	// Set random prices
	input.OracleKeeper.SetBiqExchangeRate(input.Ctx, core.MicroBSDRDenom, sdk.NewDec(1))
	input.OracleKeeper.SetBiqExchangeRate(input.Ctx, core.MicroBKRWDenom, sdk.NewDec(10))
	input.OracleKeeper.SetBiqExchangeRate(input.Ctx, core.MicroBGBPDenom, sdk.NewDec(100))
	input.OracleKeeper.SetBiqExchangeRate(input.Ctx, core.MicroBCNYDenom, sdk.NewDec(1000))

	// Record tax proceeds
	input.TreasuryKeeper.RecordEpochTaxProceeds(input.Ctx, sdk.Coins{
		sdk.NewCoin(core.MicroBSDRDenom, taxAmount),
		sdk.NewCoin(core.MicroBKRWDenom, taxAmount),
		sdk.NewCoin(core.MicroBGBPDenom, taxAmount),
		sdk.NewCoin(core.MicroBCNYDenom, taxAmount),
	}.Sort())

	// Update Indicators
	input.TreasuryKeeper.UpdateIndicators(input.Ctx)

	// Get Tax Rawards (TR)
	TR := input.TreasuryKeeper.GetTR(input.Ctx, input.TreasuryKeeper.GetEpoch(input.Ctx))
	require.Equal(t, sdk.NewDec(1111).MulInt64(core.MicroUnit), TR)
}

func TestSeigniorageRewardsForEpoch(t *testing.T) {
	input, _ := setupValidators(t)

	amt := sdk.NewInt(1000)
	sdrRate := sdk.NewDec(10)

	// Add seigniorage
	input.TreasuryKeeper.RecordEpochInitialIssuance(input.Ctx)

	// Set random prices
	input.OracleKeeper.SetBiqExchangeRate(input.Ctx, core.MicroBSDRDenom, sdrRate)
	input.Ctx = input.Ctx.WithBlockHeight(int64(core.BlocksPerWeek))

	// Add seigniorage
	err := input.BankKeeper.BurnCoins(input.Ctx, faucetAccountName, sdk.NewCoins(sdk.NewCoin(core.MicroBiqDenom, amt)))
	require.NoError(t, err)

	// Update Indicators
	input.TreasuryKeeper.UpdateIndicators(input.Ctx)

	// Get seigniorage rewards (SR)
	SR := input.TreasuryKeeper.GetSR(input.Ctx, input.TreasuryKeeper.GetEpoch(input.Ctx))
	miningRewardWeight := input.TreasuryKeeper.GetRewardWeight(input.Ctx)
	require.Equal(t, sdrRate.MulInt(amt).Mul(miningRewardWeight), SR)
}

func TestMiningRewardsForEpoch(t *testing.T) {
	input, _ := setupValidators(t)

	amt := sdk.NewInt(1000).MulRaw(core.MicroUnit)
	input.TreasuryKeeper.RecordEpochInitialIssuance(input.Ctx)

	// Set random prices
	input.OracleKeeper.SetBiqExchangeRate(input.Ctx, core.MicroBSDRDenom, sdk.NewDec(1))
	input.OracleKeeper.SetBiqExchangeRate(input.Ctx, core.MicroBKRWDenom, sdk.NewDec(10))
	input.OracleKeeper.SetBiqExchangeRate(input.Ctx, core.MicroBGBPDenom, sdk.NewDec(100))
	input.OracleKeeper.SetBiqExchangeRate(input.Ctx, core.MicroBCNYDenom, sdk.NewDec(1000))

	input.Ctx = input.Ctx.WithBlockHeight(int64(core.BlocksPerWeek))

	// Record tax proceeds
	input.TreasuryKeeper.RecordEpochTaxProceeds(input.Ctx, sdk.Coins{
		sdk.NewCoin(core.MicroBSDRDenom, amt),
		sdk.NewCoin(core.MicroBKRWDenom, amt),
		sdk.NewCoin(core.MicroBGBPDenom, amt),
		sdk.NewCoin(core.MicroBCNYDenom, amt),
	}.Sort())

	// Add seigniorage
	err := input.BankKeeper.BurnCoins(input.Ctx, faucetAccountName, sdk.NewCoins(sdk.NewCoin(core.MicroBiqDenom, amt)))
	require.NoError(t, err)

	input.TreasuryKeeper.UpdateIndicators(input.Ctx)

	epoch := input.TreasuryKeeper.GetEpoch(input.Ctx)

	tProceeds := input.TreasuryKeeper.GetTR(input.Ctx, epoch)
	sProceeds := input.TreasuryKeeper.GetSR(input.Ctx, epoch)
	mProceeds := tProceeds.Add(sProceeds)

	miningRewardWeight := input.TreasuryKeeper.GetRewardWeight(input.Ctx)
	require.Equal(t, sdk.NewDec(1111).MulInt64(core.MicroUnit).Add(miningRewardWeight.MulInt(amt)), mProceeds)
}

func TestLoadIndicatorByEpoch(t *testing.T) {
	input := CreateTestInput(t)

	TRArr := []sdk.Dec{
		sdk.NewDec(100),
		sdk.NewDec(200),
		sdk.NewDec(300),
		sdk.NewDec(400),
	}

	for epoch, TR := range TRArr {
		input.TreasuryKeeper.SetTR(input.Ctx, int64(epoch), TR)
	}

	SRArr := []sdk.Dec{
		sdk.NewDec(10),
		sdk.NewDec(20),
		sdk.NewDec(30),
		sdk.NewDec(40),
	}

	for epoch, SR := range SRArr {
		input.TreasuryKeeper.SetSR(input.Ctx, int64(epoch), SR)
	}

	TSLArr := []sdk.Int{
		sdk.NewInt(1000000),
		sdk.NewInt(2000000),
		sdk.NewInt(3000000),
		sdk.NewInt(4000000),
	}

	for epoch, TSL := range TSLArr {
		input.TreasuryKeeper.SetTSL(input.Ctx, int64(epoch), TSL)
	}

	for epoch := int64(0); epoch < 4; epoch++ {
		require.Equal(t, TRArr[epoch].QuoInt(TSLArr[epoch]), TRL(input.Ctx, epoch, input.TreasuryKeeper))
		require.Equal(t, SRArr[epoch], SR(input.Ctx, epoch, input.TreasuryKeeper))
		require.Equal(t, TRArr[epoch].Add(SRArr[epoch]), MR(input.Ctx, epoch, input.TreasuryKeeper))
	}

	// empty epoch load test
	require.Equal(t, sdk.ZeroDec(), TRL(input.Ctx, 5, input.TreasuryKeeper))
	require.Equal(t, sdk.ZeroDec(), SR(input.Ctx, 5, input.TreasuryKeeper))
	require.Equal(t, sdk.ZeroDec(), MR(input.Ctx, 5, input.TreasuryKeeper))
}

func linearFn(_ sdk.Context, _ Keeper, epoch int64) sdk.Dec {
	return sdk.NewDec(epoch)
}

func TestSumIndicator(t *testing.T) {
	input := CreateTestInput(t)

	SRArr := []sdk.Dec{
		sdk.NewDec(100),
		sdk.NewDec(200),
		sdk.NewDec(300),
		sdk.NewDec(400),
		sdk.NewDec(500),
		sdk.NewDec(600),
	}

	for epoch, SR := range SRArr {
		input.TreasuryKeeper.SetSR(input.Ctx, int64(epoch), SR)
	}

	// Case 1: at epoch 0 and summing over 0 epochs
	rval := input.TreasuryKeeper.sumIndicator(input.Ctx, 0, SR)
	require.Equal(t, sdk.ZeroDec(), rval)

	// Case 2: at epoch 0 and summing over negative epochs
	rval = input.TreasuryKeeper.sumIndicator(input.Ctx, -1, SR)
	require.Equal(t, sdk.ZeroDec(), rval)

	// Case 3: at epoch 3 and summing over 3, 4, 5 epochs; all should have the same rval
	input.Ctx = input.Ctx.WithBlockHeight(int64(core.BlocksPerWeek * 3))
	rval = input.TreasuryKeeper.sumIndicator(input.Ctx, 4, SR)
	rval2 := input.TreasuryKeeper.sumIndicator(input.Ctx, 5, SR)
	rval3 := input.TreasuryKeeper.sumIndicator(input.Ctx, 6, SR)
	require.Equal(t, sdk.NewDec(1000), rval)
	require.Equal(t, rval, rval2)
	require.Equal(t, rval2, rval3)

	// Case 4: at epoch 3 and summing over 0 epochs
	rval = input.TreasuryKeeper.sumIndicator(input.Ctx, 0, SR)
	require.Equal(t, sdk.ZeroDec(), rval)

	// Case 5. Sum up to 6
	input.Ctx = input.Ctx.WithBlockHeight(int64(core.BlocksPerWeek * 5))
	rval = input.TreasuryKeeper.sumIndicator(input.Ctx, 6, SR)
	require.Equal(t, sdk.NewDec(2100), rval)
}

func TestRollingAverageIndicator(t *testing.T) {
	input := CreateTestInput(t)
	SRArr := []sdk.Dec{
		sdk.NewDec(100),
		sdk.NewDec(200),
		sdk.NewDec(300),
		sdk.NewDec(400),
	}

	for epoch, SR := range SRArr {
		input.TreasuryKeeper.SetSR(input.Ctx, int64(epoch), SR)
	}

	// Case 1: at epoch 0 and averaging over 0 epochs
	rval := input.TreasuryKeeper.rollingAverageIndicator(input.Ctx, 0, SR)
	require.Equal(t, sdk.ZeroDec(), rval)

	// Case 2: at epoch 0 and averaging over negative epochs
	rval = input.TreasuryKeeper.rollingAverageIndicator(input.Ctx, -1, SR)
	require.Equal(t, sdk.ZeroDec(), rval)

	// Case 3: at epoch 3 and averaging over 3, 4, 5 epochs; all should have the same rval
	input.Ctx = input.Ctx.WithBlockHeight(int64(core.BlocksPerWeek * 3))
	rval = input.TreasuryKeeper.rollingAverageIndicator(input.Ctx, 4, SR)
	rval2 := input.TreasuryKeeper.rollingAverageIndicator(input.Ctx, 5, SR)
	rval3 := input.TreasuryKeeper.rollingAverageIndicator(input.Ctx, 6, SR)
	require.Equal(t, sdk.NewDec(250), rval)
	require.Equal(t, rval, rval2)
	require.Equal(t, rval2, rval3)

	// Case 4: at epoch 3 and averaging over 0 epochs
	rval = input.TreasuryKeeper.rollingAverageIndicator(input.Ctx, 0, SR)
	require.Equal(t, sdk.ZeroDec(), rval)

	// Case 5: at epoch 3 and averaging over 1 epoch
	rval = input.TreasuryKeeper.rollingAverageIndicator(input.Ctx, 1, SR)
	require.Equal(t, sdk.NewDec(400), rval)
}
