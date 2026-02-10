package loan

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	loansimulation "spider/x/loan/simulation"
	"spider/x/loan/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	loanGenesis := types.GenesisState{
		Params: types.DefaultParams(),
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&loanGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgRequestLoan          = "op_weight_msg_loan"
		defaultWeightMsgRequestLoan int = 100
	)

	var weightMsgRequestLoan int
	simState.AppParams.GetOrGenerate(opWeightMsgRequestLoan, &weightMsgRequestLoan, nil,
		func(_ *rand.Rand) {
			weightMsgRequestLoan = defaultWeightMsgRequestLoan
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgRequestLoan,
		loansimulation.SimulateMsgRequestLoan(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgApproveLoan          = "op_weight_msg_loan"
		defaultWeightMsgApproveLoan int = 100
	)

	var weightMsgApproveLoan int
	simState.AppParams.GetOrGenerate(opWeightMsgApproveLoan, &weightMsgApproveLoan, nil,
		func(_ *rand.Rand) {
			weightMsgApproveLoan = defaultWeightMsgApproveLoan
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgApproveLoan,
		loansimulation.SimulateMsgApproveLoan(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgCancelLoan          = "op_weight_msg_loan"
		defaultWeightMsgCancelLoan int = 100
	)

	var weightMsgCancelLoan int
	simState.AppParams.GetOrGenerate(opWeightMsgCancelLoan, &weightMsgCancelLoan, nil,
		func(_ *rand.Rand) {
			weightMsgCancelLoan = defaultWeightMsgCancelLoan
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCancelLoan,
		loansimulation.SimulateMsgCancelLoan(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgRepayLoan          = "op_weight_msg_loan"
		defaultWeightMsgRepayLoan int = 100
	)

	var weightMsgRepayLoan int
	simState.AppParams.GetOrGenerate(opWeightMsgRepayLoan, &weightMsgRepayLoan, nil,
		func(_ *rand.Rand) {
			weightMsgRepayLoan = defaultWeightMsgRepayLoan
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgRepayLoan,
		loansimulation.SimulateMsgRepayLoan(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgLiquidateLoan          = "op_weight_msg_loan"
		defaultWeightMsgLiquidateLoan int = 100
	)

	var weightMsgLiquidateLoan int
	simState.AppParams.GetOrGenerate(opWeightMsgLiquidateLoan, &weightMsgLiquidateLoan, nil,
		func(_ *rand.Rand) {
			weightMsgLiquidateLoan = defaultWeightMsgLiquidateLoan
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgLiquidateLoan,
		loansimulation.SimulateMsgLiquidateLoan(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
