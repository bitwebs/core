package simulation_test

import (
	"fmt"
	"testing"

	gogotypes "github.com/gogo/protobuf/types"
	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/crypto/ed25519"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/kv"

	core "github.com/bitwebs/iq-core/types"
	"github.com/bitwebs/iq-core/x/oracle/keeper"
	sim "github.com/bitwebs/iq-core/x/oracle/simulation"
	"github.com/bitwebs/iq-core/x/oracle/types"
)

var (
	delPk      = ed25519.GenPrivKey().PubKey()
	feederAddr = sdk.AccAddress(delPk.Address())
	valAddr    = sdk.ValAddress(delPk.Address())
)

func TestDecodeDistributionStore(t *testing.T) {
	cdc := keeper.MakeTestCodec(t)
	dec := sim.NewDecodeStore(cdc)

	exchangeRate := sdk.NewDecWithPrec(1234, 1)
	missCounter := uint64(23)

	aggregatePrevote := types.NewAggregateExchangeRatePrevote(types.AggregateVoteHash([]byte("12345")), valAddr, 123)
	aggregateVote := types.NewAggregateExchangeRateVote(types.ExchangeRateTuples{
		{Denom: core.MicroBKRWDenom, ExchangeRate: sdk.NewDecWithPrec(1234, 1)},
		{Denom: core.MicroBKRWDenom, ExchangeRate: sdk.NewDecWithPrec(4321, 1)},
	}, valAddr)

	tobinTax := sdk.NewDecWithPrec(2, 2)

	kvPairs := kv.Pairs{
		Pairs: []kv.Pair{
			{Key: types.ExchangeRateKey, Value: cdc.MustMarshal(&sdk.DecProto{Dec: exchangeRate})},
			{Key: types.FeederDelegationKey, Value: feederAddr.Bytes()},
			{Key: types.MissCounterKey, Value: cdc.MustMarshal(&gogotypes.UInt64Value{Value: missCounter})},
			{Key: types.AggregateExchangeRatePrevoteKey, Value: cdc.MustMarshal(&aggregatePrevote)},
			{Key: types.AggregateExchangeRateVoteKey, Value: cdc.MustMarshal(&aggregateVote)},
			{Key: types.TobinTaxKey, Value: cdc.MustMarshal(&sdk.DecProto{Dec: tobinTax})},
			{Key: []byte{0x99}, Value: []byte{0x99}},
		},
	}

	tests := []struct {
		name        string
		expectedLog string
	}{
		{"ExchangeRate", fmt.Sprintf("%v\n%v", exchangeRate, exchangeRate)},
		{"FeederDelegation", fmt.Sprintf("%v\n%v", feederAddr, feederAddr)},
		{"MissCounter", fmt.Sprintf("%v\n%v", missCounter, missCounter)},
		{"AggregatePrevote", fmt.Sprintf("%v\n%v", aggregatePrevote, aggregatePrevote)},
		{"AggregateVote", fmt.Sprintf("%v\n%v", aggregateVote, aggregateVote)},
		{"TobinTax", fmt.Sprintf("%v\n%v", tobinTax, tobinTax)},
		{"other", ""},
	}

	for i, tt := range tests {
		i, tt := i, tt
		t.Run(tt.name, func(t *testing.T) {
			switch i {
			case len(tests) - 1:
				require.Panics(t, func() { dec(kvPairs.Pairs[i], kvPairs.Pairs[i]) }, tt.name)
			default:
				require.Equal(t, tt.expectedLog, dec(kvPairs.Pairs[i], kvPairs.Pairs[i]), tt.name)
			}
		})
	}
}
