package loan

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"spider/x/loan/types"
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
					RpcMethod: "ListLoan",
					Use:       "list-loan",
					Short:     "List all loan",
				},
				{
					RpcMethod:      "GetLoan",
					Use:            "get-loan [id]",
					Short:          "Gets a loan",
					Alias:          []string{"show-loan"},
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "borrower"}},
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
					RpcMethod:      "RequestLoan",
					Use:            "request-loan [deadline] [amount] [fee] [collateral-type] [collateral-coin] [collateral-nft-class] [collateral-nft-id]",
					Short:          "Send a request-loan tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "deadline"}, {ProtoField: "amount"}, {ProtoField: "fee"}, {ProtoField: "collateral_type"}, {ProtoField: "collateral_coin"}, {ProtoField: "collateral_nft_class"}, {ProtoField: "collateral_nft_id"}},
				},
				{
					RpcMethod:      "ApproveLoan",
					Use:            "approve-loan [seq] [borrower] [public-liquidation-delay] [public-liquidation-reward]",
					Short:          "Send a approve-loan tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "seq"}, {ProtoField: "borrower"}, {ProtoField: "public_liquidation_delay"}, {ProtoField: "public_liquidation_reward"}},
				},
				{
					RpcMethod:      "CancelLoan",
					Use:            "cancel-loan [seq]",
					Short:          "Send a cancel-loan tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "seq"}},
				},
				{
					RpcMethod:      "RepayLoan",
					Use:            "repay-loan [seq] [borrower]",
					Short:          "Send a repay-loan tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "seq"}, {ProtoField: "borrower"}},
				},
				{
					RpcMethod:      "LiquidateLoan",
					Use:            "liquidate-loan [seq] [borrower]",
					Short:          "Send a liquidate-loan tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "seq"}, {ProtoField: "borrower"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
