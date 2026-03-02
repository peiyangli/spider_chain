package keeper

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"spider/x/snft/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/x/nft"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	PermCreate = 0x0001
)

func get_name_space(name string) string {
	pos := strings.IndexByte(name, types.NftClassidNamespaceSeparator)
	if pos < 1 {
		return ""
	}
	return name[:pos]
}

func (k msgServer) CreateClassOwner(ctx context.Context, msg *types.MsgCreateClassOwner) (*types.MsgCreateClassOwnerResponse, error) {
	creatorAddr, err := k.addressCodec.StringToBytes(msg.Owner)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	// NftClassidNamespaceSeparator
	// todo check msg.ClassId: lower case
	//to lower
	msg.ClassId = strings.ToLower(msg.ClassId)

	nskey := get_name_space(msg.ClassId)
	var fees []sdk.Coin
	if len(nskey) > 0 {
		//todo check
		ns, err := k.ClassNamespace.Get(ctx, nskey)
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

	// Check if the value already exists
	ok, err := k.ClassOwner.Has(ctx, msg.ClassId)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	} else if ok {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
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

	err = k.nftKeeper.SaveClass(ctx, nft.Class{
		Id: msg.ClassId,
	})
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	var classOwner = types.ClassOwner{
		Owner:        msg.Owner,
		ClassId:      msg.ClassId,
		PendingOwner: msg.PendingOwner,
		CreationFee:  msg.CreationFee,
	}

	if err := k.ClassOwner.Set(ctx, classOwner.ClassId, classOwner); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	return &types.MsgCreateClassOwnerResponse{}, nil
}

func (k msgServer) UpdateClassOwner(ctx context.Context, msg *types.MsgUpdateClassOwner) (*types.MsgUpdateClassOwnerResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Owner); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid signer address: %s", err))
	}

	// Check if the value exists
	val, err := k.ClassOwner.Get(ctx, msg.ClassId)
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

	var classOwner = types.ClassOwner{
		Owner:        msg.Owner,
		ClassId:      msg.ClassId,
		PendingOwner: msg.PendingOwner,
		CreationFee:  msg.CreationFee,
	}

	if err := k.ClassOwner.Set(ctx, classOwner.ClassId, classOwner); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to update classOwner")
	}

	return &types.MsgUpdateClassOwnerResponse{}, nil
}

func (k msgServer) DeleteClassOwner(ctx context.Context, msg *types.MsgDeleteClassOwner) (*types.MsgDeleteClassOwnerResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Owner); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid signer address: %s", err))
	}

	// Check if the value exists
	val, err := k.ClassOwner.Get(ctx, msg.ClassId)
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

	if err := k.ClassOwner.Remove(ctx, msg.ClassId); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to remove classOwner")
	}

	return &types.MsgDeleteClassOwnerResponse{}, nil
}
