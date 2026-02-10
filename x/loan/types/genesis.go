package types

import "fmt"

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:  DefaultParams(),
		LoanMap: []Loan{}}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	loanIndexMap := make(map[string]struct{})

	for _, elem := range gs.LoanMap {
		index := fmt.Sprint(elem.Borrower)
		if _, ok := loanIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for loan")
		}
		loanIndexMap[index] = struct{}{}
	}

	return gs.Params.Validate()
}
