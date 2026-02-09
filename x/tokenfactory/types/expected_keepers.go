package types

import (
	"context"
	"spider/x/official/types"

	"cosmossdk.io/core/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type OfficialKeeper interface {
	// TODO Add methods imported from official should be defined here
	GetOperator(ctx context.Context, address, module string) (types.OperatorI, error)
}

// AuthKeeper defines the expected interface for the Auth module.
type AuthKeeper interface {
	AddressCodec() address.Codec
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI // only used for simulation
	// Methods imported from account should be defined here
	GetModuleAddress(moduleName string) sdk.AccAddress
	GetModuleAccount(ctx context.Context, moduleName string) sdk.ModuleAccountI
}

// BankKeeper defines the expected interface for the Bank module.
type BankKeeper interface {
	SpendableCoins(context.Context, sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here
	SendCoins(ctx context.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
}

// ParamSubspace defines the expected Subspace interface for parameters.
type ParamSubspace interface {
	Get(context.Context, []byte, interface{})
	Set(context.Context, []byte, interface{})
}
