package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	core "github.com/bitwebs/iq-core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestIqPoolDeltaUpdate(t *testing.T) {
	input := CreateTestInput(t)

	iqPoolDelta := input.MarketKeeper.GetIqPoolDelta(input.Ctx)
	require.Equal(t, sdk.ZeroDec(), iqPoolDelta)

	diff := sdk.NewDec(10)
	input.MarketKeeper.SetIqPoolDelta(input.Ctx, diff)

	iqPoolDelta = input.MarketKeeper.GetIqPoolDelta(input.Ctx)
	require.Equal(t, diff, iqPoolDelta)
}

// TestReplenishPools tests that
// each pools move towards base pool
func TestReplenishPools(t *testing.T) {
	input := CreateTestInput(t)
	input.OracleKeeper.SetBiqExchangeRate(input.Ctx, core.MicroBSDRDenom, sdk.OneDec())

	basePool := input.MarketKeeper.BasePool(input.Ctx)
	iqPoolDelta := input.MarketKeeper.GetIqPoolDelta(input.Ctx)
	require.True(t, iqPoolDelta.IsZero())

	// Positive delta
	diff := basePool.QuoInt64((int64)(core.BlocksPerDay))
	input.MarketKeeper.SetIqPoolDelta(input.Ctx, diff)

	input.MarketKeeper.ReplenishPools(input.Ctx)

	iqPoolDelta = input.MarketKeeper.GetIqPoolDelta(input.Ctx)
	replenishAmt := diff.QuoInt64((int64)(input.MarketKeeper.PoolRecoveryPeriod(input.Ctx)))
	expectedDelta := diff.Sub(replenishAmt)
	require.Equal(t, expectedDelta, iqPoolDelta)

	// Negative delta
	diff = diff.Neg()
	input.MarketKeeper.SetIqPoolDelta(input.Ctx, diff)

	input.MarketKeeper.ReplenishPools(input.Ctx)

	iqPoolDelta = input.MarketKeeper.GetIqPoolDelta(input.Ctx)
	replenishAmt = diff.QuoInt64((int64)(input.MarketKeeper.PoolRecoveryPeriod(input.Ctx)))
	expectedDelta = diff.Sub(replenishAmt)
	require.Equal(t, expectedDelta, iqPoolDelta)
}
