package keeper_test

import (
	"testing"

	"spider/x/snft/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params:        types.DefaultParams(),
		ClassOwnerMap: []types.ClassOwner{{ClassId: "0"}, {ClassId: "1"}}, ClassNamespaceMap: []types.ClassNamespace{{Namespace: "0"}, {Namespace: "1"}}}

	f := initFixture(t)
	err := f.keeper.InitGenesis(f.ctx, genesisState)
	require.NoError(t, err)
	got, err := f.keeper.ExportGenesis(f.ctx)
	require.NoError(t, err)
	require.NotNil(t, got)

	require.EqualExportedValues(t, genesisState.Params, got.Params)
	require.EqualExportedValues(t, genesisState.ClassOwnerMap, got.ClassOwnerMap)
	require.EqualExportedValues(t, genesisState.ClassNamespaceMap, got.ClassNamespaceMap)

}
