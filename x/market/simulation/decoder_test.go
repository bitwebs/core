package simulation

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/bitwebs/iq-core/x/market/keeper"
	"github.com/bitwebs/iq-core/x/market/types"
)

func TestDecodeDistributionStore(t *testing.T) {
	cdc := keeper.MakeTestCodec(t)
	dec := NewDecodeStore(cdc)

	iqDelta := sdk.NewDecWithPrec(12, 2)

	kvPairs := kv.Pairs{
		Pairs: []kv.Pair{
			{Key: types.IqPoolDeltaKey, Value: cdc.MustMarshal(&sdk.DecProto{Dec: iqDelta})},
			{Key: []byte{0x99}, Value: []byte{0x99}},
		},
	}

	tests := []struct {
		name        string
		expectedLog string
	}{
		{"IqPoolDelta", fmt.Sprintf("%v\n%v", iqDelta, iqDelta)},
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
