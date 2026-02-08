package identity

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"spider/x/identity/types"
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
					RpcMethod: "ListIdentity",
					Use:       "list-identity",
					Short:     "List all identity",
				},
				{
					RpcMethod:      "GetIdentity",
					Use:            "get-identity [id]",
					Short:          "Gets a identity",
					Alias:          []string{"show-identity"},
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "uid"}},
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
					RpcMethod:      "CreateIdentity",
					Use:            "create-identity [uid] [owner] [idkey] [msgkey]",
					Short:          "Create a new identity",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "uid"}, {ProtoField: "owner"}, {ProtoField: "idkey"}, {ProtoField: "msgkey", Varargs: true}},
				},
				{
					RpcMethod:      "UpdateIdentity",
					Use:            "update-identity [uid] [owner] [idkey] [msgkey]",
					Short:          "Update identity",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "uid"}, {ProtoField: "owner"}, {ProtoField: "idkey"}, {ProtoField: "msgkey", Varargs: true}},
				},
				{
					RpcMethod:      "DeleteIdentity",
					Use:            "delete-identity [uid]",
					Short:          "Delete identity",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "uid"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
