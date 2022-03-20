package cli

import (
	"context"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitwebs/iq-core/x/market/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	marketQueryCmd := &cobra.Command{
		Use:                        "market",
		Short:                      "Querying commands for the market module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	marketQueryCmd.AddCommand(
		GetCmdQuerySwap(),
		GetCmdQueryIqPoolDelta(),
		GetCmdQueryParams(),
	)

	return marketQueryCmd
}

// GetCmdQuerySwap implements the query swap simulation result command.
func GetCmdQuerySwap() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "swap [offer-coin] [ask-denom]",
		Args:  cobra.ExactArgs(2),
		Short: "Query a quote for a swap operation",
		Long: strings.TrimSpace(`
Query a quote for how many coins can be received in a swap operation. Note; rates are dynamic and can quickly change.

$ iqd query swap 5000000ubiq ubusd
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			// parse offerCoin
			offerCoinStr := args[0]
			_, err = sdk.ParseCoinNormalized(offerCoinStr)
			if err != nil {
				return err
			}

			askDenom := args[1]

			res, err := queryClient.Swap(context.Background(),
				&types.QuerySwapRequest{OfferCoin: offerCoinStr, AskDenom: askDenom},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryIqPoolDelta implements the query mint pool delta command.
func GetCmdQueryIqPoolDelta() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "iq-pool-delta",
		Args:  cobra.NoArgs,
		Short: "Query iq pool delta",
		Long: `Query iq pool delta, which is ubusd amount used for swap operation from the IqPool.
It can be negative if the market wants more Iq than Biq, and vice versa if the market wants more Biq.

$ iqd query market iq-pool-delta
	`,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.IqPoolDelta(context.Background(),
				&types.QueryIqPoolDeltaRequest{},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryParams implements the query params command.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current market params",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Params(context.Background(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
