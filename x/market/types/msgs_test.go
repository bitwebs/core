package types

import (
	"testing"

	core "github.com/bitwebs/iq-core/types"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestMsgSwap(t *testing.T) {
	addrs := []sdk.AccAddress{
		sdk.AccAddress([]byte("addr1_______________")),
	}

	overflowOfferAmt, _ := sdk.NewIntFromString("100000000000000000000000000000000000000000000000000000000")

	tests := []struct {
		trader      sdk.AccAddress
		offerCoin   sdk.Coin
		askDenom    string
		expectedErr string
	}{
		{addrs[0], sdk.NewCoin(core.MicroBiqDenom, sdk.OneInt()), core.MicroBSDRDenom, ""},
		{sdk.AccAddress{}, sdk.NewCoin(core.MicroBiqDenom, sdk.OneInt()), core.MicroBSDRDenom, "Invalid trader address (empty address string is not allowed): invalid address"},
		{addrs[0], sdk.NewCoin(core.MicroBiqDenom, sdk.ZeroInt()), core.MicroBSDRDenom, "0ubiq: invalid coins"},
		{addrs[0], sdk.NewCoin(core.MicroBiqDenom, overflowOfferAmt), core.MicroBSDRDenom, "100000000000000000000000000000000000000000000000000000000ubiq: invalid coins"},
		{addrs[0], sdk.NewCoin(core.MicroBiqDenom, sdk.OneInt()), core.MicroBiqDenom, "ubiq: recursive swap"},
	}

	for _, tc := range tests {
		msg := NewMsgSwap(tc.trader, tc.offerCoin, tc.askDenom)
		if tc.expectedErr == "" {
			require.Nil(t, msg.ValidateBasic())
		} else {
			require.EqualError(t, msg.ValidateBasic(), tc.expectedErr)
		}
	}
}

func TestMsgSwapSend(t *testing.T) {
	addrs := []sdk.AccAddress{
		sdk.AccAddress([]byte("addr1_______________")),
		sdk.AccAddress([]byte("addr2_______________")),
	}

	overflowOfferAmt, _ := sdk.NewIntFromString("100000000000000000000000000000000000000000000000000000000")

	tests := []struct {
		fromAddress sdk.AccAddress
		toAddress   sdk.AccAddress
		offerCoin   sdk.Coin
		askDenom    string
		expectedErr string
	}{
		{addrs[0], addrs[0], sdk.NewCoin(core.MicroBiqDenom, sdk.OneInt()), core.MicroBSDRDenom, ""},
		{addrs[0], sdk.AccAddress{}, sdk.NewCoin(core.MicroBiqDenom, sdk.OneInt()), core.MicroBSDRDenom, "Invalid to address (empty address string is not allowed): invalid address"},
		{sdk.AccAddress{}, addrs[0], sdk.NewCoin(core.MicroBiqDenom, sdk.OneInt()), core.MicroBSDRDenom, "Invalid from address (empty address string is not allowed): invalid address"},
		{addrs[0], addrs[0], sdk.NewCoin(core.MicroBiqDenom, sdk.ZeroInt()), core.MicroBSDRDenom, "0ubiq: invalid coins"},
		{addrs[0], addrs[0], sdk.NewCoin(core.MicroBiqDenom, overflowOfferAmt), core.MicroBSDRDenom, "100000000000000000000000000000000000000000000000000000000ubiq: invalid coins"},
		{addrs[0], addrs[0], sdk.NewCoin(core.MicroBiqDenom, sdk.OneInt()), core.MicroBiqDenom, "ubiq: recursive swap"},
	}

	for _, tc := range tests {
		msg := NewMsgSwapSend(tc.fromAddress, tc.toAddress, tc.offerCoin, tc.askDenom)
		if tc.expectedErr == "" {
			require.Nil(t, msg.ValidateBasic())
		} else {
			require.EqualError(t, msg.ValidateBasic(), tc.expectedErr)
		}
	}
}
