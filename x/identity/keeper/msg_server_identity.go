package keeper

import (
	"context"
	"errors"
	"fmt"

	"spider/x/identity/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	PermCreate = 0x0001
	PermDelete = 0x0002
)

func (k msgServer) CreateIdentity(ctx context.Context, msg *types.MsgCreateIdentity) (*types.MsgCreateIdentityResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	//only operator can create uid->pub
	operator, err := k.officialKeeper.GetOperator(ctx, msg.Creator, types.ModuleName)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}
	if PermCreate&operator.GetPermissions() == 0 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "no perm")
	}

	// Check if the value already exists
	ok, err := k.Identity.Has(ctx, msg.Uid)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	} else if ok {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	var identity = types.Identity{
		Creator: msg.Creator,
		Uid:     msg.Uid,
		Owner:   msg.Owner,
		Idkey:   msg.Idkey,
		Msgkey:  msg.Msgkey,
	}

	if err := k.Identity.Set(ctx, identity.Uid, identity); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	return &types.MsgCreateIdentityResponse{}, nil
}

func (k msgServer) UpdateIdentity(ctx context.Context, msg *types.MsgUpdateIdentity) (*types.MsgUpdateIdentityResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid signer address: %s", err))
	}

	// Check if the value exists
	val, err := k.Identity.Get(ctx, msg.Uid)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
		}

		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != val.Owner {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	var identity = types.Identity{
		Creator: val.Creator,
		Uid:     val.Uid,
		Owner:   val.Owner,
		Idkey:   val.Idkey,
		Msgkey:  msg.Msgkey, //only msgkey
	}

	if err := k.Identity.Set(ctx, identity.Uid, identity); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to update identity")
	}

	return &types.MsgUpdateIdentityResponse{}, nil
}

func (k msgServer) DeleteIdentity(ctx context.Context, msg *types.MsgDeleteIdentity) (*types.MsgDeleteIdentityResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid signer address: %s", err))
	}
	// Check if the value exists
	val, err := k.Identity.Get(ctx, msg.Uid)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
		}

		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != val.Owner {
		//operator
		operator, err := k.officialKeeper.GetOperator(ctx, msg.Creator, types.ModuleName)
		if err != nil {
			return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
		}
		if PermDelete&operator.GetPermissions() == 0 {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "no perm nor owner")
		}
	}

	if err := k.Identity.Remove(ctx, msg.Uid); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to remove identity")
	}

	return &types.MsgDeleteIdentityResponse{}, nil
}
