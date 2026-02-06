package keeper

import (
	"context"

	"spider/x/official/types"

	"cosmossdk.io/collections"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	for _, elem := range genState.OperatorMap {
		var key = collections.Join(elem.Address, elem.Module)
		if err := k.Operator.Set(ctx, key, elem); err != nil {
			return err
		}
	}

	return k.Params.Set(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis.
func (k Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	var err error

	genesis := types.DefaultGenesis()
	genesis.Params, err = k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}
	if err := k.Operator.Walk(ctx, nil, func(_ collections.Pair[string, string], val types.Operator) (stop bool, err error) {
		genesis.OperatorMap = append(genesis.OperatorMap, val)
		return false, nil
	}); err != nil {
		return nil, err
	}

	return genesis, nil
}
