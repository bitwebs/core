package v05_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitwebs/iq-core/app"
	core "github.com/bitwebs/iq-core/types"
	v04oracle "github.com/bitwebs/iq-core/x/oracle/legacy/v04"
	v05oracle "github.com/bitwebs/iq-core/x/oracle/legacy/v05"
)

func TestMigrate(t *testing.T) {
	sdk.GetConfig().SetBech32PrefixForAccount(core.Bech32PrefixAccAddr, core.Bech32PrefixAccPub)
	sdk.GetConfig().SetBech32PrefixForValidator(core.Bech32PrefixValAddr, core.Bech32PrefixValPub)

	encodingConfig := app.MakeEncodingConfig()
	clientCtx := client.Context{}.
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithCodec(encodingConfig.Marshaler)

	voter, err := sdk.ValAddressFromBech32("iqvaloper1mx72uukvzqtzhc6gde7shrjqfu5srk22v3yx7a")
	require.NoError(t, err)
	feeder, err := sdk.AccAddressFromBech32("iq1mx72uukvzqtzhc6gde7shrjqfu5srk22v7gmww")
	require.NoError(t, err)

	voter2, err := sdk.ValAddressFromBech32("iqvaloper13vs2znvhdcy948ejsh7p8p22j8l4n4y07qkhsn")
	require.NoError(t, err)
	feeder2, err := sdk.AccAddressFromBech32("iq13vs2znvhdcy948ejsh7p8p22j8l4n4y07062qq")
	require.NoError(t, err)

	voteHash, err := v04oracle.VoteHashFromHexString("24738fdea72142136dde59c1e1f79f32c53dee12")
	require.NoError(t, err)

	oracleGenState := v04oracle.GenesisState{
		AggregateExchangeRatePrevotes: []v04oracle.AggregateExchangeRatePrevote{
			{
				Hash:        voteHash,
				SubmitBlock: 100,
				Voter:       voter2,
			},
			{
				Hash:        voteHash,
				SubmitBlock: 100,
				Voter:       voter,
			},
		},
		AggregateExchangeRateVotes: []v04oracle.AggregateExchangeRateVote{
			{
				ExchangeRateTuples: v04oracle.ExchangeRateTuples{
					{
						Denom:        core.MicroBSDRDenom,
						ExchangeRate: sdk.NewDec(1800),
					},
					{
						Denom:        core.MicroBUSDDenom,
						ExchangeRate: sdk.NewDec(1700),
					},
				},
				Voter: voter2,
			},
			{
				ExchangeRateTuples: v04oracle.ExchangeRateTuples{
					{
						Denom:        core.MicroBSDRDenom,
						ExchangeRate: sdk.NewDec(1800),
					},
					{
						Denom:        core.MicroBUSDDenom,
						ExchangeRate: sdk.NewDec(1700),
					},
				},
				Voter: voter,
			},
		},
		ExchangeRates: map[string]sdk.Dec{
			core.MicroBSDRDenom: sdk.NewDec(1800),
			core.MicroBUSDDenom: sdk.NewDec(1700),
		},
		FeederDelegations: map[string]sdk.AccAddress{
			"iqvaloper13vs2znvhdcy948ejsh7p8p22j8l4n4y07qkhsn": feeder2,
			"iqvaloper1mx72uukvzqtzhc6gde7shrjqfu5srk22v3yx7a": feeder,
		},
		MissCounters: map[string]int64{
			"iqvaloper13vs2znvhdcy948ejsh7p8p22j8l4n4y07qkhsn": 321,
			"iqvaloper1mx72uukvzqtzhc6gde7shrjqfu5srk22v3yx7a": 123,
		},
		Params: v04oracle.Params{
			MinValidPerWindow:        sdk.NewDecWithPrec(5, 2),
			RewardBand:               sdk.NewDecWithPrec(7, 2),
			RewardDistributionWindow: 100,
			SlashFraction:            sdk.NewDecWithPrec(1, 3),
			SlashWindow:              100,
			VotePeriod:               100,
			VoteThreshold:            sdk.NewDecWithPrec(5, 1),
			Whitelist: v04oracle.DenomList{
				{
					Name:     core.MicroBSDRDenom,
					TobinTax: sdk.NewDecWithPrec(1, 2),
				},
				{
					Name:     core.MicroBUSDDenom,
					TobinTax: sdk.NewDecWithPrec(2, 2),
				},
			},
		},
		TobinTaxes: map[string]sdk.Dec{
			core.MicroBSDRDenom: sdk.NewDecWithPrec(1, 2),
			core.MicroBUSDDenom: sdk.NewDecWithPrec(2, 2),
		},
		ExchangeRatePrevotes: []v04oracle.ExchangeRatePrevote{},
		ExchangeRateVotes:    []v04oracle.ExchangeRateVote{},
	}

	migrated := v05oracle.Migrate(oracleGenState)

	bz, err := clientCtx.Codec.MarshalJSON(migrated)
	require.NoError(t, err)

	// Indent the JSON bz correctly.
	var jsonObj map[string]interface{}
	err = json.Unmarshal(bz, &jsonObj)
	require.NoError(t, err)
	indentedBz, err := json.MarshalIndent(jsonObj, "", "\t")
	require.NoError(t, err)

	// Make sure about:
	// ExchangeRateVotes removed
	// ExchangeRatePrevotes removed
	expected := `{
	"aggregate_exchange_rate_prevotes": [
		{
			"hash": "24738fdea72142136dde59c1e1f79f32c53dee12",
			"submit_block": "100",
			"voter": "iqvaloper13vs2znvhdcy948ejsh7p8p22j8l4n4y07qkhsn"
		},
		{
			"hash": "24738fdea72142136dde59c1e1f79f32c53dee12",
			"submit_block": "100",
			"voter": "iqvaloper1mx72uukvzqtzhc6gde7shrjqfu5srk22v3yx7a"
		}
	],
	"aggregate_exchange_rate_votes": [
		{
			"exchange_rate_tuples": [
				{
					"denom": "usdr",
					"exchange_rate": "1800.000000000000000000"
				},
				{
					"denom": "uusd",
					"exchange_rate": "1700.000000000000000000"
				}
			],
			"voter": "iqvaloper13vs2znvhdcy948ejsh7p8p22j8l4n4y07qkhsn"
		},
		{
			"exchange_rate_tuples": [
				{
					"denom": "usdr",
					"exchange_rate": "1800.000000000000000000"
				},
				{
					"denom": "uusd",
					"exchange_rate": "1700.000000000000000000"
				}
			],
			"voter": "iqvaloper1mx72uukvzqtzhc6gde7shrjqfu5srk22v3yx7a"
		}
	],
	"exchange_rates": [
		{
			"denom": "usdr",
			"exchange_rate": "1800.000000000000000000"
		},
		{
			"denom": "uusd",
			"exchange_rate": "1700.000000000000000000"
		}
	],
	"feeder_delegations": [
		{
			"feeder_address": "iq13vs2znvhdcy948ejsh7p8p22j8l4n4y07062qq",
			"validator_address": "iqvaloper13vs2znvhdcy948ejsh7p8p22j8l4n4y07qkhsn"
		},
		{
			"feeder_address": "iq1mx72uukvzqtzhc6gde7shrjqfu5srk22v7gmww",
			"validator_address": "iqvaloper1mx72uukvzqtzhc6gde7shrjqfu5srk22v3yx7a"
		}
	],
	"miss_counters": [
		{
			"miss_counter": "321",
			"validator_address": "iqvaloper13vs2znvhdcy948ejsh7p8p22j8l4n4y07qkhsn"
		},
		{
			"miss_counter": "123",
			"validator_address": "iqvaloper1mx72uukvzqtzhc6gde7shrjqfu5srk22v3yx7a"
		}
	],
	"params": {
		"min_valid_per_window": "0.050000000000000000",
		"reward_band": "0.070000000000000000",
		"reward_distribution_window": "100",
		"slash_fraction": "0.001000000000000000",
		"slash_window": "100",
		"vote_period": "100",
		"vote_threshold": "0.500000000000000000",
		"whitelist": [
			{
				"name": "usdr",
				"tobin_tax": "0.010000000000000000"
			},
			{
				"name": "uusd",
				"tobin_tax": "0.020000000000000000"
			}
		]
	},
	"tobin_taxes": [
		{
			"denom": "usdr",
			"tobin_tax": "0.010000000000000000"
		},
		{
			"denom": "uusd",
			"tobin_tax": "0.020000000000000000"
		}
	]
}`

	assert.JSONEq(t, expected, string(indentedBz))
}
