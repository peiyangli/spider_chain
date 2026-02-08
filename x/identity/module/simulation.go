package identity

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"spider/testutil/sample"
	identitysimulation "spider/x/identity/simulation"
	"spider/x/identity/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	identityGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		IdentityMap: []types.Identity{{Creator: sample.AccAddress(),
			Uid: "0",
		}, {Creator: sample.AccAddress(),
			Uid: "1",
		}}}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&identityGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgCreateIdentity          = "op_weight_msg_identity"
		defaultWeightMsgCreateIdentity int = 100
	)

	var weightMsgCreateIdentity int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateIdentity, &weightMsgCreateIdentity, nil,
		func(_ *rand.Rand) {
			weightMsgCreateIdentity = defaultWeightMsgCreateIdentity
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateIdentity,
		identitysimulation.SimulateMsgCreateIdentity(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgUpdateIdentity          = "op_weight_msg_identity"
		defaultWeightMsgUpdateIdentity int = 100
	)

	var weightMsgUpdateIdentity int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateIdentity, &weightMsgUpdateIdentity, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateIdentity = defaultWeightMsgUpdateIdentity
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateIdentity,
		identitysimulation.SimulateMsgUpdateIdentity(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgDeleteIdentity          = "op_weight_msg_identity"
		defaultWeightMsgDeleteIdentity int = 100
	)

	var weightMsgDeleteIdentity int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteIdentity, &weightMsgDeleteIdentity, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteIdentity = defaultWeightMsgDeleteIdentity
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteIdentity,
		identitysimulation.SimulateMsgDeleteIdentity(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
