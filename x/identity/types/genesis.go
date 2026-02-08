package types

import "fmt"

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:      DefaultParams(),
		IdentityMap: []Identity{}}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	identityIndexMap := make(map[string]struct{})

	for _, elem := range gs.IdentityMap {
		index := fmt.Sprint(elem.Uid)
		if _, ok := identityIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for identity")
		}
		identityIndexMap[index] = struct{}{}
	}

	return gs.Params.Validate()
}
