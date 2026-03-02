package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"

	"spider/x/snft/types"
)

type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec
	// Address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority []byte

	Schema collections.Schema
	Params collections.Item[types.Params]

	authKeeper     types.AuthKeeper
	bankKeeper     types.BankKeeper
	nftKeeper      types.NftKeeper
	officialKeeper types.OfficialKeeper
	ClassOwner     collections.Map[string, types.ClassOwner]
	ClassNamespace collections.Map[string, types.ClassNamespace]
}

func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
	authority []byte,

	authKeeper types.AuthKeeper,
	bankKeeper types.BankKeeper,
	nftKeeper types.NftKeeper,
	officialKeeper types.OfficialKeeper,
) Keeper {
	if _, err := addressCodec.BytesToString(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address %s: %s", authority, err))
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		storeService: storeService,
		cdc:          cdc,
		addressCodec: addressCodec,
		authority:    authority,

		authKeeper:     authKeeper,
		bankKeeper:     bankKeeper,
		nftKeeper:      nftKeeper,
		officialKeeper: officialKeeper,
		Params:         collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		ClassOwner:     collections.NewMap(sb, types.ClassOwnerKey, "classOwner", collections.StringKey, codec.CollValue[types.ClassOwner](cdc)), ClassNamespace: collections.NewMap(sb, types.ClassNamespaceKey, "classNamespace", collections.StringKey, codec.CollValue[types.ClassNamespace](cdc))}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() []byte {
	return k.authority
}
