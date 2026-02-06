package official

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"spider/testutil/sample"
	officialsimulation "spider/x/official/simulation"
	"spider/x/official/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	officialGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		OperatorMap: []types.Operator{{Creator: sample.AccAddress(),
			Address: "0",
		}, {Creator: sample.AccAddress(),
			Address: "1",
		}}}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&officialGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgCreateOperator          = "op_weight_msg_official"
		defaultWeightMsgCreateOperator int = 100
	)

	var weightMsgCreateOperator int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateOperator, &weightMsgCreateOperator, nil,
		func(_ *rand.Rand) {
			weightMsgCreateOperator = defaultWeightMsgCreateOperator
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateOperator,
		officialsimulation.SimulateMsgCreateOperator(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgUpdateOperator          = "op_weight_msg_official"
		defaultWeightMsgUpdateOperator int = 100
	)

	var weightMsgUpdateOperator int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateOperator, &weightMsgUpdateOperator, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateOperator = defaultWeightMsgUpdateOperator
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateOperator,
		officialsimulation.SimulateMsgUpdateOperator(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgDeleteOperator          = "op_weight_msg_official"
		defaultWeightMsgDeleteOperator int = 100
	)

	var weightMsgDeleteOperator int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteOperator, &weightMsgDeleteOperator, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteOperator = defaultWeightMsgDeleteOperator
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteOperator,
		officialsimulation.SimulateMsgDeleteOperator(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
