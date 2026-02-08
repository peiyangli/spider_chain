package types_test

import (
	"testing"

	"spider/x/identity/types"

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
			genState: &types.GenesisState{IdentityMap: []types.Identity{{Uid: "0"}, {Uid: "1"}}},
			valid:    true,
		}, {
			desc: "duplicated identity",
			genState: &types.GenesisState{
				IdentityMap: []types.Identity{
					{
						Uid: "0",
					},
					{
						Uid: "0",
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
