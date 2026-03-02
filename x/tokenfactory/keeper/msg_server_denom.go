package keeper

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"spider/x/tokenfactory/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	PermCreate = 0x0001

	NamespaceSeparator = '/'
)

func get_namespace(name string) string {
	pos := strings.IndexByte(name, NamespaceSeparator)
	if pos < 1 {
		return ""
	}
	return name[:pos]
}

func (k msgServer) CreateDenom(ctx context.Context, msg *types.MsgCreateDenom) (*types.MsgCreateDenomResponse, error) {
	creatorAddr, err := k.addressCodec.StringToBytes(msg.Owner)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}
	//to lower
	msg.Denom = strings.ToLower(msg.Denom)

	if msg.Denom == sdk.DefaultBondDenom {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("no bond-denom:  %s", sdk.DefaultBondDenom))
	}

	nskey := get_namespace(msg.Denom)
	var fees []sdk.Coin
	if len(nskey) > 0 {
		//todo check
		ns, err := k.Namespace.Get(ctx, nskey)
		if err != nil {
			return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
		}
		if !ns.CreationFee.IsPositive() {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "not opening for create")
		}
		if msg.CreationFee.IsLT(ns.CreationFee) {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("need creation-fee: %v", ns.CreationFee))
		}

		fees = []sdk.Coin{ns.CreationFee}
	} else {
		operator, err := k.officialKeeper.GetOperator(ctx, msg.Owner, types.ModuleName)
		if err != nil {
			return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
		}
		if PermCreate&operator.GetPermissions() == 0 {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "no perm")
		}
	}

	//check coin type
	var coin = sdk.Coin{
		Denom:  msg.Denom,
		Amount: math.NewInt(1),
	}
	err = coin.Validate()
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("not a coin type: %s", err))
	}

	if len(fees) > 0 {
		if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, creatorAddr, types.ModuleName, fees); err != nil {
			return nil, errorsmod.Wrap(err, "failed to pay")
		}
		// Burn the coins
		if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, fees); err != nil {
			return nil, errorsmod.Wrap(err, "failed to burn coins")
		}
	}

	// Check if the value already exists
	ok, err := k.Denom.Has(ctx, msg.Denom)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	} else if ok {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	//

	var denom = types.Denom{
		Owner:              msg.Owner,
		Denom:              msg.Denom,
		Description:        msg.Description,
		Ticker:             msg.Ticker,
		Precision:          msg.Precision,
		Url:                msg.Url,
		MaxSupply:          msg.MaxSupply,
		Supply:             0,
		CanChangeMaxSupply: msg.CanChangeMaxSupply,
	}

	if err := k.Denom.Set(ctx, denom.Denom, denom); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	return &types.MsgCreateDenomResponse{}, nil
}

func (k msgServer) UpdateDenom(ctx context.Context, msg *types.MsgUpdateDenom) (*types.MsgUpdateDenomResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Owner); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid signer address: %s", err))
	}
	//tolower
	msg.Denom = strings.ToLower(msg.Denom)
	// Check if the value exists
	val, err := k.Denom.Get(ctx, msg.Denom)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
		}

		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	// Checks if the msg owner is the same as the current owner
	if msg.Owner != val.Owner {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	if !val.CanChangeMaxSupply && val.MaxSupply != msg.MaxSupply {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "cannot change maxsupply")
	}
	if !val.CanChangeMaxSupply && msg.CanChangeMaxSupply {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "Cannot revert change maxsupply flag")
	}

	var denom = types.Denom{
		Owner:              msg.Owner,
		Denom:              val.Denom,
		Description:        msg.Description,
		Ticker:             val.Ticker,
		Precision:          val.Precision,
		Url:                msg.Url,
		MaxSupply:          msg.MaxSupply,
		Supply:             val.Supply,
		CanChangeMaxSupply: msg.CanChangeMaxSupply,
	}

	if err := k.Denom.Set(ctx, denom.Denom, denom); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to update denom")
	}

	return &types.MsgUpdateDenomResponse{}, nil
}
