package types

import "fmt"

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:      DefaultParams(),
		OperatorMap: []Operator{}}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	operatorIndexMap := make(map[string]struct{})

	for _, elem := range gs.OperatorMap {
		index := fmt.Sprint(elem.Address)
		if _, ok := operatorIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for operator")
		}
		operatorIndexMap[index] = struct{}{}
	}

	return gs.Params.Validate()
}
