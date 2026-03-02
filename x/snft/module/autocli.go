package snft

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"spider/x/snft/types"
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
					RpcMethod: "ListClassOwner",
					Use:       "list-class-owner",
					Short:     "List all class_owner",
				},
				{
					RpcMethod:      "GetClassOwner",
					Use:            "get-class-owner [id]",
					Short:          "Gets a class_owner",
					Alias:          []string{"show-class-owner"},
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "class_id"}},
				},
				{
					RpcMethod: "ListClassNamespace",
					Use:       "list-class-namespace",
					Short:     "List all class_namespace",
				},
				{
					RpcMethod:      "GetClassNamespace",
					Use:            "get-class-namespace [id]",
					Short:          "Gets a class_namespace",
					Alias:          []string{"show-class-namespace"},
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "namespace"}},
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
					RpcMethod:      "CreateClassOwner",
					Use:            "create-class-owner [class_id] [pending-owner] [creation-fee]",
					Short:          "Create a new class_owner",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "class_id"}, {ProtoField: "pending_owner"}, {ProtoField: "creation_fee"}},
				},
				{
					RpcMethod:      "UpdateClassOwner",
					Use:            "update-class-owner [class_id] [pending-owner] [creation-fee]",
					Short:          "Update class_owner",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "class_id"}, {ProtoField: "pending_owner"}, {ProtoField: "creation_fee"}},
				},
				{
					RpcMethod:      "DeleteClassOwner",
					Use:            "delete-class-owner [class_id]",
					Short:          "Delete class_owner",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "class_id"}},
				},
				{
					RpcMethod:      "RespondClassOwnerTransfer",
					Use:            "respond-class-owner-transfer [class-id] [decision]",
					Short:          "Send a RespondClassOwnerTransfer tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "class_id"}, {ProtoField: "decision"}},
				},
				{
					RpcMethod:      "MintAndSend",
					Use:            "mint-and-send [class-id] [nft-id] [uri] [uri-hash] [recipient]",
					Short:          "Send a MintAndSend tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "class_id"}, {ProtoField: "nft_id"}, {ProtoField: "uri"}, {ProtoField: "uri_hash"}, {ProtoField: "recipient"}},
				},
				{
					RpcMethod:      "CreateClassNamespace",
					Use:            "create-class-namespace [namespace] [creation-fee]",
					Short:          "Create a new class_namespace",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "namespace"}, {ProtoField: "creation_fee"}},
				},
				{
					RpcMethod:      "UpdateClassNamespace",
					Use:            "update-class-namespace [namespace] [creation-fee]",
					Short:          "Update class_namespace",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "namespace"}, {ProtoField: "creation_fee"}},
				},
				{
					RpcMethod:      "DeleteClassNamespace",
					Use:            "delete-class-namespace [namespace]",
					Short:          "Delete class_namespace",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "namespace"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
