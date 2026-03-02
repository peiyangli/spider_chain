package keeper

import (
	"context"
	"errors"
	"fmt"

	"spider/x/snft/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateClassNamespace(ctx context.Context, msg *types.MsgCreateClassNamespace) (*types.MsgCreateClassNamespaceResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	//todo a-z lower case

	_, err := k.officialKeeper.GetOperator(ctx, msg.Creator, types.ModuleName)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	// Check if the value already exists
	ok, err := k.ClassNamespace.Has(ctx, msg.Namespace)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	} else if ok {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	var classNamespace = types.ClassNamespace{
		Creator:     msg.Creator,
		Namespace:   msg.Namespace,
		CreationFee: msg.CreationFee,
	}

	if err := k.ClassNamespace.Set(ctx, classNamespace.Namespace, classNamespace); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	return &types.MsgCreateClassNamespaceResponse{}, nil
}

func (k msgServer) UpdateClassNamespace(ctx context.Context, msg *types.MsgUpdateClassNamespace) (*types.MsgUpdateClassNamespaceResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid signer address: %s", err))
	}

	_, err := k.officialKeeper.GetOperator(ctx, msg.Creator, types.ModuleName)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	// Check if the value exists
	val, err := k.ClassNamespace.Get(ctx, msg.Namespace)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
		}

		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != val.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	var classNamespace = types.ClassNamespace{
		Creator:     msg.Creator,
		Namespace:   msg.Namespace,
		CreationFee: msg.CreationFee,
	}

	if err := k.ClassNamespace.Set(ctx, classNamespace.Namespace, classNamespace); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to update classNamespace")
	}

	return &types.MsgUpdateClassNamespaceResponse{}, nil
}

func (k msgServer) DeleteClassNamespace(ctx context.Context, msg *types.MsgDeleteClassNamespace) (*types.MsgDeleteClassNamespaceResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid signer address: %s", err))
	}

	_, err := k.officialKeeper.GetOperator(ctx, msg.Creator, types.ModuleName)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	// Check if the value exists
	val, err := k.ClassNamespace.Get(ctx, msg.Namespace)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
		}

		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != val.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	if err := k.ClassNamespace.Remove(ctx, msg.Namespace); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to remove classNamespace")
	}

	return &types.MsgDeleteClassNamespaceResponse{}, nil
}
