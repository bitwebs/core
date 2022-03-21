package wasm_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	core "github.com/bitwebs/iq-core/types"
	"github.com/bitwebs/iq-core/x/oracle/keeper"
	"github.com/bitwebs/iq-core/x/oracle/wasm"
)

func TestQueryExchangeRates(t *testing.T) {
	input := keeper.CreateTestInput(t)

	KRWExchangeRate := sdk.NewDec(1700)
	USDExchangeRate := sdk.NewDecWithPrec(17, 1)
	SDRExchangeRate := sdk.NewDecWithPrec(19, 1)
	input.OracleKeeper.SetBiqExchangeRate(input.Ctx, core.MicroBKRWDenom, KRWExchangeRate)
	input.OracleKeeper.SetBiqExchangeRate(input.Ctx, core.MicroBUSDDenom, USDExchangeRate)
	input.OracleKeeper.SetBiqExchangeRate(input.Ctx, core.MicroBSDRDenom, SDRExchangeRate)

	querier := wasm.NewWasmQuerier(input.OracleKeeper)
	var err error

	// empty data will occur error
	_, err = querier.QueryCustom(input.Ctx, []byte{})
	require.Error(t, err)

	// not existing quote denom query
	queryParams := wasm.ExchangeRateQueryParams{
		BaseDenom:   core.MicroBiqDenom,
		QuoteDenoms: []string{core.MicroBMNTDenom},
	}
	bz, err := json.Marshal(wasm.CosmosQuery{
		ExchangeRates: &queryParams,
	})
	require.NoError(t, err)

	res, err := querier.QueryCustom(input.Ctx, bz)
	require.NoError(t, err)

	var exchangeRatesResponse wasm.ExchangeRatesQueryResponse
	err = json.Unmarshal(res, &exchangeRatesResponse)
	require.NoError(t, err)
	require.Equal(t, wasm.ExchangeRatesQueryResponse{
		BaseDenom:     core.MicroBiqDenom,
		ExchangeRates: nil,
	}, exchangeRatesResponse)

	// not existing base denom query
	queryParams = wasm.ExchangeRateQueryParams{
		BaseDenom:   core.MicroBCNYDenom,
		QuoteDenoms: []string{core.MicroBKRWDenom, core.MicroBUSDDenom, core.MicroBSDRDenom},
	}
	bz, err = json.Marshal(wasm.CosmosQuery{
		ExchangeRates: &queryParams,
	})
	require.NoError(t, err)

	res, err = querier.QueryCustom(input.Ctx, bz)
	require.Error(t, err)

	// valid query biq exchange rates
	queryParams = wasm.ExchangeRateQueryParams{
		BaseDenom:   core.MicroBiqDenom,
		QuoteDenoms: []string{core.MicroBKRWDenom, core.MicroBUSDDenom, core.MicroBSDRDenom},
	}
	bz, err = json.Marshal(wasm.CosmosQuery{
		ExchangeRates: &queryParams,
	})
	require.NoError(t, err)

	res, err = querier.QueryCustom(input.Ctx, bz)
	require.NoError(t, err)

	err = json.Unmarshal(res, &exchangeRatesResponse)
	require.NoError(t, err)
	require.Equal(t, exchangeRatesResponse, wasm.ExchangeRatesQueryResponse{
		BaseDenom: core.MicroBiqDenom,
		ExchangeRates: []wasm.ExchangeRateItem{
			{
				ExchangeRate: KRWExchangeRate.String(),
				QuoteDenom:   core.MicroBKRWDenom,
			},
			{
				ExchangeRate: USDExchangeRate.String(),
				QuoteDenom:   core.MicroBUSDDenom,
			},
			{
				ExchangeRate: SDRExchangeRate.String(),
				QuoteDenom:   core.MicroBSDRDenom,
			},
		},
	})
}
