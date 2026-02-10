package keeper

import (
	"context"

	"spider/x/loan/types"

	"cosmossdk.io/collections"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	for _, elem := range genState.LoanMap {
		if err := k.Loan.Set(ctx, collections.Join3(elem.Borrower, elem.Status, elem.Seq), elem); err != nil {
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
	if err := k.Loan.Walk(ctx, nil, func(_ collections.Triple[string, uint64, uint64], val types.Loan) (stop bool, err error) {
		genesis.LoanMap = append(genesis.LoanMap, val)
		return false, nil
	}); err != nil {
		return nil, err
	}

	return genesis, nil
}
