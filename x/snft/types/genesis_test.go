package types_test

import (
	"testing"

	"spider/x/snft/types"

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
			genState: &types.GenesisState{ClassOwnerMap: []types.ClassOwner{{ClassId: "0"}, {ClassId: "1"}}, ClassNamespaceMap: []types.ClassNamespace{{Namespace: "0"}, {Namespace: "1"}}},
			valid:    true,
		}, {
			desc: "duplicated classOwner",
			genState: &types.GenesisState{
				ClassOwnerMap: []types.ClassOwner{
					{
						ClassId: "0",
					},
					{
						ClassId: "0",
					},
				},
				ClassNamespaceMap: []types.ClassNamespace{{Namespace: "0"}, {Namespace: "1"}}},
			valid: false,
		}, {
			desc: "duplicated classNamespace",
			genState: &types.GenesisState{
				ClassNamespaceMap: []types.ClassNamespace{
					{
						Namespace: "0",
					},
					{
						Namespace: "0",
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
