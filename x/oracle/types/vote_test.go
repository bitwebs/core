package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseExchangeRateTuples(t *testing.T) {
	valid := "123.0ubbiq,123.123ubkrw"
	_, err := ParseExchangeRateTuples(valid)
	require.NoError(t, err)

	duplicatedDenom := "100.0ubbiq,123.123ubkrw,121233.123ubkrw"
	_, err = ParseExchangeRateTuples(duplicatedDenom)
	require.Error(t, err)

	invalidCoins := "123.123"
	_, err = ParseExchangeRateTuples(invalidCoins)
	require.Error(t, err)

	invalidCoinsWithValid := "123.0ubbiq,123.1"
	_, err = ParseExchangeRateTuples(invalidCoinsWithValid)
	require.Error(t, err)

	abstainCoinsWithValid := "0.0ubbiq,123.1ubkrw"
	_, err = ParseExchangeRateTuples(abstainCoinsWithValid)
	require.NoError(t, err)
}
