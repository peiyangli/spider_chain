package official

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"spider/x/official/types"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: types.Query_serviceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				{
					RpcMethod: "ListOperator",
					Use:       "list-operator",
					Short:     "List all operator",
				},
				{
					RpcMethod:      "GetOperator",
					Use:            "get-operator [id]",
					Short:          "Gets a operator",
					Alias:          []string{"show-operator"},
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "address"}},
				},
				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              types.Msg_serviceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod:      "CreateOperator",
					Use:            "create-operator [address] [module] [name] [role] [permissions]",
					Short:          "Create a new operator",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "address"}, {ProtoField: "module"}, {ProtoField: "name"}, {ProtoField: "role"}, {ProtoField: "permissions"}},
				},
				{
					RpcMethod:      "UpdateOperator",
					Use:            "update-operator [address] [module] [name] [role] [permissions]",
					Short:          "Update operator",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "address"}, {ProtoField: "module"}, {ProtoField: "name"}, {ProtoField: "role"}, {ProtoField: "permissions"}},
				},
				{
					RpcMethod:      "DeleteOperator",
					Use:            "delete-operator [address]",
					Short:          "Delete operator",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "address"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
