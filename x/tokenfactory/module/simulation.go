package tokenfactory

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"spider/testutil/sample"
	tokenfactorysimulation "spider/x/tokenfactory/simulation"
	"spider/x/tokenfactory/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	tokenfactoryGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		DenomMap: []types.Denom{{Owner: sample.AccAddress(),
			Denom: "0",
		}, {Owner: sample.AccAddress(),
			Denom: "1",
		}}, NamespaceMap: []types.Namespace{{Creator: sample.AccAddress(),
			Namespace: "0",
		}, {Creator: sample.AccAddress(),
			Namespace: "1",
		}}}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&tokenfactoryGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgCreateDenom          = "op_weight_msg_tokenfactory"
		defaultWeightMsgCreateDenom int = 100
	)

	var weightMsgCreateDenom int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateDenom, &weightMsgCreateDenom, nil,
		func(_ *rand.Rand) {
			weightMsgCreateDenom = defaultWeightMsgCreateDenom
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateDenom,
		tokenfactorysimulation.SimulateMsgCreateDenom(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgUpdateDenom          = "op_weight_msg_tokenfactory"
		defaultWeightMsgUpdateDenom int = 100
	)

	var weightMsgUpdateDenom int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateDenom, &weightMsgUpdateDenom, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateDenom = defaultWeightMsgUpdateDenom
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateDenom,
		tokenfactorysimulation.SimulateMsgUpdateDenom(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgDeleteDenom          = "op_weight_msg_tokenfactory"
		defaultWeightMsgDeleteDenom int = 100
	)

	var weightMsgDeleteDenom int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteDenom, &weightMsgDeleteDenom, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteDenom = defaultWeightMsgDeleteDenom
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteDenom,
		tokenfactorysimulation.SimulateMsgDeleteDenom(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgMintAndSend          = "op_weight_msg_tokenfactory"
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
		tokenfactorysimulation.SimulateMsgMintAndSend(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgUpdateOwner          = "op_weight_msg_tokenfactory"
		defaultWeightMsgUpdateOwner int = 100
	)

	var weightMsgUpdateOwner int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateOwner, &weightMsgUpdateOwner, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateOwner = defaultWeightMsgUpdateOwner
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateOwner,
		tokenfactorysimulation.SimulateMsgUpdateOwner(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgCreateNamespace          = "op_weight_msg_tokenfactory"
		defaultWeightMsgCreateNamespace int = 100
	)

	var weightMsgCreateNamespace int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateNamespace, &weightMsgCreateNamespace, nil,
		func(_ *rand.Rand) {
			weightMsgCreateNamespace = defaultWeightMsgCreateNamespace
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateNamespace,
		tokenfactorysimulation.SimulateMsgCreateNamespace(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgUpdateNamespace          = "op_weight_msg_tokenfactory"
		defaultWeightMsgUpdateNamespace int = 100
	)

	var weightMsgUpdateNamespace int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateNamespace, &weightMsgUpdateNamespace, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateNamespace = defaultWeightMsgUpdateNamespace
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateNamespace,
		tokenfactorysimulation.SimulateMsgUpdateNamespace(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgDeleteNamespace          = "op_weight_msg_tokenfactory"
		defaultWeightMsgDeleteNamespace int = 100
	)

	var weightMsgDeleteNamespace int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteNamespace, &weightMsgDeleteNamespace, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteNamespace = defaultWeightMsgDeleteNamespace
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteNamespace,
		tokenfactorysimulation.SimulateMsgDeleteNamespace(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
