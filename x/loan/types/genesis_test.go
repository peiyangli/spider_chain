package types_test

import (
	"testing"

	"spider/x/loan/types"

	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	tests := []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc:     "valid genesis state",
			genState: &types.GenesisState{LoanMap: []types.Loan{{Borrower: "0"}, {Borrower: "1"}}},
			valid:    true,
		}, {
			desc: "duplicated loan",
			genState: &types.GenesisState{
				LoanMap: []types.Loan{
					{
						Borrower: "0",
					},
					{
						Borrower: "0",
					},
				},
			},
			valid: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
