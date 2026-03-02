package keeper

import (
	"context"

	"spider/x/snft/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	for _, elem := range genState.ClassOwnerMap {
		if err := k.ClassOwner.Set(ctx, elem.ClassId, elem); err != nil {
			return err
		}
	}
	for _, elem := range genState.ClassNamespaceMap {
		if err := k.ClassNamespace.Set(ctx, elem.Namespace, elem); err != nil {
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
	if err := k.ClassOwner.Walk(ctx, nil, func(_ string, val types.ClassOwner) (stop bool, err error) {
		genesis.ClassOwnerMap = append(genesis.ClassOwnerMap, val)
		return false, nil
	}); err != nil {
		return nil, err
	}
	if err := k.ClassNamespace.Walk(ctx, nil, func(_ string, val types.ClassNamespace) (stop bool, err error) {
		genesis.ClassNamespaceMap = append(genesis.ClassNamespaceMap, val)
		return false, nil
	}); err != nil {
		return nil, err
	}

	return genesis, nil
}
