package keeper_test

import (
	"testing"

	"spider/x/tokenfactory/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params:   types.DefaultParams(),
		DenomMap: []types.Denom{{Denom: "0"}, {Denom: "1"}}, NamespaceMap: []types.Namespace{{Namespace: "0"}, {Namespace: "1"}}}

	f := initFixture(t)
	err := f.keeper.InitGenesis(f.ctx, genesisState)
	require.NoError(t, err)
	got, err := f.keeper.ExportGenesis(f.ctx)
	require.NoError(t, err)
	require.NotNil(t, got)

	require.EqualExportedValues(t, genesisState.Params, got.Params)
	require.EqualExportedValues(t, genesisState.DenomMap, got.DenomMap)
	require.EqualExportedValues(t, genesisState.NamespaceMap, got.NamespaceMap)

}
