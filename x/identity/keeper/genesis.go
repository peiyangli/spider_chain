package keeper

import (
	"context"

	"spider/x/identity/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	for _, elem := range genState.IdentityMap {
		if err := k.Identity.Set(ctx, elem.Uid, elem); err != nil {
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
	if err := k.Identity.Walk(ctx, nil, func(_ string, val types.Identity) (stop bool, err error) {
		genesis.IdentityMap = append(genesis.IdentityMap, val)
		return false, nil
	}); err != nil {
		return nil, err
	}

	return genesis, nil
}
