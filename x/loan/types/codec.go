package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterInterfaces(registrar codectypes.InterfaceRegistry) {
	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgLiquidateLoan{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRepayLoan{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCancelLoan{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgApproveLoan{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRequestLoan{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
	)
	msgservice.RegisterMsgServiceDesc(registrar, &_Msg_serviceDesc)
}
