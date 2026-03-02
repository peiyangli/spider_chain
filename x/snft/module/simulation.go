package snft

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"spider/testutil/sample"
	snftsimulation "spider/x/snft/simulation"
	"spider/x/snft/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	snftGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		ClassOwnerMap: []types.ClassOwner{{Owner: sample.AccAddress(),
			ClassId: "0",
		}, {Owner: sample.AccAddress(),
			ClassId: "1",
		}}, ClassNamespaceMap: []types.ClassNamespace{{Creator: sample.AccAddress(),
			Namespace: "0",
		}, {Creator: sample.AccAddress(),
			Namespace: "1",
		}}}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&snftGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgCreateClassOwner          = "op_weight_msg_snft"
		defaultWeightMsgCreateClassOwner int = 100
	)

	var weightMsgCreateClassOwner int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateClassOwner, &weightMsgCreateClassOwner, nil,
		func(_ *rand.Rand) {
			weightMsgCreateClassOwner = defaultWeightMsgCreateClassOwner
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateClassOwner,
		snftsimulation.SimulateMsgCreateClassOwner(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgUpdateClassOwner          = "op_weight_msg_snft"
		defaultWeightMsgUpdateClassOwner int = 100
	)

	var weightMsgUpdateClassOwner int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateClassOwner, &weightMsgUpdateClassOwner, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateClassOwner = defaultWeightMsgUpdateClassOwner
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateClassOwner,
		snftsimulation.SimulateMsgUpdateClassOwner(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgDeleteClassOwner          = "op_weight_msg_snft"
		defaultWeightMsgDeleteClassOwner int = 100
	)

	var weightMsgDeleteClassOwner int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteClassOwner, &weightMsgDeleteClassOwner, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteClassOwner = defaultWeightMsgDeleteClassOwner
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteClassOwner,
		snftsimulation.SimulateMsgDeleteClassOwner(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgRespondClassOwnerTransfer          = "op_weight_msg_snft"
		defaultWeightMsgRespondClassOwnerTransfer int = 100
	)

	var weightMsgRespondClassOwnerTransfer int
	simState.AppParams.GetOrGenerate(opWeightMsgRespondClassOwnerTransfer, &weightMsgRespondClassOwnerTransfer, nil,
		func(_ *rand.Rand) {
			weightMsgRespondClassOwnerTransfer = defaultWeightMsgRespondClassOwnerTransfer
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgRespondClassOwnerTransfer,
		snftsimulation.SimulateMsgRespondClassOwnerTransfer(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgMintAndSend          = "op_weight_msg_snft"
		defaultWeightMsgMintAndSend int = 100
	)

	var weightMsgMintAndSend int
	simState.AppParams.GetOrGenerate(opWeightMsgMintAndSend, &weightMsgMintAndSend, nil,
		func(_ *rand.Rand) {
			weightMsgMintAndSend = defaultWeightMsgMintAndSend
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgMintAndSend,
		snftsimulation.SimulateMsgMintAndSend(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgCreateClassNamespace          = "op_weight_msg_snft"
		defaultWeightMsgCreateClassNamespace int = 100
	)

	var weightMsgCreateClassNamespace int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateClassNamespace, &weightMsgCreateClassNamespace, nil,
		func(_ *rand.Rand) {
			weightMsgCreateClassNamespace = defaultWeightMsgCreateClassNamespace
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateClassNamespace,
		snftsimulation.SimulateMsgCreateClassNamespace(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgUpdateClassNamespace          = "op_weight_msg_snft"
		defaultWeightMsgUpdateClassNamespace int = 100
	)

	var weightMsgUpdateClassNamespace int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateClassNamespace, &weightMsgUpdateClassNamespace, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateClassNamespace = defaultWeightMsgUpdateClassNamespace
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateClassNamespace,
		snftsimulation.SimulateMsgUpdateClassNamespace(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgDeleteClassNamespace          = "op_weight_msg_snft"
		defaultWeightMsgDeleteClassNamespace int = 100
	)

	var weightMsgDeleteClassNamespace int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteClassNamespace, &weightMsgDeleteClassNamespace, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteClassNamespace = defaultWeightMsgDeleteClassNamespace
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteClassNamespace,
		snftsimulation.SimulateMsgDeleteClassNamespace(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
