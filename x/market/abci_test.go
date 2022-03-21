package market

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/bitwebs/iq-core/x/market/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestReplenishPools(t *testing.T) {
	input := keeper.CreateTestInput(t)

	iqDelta := sdk.NewDecWithPrec(17987573223725367, 3)
	input.MarketKeeper.SetIqPoolDelta(input.Ctx, iqDelta)

	for i := 0; i < 100; i++ {
		iqDelta = input.MarketKeeper.GetIqPoolDelta(input.Ctx)

		poolRecoveryPeriod := int64(input.MarketKeeper.PoolRecoveryPeriod(input.Ctx))
		iqRegressionAmt := iqDelta.QuoInt64(poolRecoveryPeriod)

		EndBlocker(input.Ctx, input.MarketKeeper)

		iqPoolDelta := input.MarketKeeper.GetIqPoolDelta(input.Ctx)
		require.Equal(t, iqDelta.Sub(iqRegressionAmt), iqPoolDelta)
	}
}
