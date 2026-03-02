package types

import "fmt"

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:        DefaultParams(),
		ClassOwnerMap: []ClassOwner{}, ClassNamespaceMap: []ClassNamespace{}}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	classOwnerIndexMap := make(map[string]struct{})

	for _, elem := range gs.ClassOwnerMap {
		index := fmt.Sprint(elem.ClassId)
		if _, ok := classOwnerIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for classOwner")
		}
		classOwnerIndexMap[index] = struct{}{}
	}
	classNamespaceIndexMap := make(map[string]struct{})

	for _, elem := range gs.ClassNamespaceMap {
		index := fmt.Sprint(elem.Namespace)
		if _, ok := classNamespaceIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for classNamespace")
		}
		classNamespaceIndexMap[index] = struct{}{}
	}

	return gs.Params.Validate()
}
