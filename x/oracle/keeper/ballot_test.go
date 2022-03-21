package keeper

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	core "github.com/bitwebs/iq-core/types"
	"github.com/bitwebs/iq-core/x/oracle/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

func TestOrganizeAggregate(t *testing.T) {
	input := CreateTestInput(t)

	power := int64(100)
	amt := sdk.TokensFromConsensusPower(power, sdk.DefaultPowerReduction)
	sh := staking.NewHandler(input.StakingKeeper)
	ctx := input.Ctx

	// Validator created
	_, err := sh(ctx, NewTestMsgCreateValidator(ValAddrs[0], ValPubKeys[0], amt))
	require.NoError(t, err)
	_, err = sh(ctx, NewTestMsgCreateValidator(ValAddrs[1], ValPubKeys[1], amt))
	require.NoError(t, err)
	_, err = sh(ctx, NewTestMsgCreateValidator(ValAddrs[2], ValPubKeys[2], amt))
	require.NoError(t, err)
	staking.EndBlocker(ctx, input.StakingKeeper)

	sdrBallot := types.ExchangeRateBallot{
		types.NewVoteForTally(sdk.NewDec(17), core.MicroBSDRDenom, ValAddrs[0], power),
		types.NewVoteForTally(sdk.NewDec(10), core.MicroBSDRDenom, ValAddrs[1], power),
		types.NewVoteForTally(sdk.NewDec(6), core.MicroBSDRDenom, ValAddrs[2], power),
	}
	krwBallot := types.ExchangeRateBallot{
		types.NewVoteForTally(sdk.NewDec(1000), core.MicroBKRWDenom, ValAddrs[0], power),
		types.NewVoteForTally(sdk.NewDec(1300), core.MicroBKRWDenom, ValAddrs[1], power),
		types.NewVoteForTally(sdk.NewDec(2000), core.MicroBKRWDenom, ValAddrs[2], power),
	}

	for i := range sdrBallot {
		input.OracleKeeper.SetAggregateExchangeRateVote(input.Ctx, ValAddrs[i],
			types.NewAggregateExchangeRateVote(types.ExchangeRateTuples{
				{Denom: sdrBallot[i].Denom, ExchangeRate: sdrBallot[i].ExchangeRate},
				{Denom: krwBallot[i].Denom, ExchangeRate: krwBallot[i].ExchangeRate},
			}, ValAddrs[i]))
	}

	// organize votes by denom
	ballotMap := input.OracleKeeper.OrganizeBallotByDenom(input.Ctx, map[string]types.Claim{
		ValAddrs[0].String(): {
			Power:     power,
			WinCount:  0,
			Recipient: ValAddrs[0],
		},
		ValAddrs[1].String(): {
			Power:     power,
			WinCount:  0,
			Recipient: ValAddrs[1],
		},
		ValAddrs[2].String(): {
			Power:     power,
			WinCount:  0,
			Recipient: ValAddrs[2],
		},
	})

	// sort each ballot for comparison
	sort.Sort(sdrBallot)
	sort.Sort(krwBallot)
	sort.Sort(ballotMap[core.MicroBSDRDenom])
	sort.Sort(ballotMap[core.MicroBKRWDenom])

	require.Equal(t, sdrBallot, ballotMap[core.MicroBSDRDenom])
	require.Equal(t, krwBallot, ballotMap[core.MicroBKRWDenom])
}

func TestClearBallots(t *testing.T) {
	input := CreateTestInput(t)

	power := int64(100)
	amt := sdk.TokensFromConsensusPower(power, sdk.DefaultPowerReduction)
	sh := staking.NewHandler(input.StakingKeeper)
	ctx := input.Ctx

	// Validator created
	_, err := sh(ctx, NewTestMsgCreateValidator(ValAddrs[0], ValPubKeys[0], amt))
	require.NoError(t, err)
	_, err = sh(ctx, NewTestMsgCreateValidator(ValAddrs[1], ValPubKeys[1], amt))
	require.NoError(t, err)
	_, err = sh(ctx, NewTestMsgCreateValidator(ValAddrs[2], ValPubKeys[2], amt))
	require.NoError(t, err)
	staking.EndBlocker(ctx, input.StakingKeeper)

	sdrBallot := types.ExchangeRateBallot{
		types.NewVoteForTally(sdk.NewDec(17), core.MicroBSDRDenom, ValAddrs[0], power),
		types.NewVoteForTally(sdk.NewDec(10), core.MicroBSDRDenom, ValAddrs[1], power),
		types.NewVoteForTally(sdk.NewDec(6), core.MicroBSDRDenom, ValAddrs[2], power),
	}
	krwBallot := types.ExchangeRateBallot{
		types.NewVoteForTally(sdk.NewDec(1000), core.MicroBKRWDenom, ValAddrs[0], power),
		types.NewVoteForTally(sdk.NewDec(1300), core.MicroBKRWDenom, ValAddrs[1], power),
		types.NewVoteForTally(sdk.NewDec(2000), core.MicroBKRWDenom, ValAddrs[2], power),
	}

	for i := range sdrBallot {
		input.OracleKeeper.SetAggregateExchangeRatePrevote(input.Ctx, ValAddrs[i], types.AggregateExchangeRatePrevote{
			Hash:        "",
			Voter:       ValAddrs[i].String(),
			SubmitBlock: uint64(input.Ctx.BlockHeight()),
		})

		input.OracleKeeper.SetAggregateExchangeRateVote(input.Ctx, ValAddrs[i],
			types.NewAggregateExchangeRateVote(types.ExchangeRateTuples{
				{Denom: sdrBallot[i].Denom, ExchangeRate: sdrBallot[i].ExchangeRate},
				{Denom: krwBallot[i].Denom, ExchangeRate: krwBallot[i].ExchangeRate},
			}, ValAddrs[i]))
	}

	input.OracleKeeper.ClearBallots(input.Ctx, 5)

	prevoteCounter := 0
	voteCounter := 0
	input.OracleKeeper.IterateAggregateExchangeRatePrevotes(input.Ctx, func(_ sdk.ValAddress, _ types.AggregateExchangeRatePrevote) bool {
		prevoteCounter++
		return false
	})
	input.OracleKeeper.IterateAggregateExchangeRateVotes(input.Ctx, func(_ sdk.ValAddress, _ types.AggregateExchangeRateVote) bool {
		voteCounter++
		return false
	})

	require.Equal(t, prevoteCounter, 3)
	require.Equal(t, voteCounter, 0)

	input.OracleKeeper.ClearBallots(input.Ctx.WithBlockHeight(input.Ctx.BlockHeight()+6), 5)

	prevoteCounter = 0
	input.OracleKeeper.IterateAggregateExchangeRatePrevotes(input.Ctx, func(_ sdk.ValAddress, _ types.AggregateExchangeRatePrevote) bool {
		prevoteCounter++
		return false
	})
	require.Equal(t, prevoteCounter, 0)
}

func TestApplyWhitelist(t *testing.T) {
	input := CreateTestInput(t)

	// no update
	input.OracleKeeper.ApplyWhitelist(input.Ctx, types.DenomList{
		types.Denom{
			Name:     "ubusd",
			TobinTax: sdk.OneDec(),
		},
		types.Denom{
			Name:     "ubkrw",
			TobinTax: sdk.OneDec(),
		},
	}, map[string]sdk.Dec{
		"ubusd": sdk.ZeroDec(),
		"ubkrw": sdk.ZeroDec(),
	})

	price, err := input.OracleKeeper.GetTobinTax(input.Ctx, "ubusd")
	require.NoError(t, err)
	require.Equal(t, price, sdk.OneDec())

	price, err = input.OracleKeeper.GetTobinTax(input.Ctx, "ubkrw")
	require.NoError(t, err)
	require.Equal(t, price, sdk.OneDec())

	metadata, ok := input.BankKeeper.GetDenomMetaData(input.Ctx, "ubusd")
	require.True(t, ok)
	require.Equal(t, metadata.Base, "ubusd")
	require.Equal(t, metadata.Display, "busd")
	require.Equal(t, len(metadata.DenomUnits), 3)
	require.Equal(t, metadata.Description, "The native stable token of the IQ Swartz.")

	metadata, ok = input.BankKeeper.GetDenomMetaData(input.Ctx, "ubkrw")
	require.True(t, ok)
	require.Equal(t, metadata.Base, "ubkrw")
	require.Equal(t, metadata.Display, "bkrw")
	require.Equal(t, len(metadata.DenomUnits), 3)
	require.Equal(t, metadata.Description, "The native stable token of the IQ Swartz.")
}
