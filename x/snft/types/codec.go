package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterInterfaces(registrar codectypes.InterfaceRegistry) {
	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateClassNamespace{},
		&MsgUpdateClassNamespace{},
		&MsgDeleteClassNamespace{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgMintAndSend{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRespondClassOwnerTransfer{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateClassOwner{},
		&MsgUpdateClassOwner{},
		&MsgDeleteClassOwner{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
	)
	msgservice.RegisterMsgServiceDesc(registrar, &_Msg_serviceDesc)
}
