package market

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitwebs/iq-core/x/market/keeper"
	"github.com/bitwebs/iq-core/x/market/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data *types.GenesisState) {
	keeper.SetParams(ctx, data.Params)
	keeper.SetIqPoolDelta(ctx, data.IqPoolDelta)

	// check if the module account exists
	moduleAcc := keeper.GetMarketAccount(ctx)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) (data *types.GenesisState) {
	params := keeper.GetParams(ctx)
	iqPoolDelta := keeper.GetIqPoolDelta(ctx)

	return types.NewGenesisState(iqPoolDelta, params)
}
