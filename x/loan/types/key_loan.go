package types

import "cosmossdk.io/collections"

// LoanKey is the prefix to retrieve all Loan
var (
	LoanKey    = collections.NewPrefix("loan/value/")
	LoanSeqKey = collections.NewPrefix("loan/seq/")
)
