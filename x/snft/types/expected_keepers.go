package types

import (
	"context"
	"spider/x/official/types"

	"cosmossdk.io/core/address"
	"cosmossdk.io/x/nft"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type NftKeeper interface {
	// TODO Add methods imported from nft should be defined here
	GetOwner(ctx context.Context, classID, nftID string) sdk.AccAddress
	Transfer(ctx context.Context, classID string, nftID string, receiver sdk.AccAddress) error

	Mint(ctx context.Context, token nft.NFT, receiver sdk.AccAddress) error

	SaveClass(ctx context.Context, class nft.Class) error
}

type OfficialKeeper interface {
	// TODO Add methods imported from official should be defined here
	GetOperator(ctx context.Context, address, module string) (types.OperatorI, error)
}

// AuthKeeper defines the expected interface for the Auth module.
type AuthKeeper interface {
	AddressCodec() address.Codec
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI // only used for simulation
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface for the Bank module.
type BankKeeper interface {
	SpendableCoins(context.Context, sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here
	BurnCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
}

// ParamSubspace defines the expected Subspace interface for parameters.
type ParamSubspace interface {
	Get(context.Context, []byte, interface{})
	Set(context.Context, []byte, interface{})
}
