package keeper

import (
	"context"
	"errors"
	"fmt"

	"spider/x/tokenfactory/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateNamespace(ctx context.Context, msg *types.MsgCreateNamespace) (*types.MsgCreateNamespaceResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	_, err := k.officialKeeper.GetOperator(ctx, msg.Creator, types.ModuleName)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	// Check if the value already exists
	ok, err := k.Namespace.Has(ctx, msg.Namespace)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	} else if ok {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	var namespace = types.Namespace{
		Creator:     msg.Creator,
		Namespace:   msg.Namespace,
		CreationFee: msg.CreationFee,
	}

	if err := k.Namespace.Set(ctx, namespace.Namespace, namespace); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	return &types.MsgCreateNamespaceResponse{}, nil
}

func (k msgServer) UpdateNamespace(ctx context.Context, msg *types.MsgUpdateNamespace) (*types.MsgUpdateNamespaceResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid signer address: %s", err))
	}

	_, err := k.officialKeeper.GetOperator(ctx, msg.Creator, types.ModuleName)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	// Check if the value exists
	val, err := k.Namespace.Get(ctx, msg.Namespace)
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

	var namespace = types.Namespace{
		Creator:     msg.Creator,
		Namespace:   msg.Namespace,
		CreationFee: msg.CreationFee,
	}

	if err := k.Namespace.Set(ctx, namespace.Namespace, namespace); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to update namespace")
	}

	return &types.MsgUpdateNamespaceResponse{}, nil
}

func (k msgServer) DeleteNamespace(ctx context.Context, msg *types.MsgDeleteNamespace) (*types.MsgDeleteNamespaceResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid signer address: %s", err))
	}

	_, err := k.officialKeeper.GetOperator(ctx, msg.Creator, types.ModuleName)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	// Check if the value exists
	val, err := k.Namespace.Get(ctx, msg.Namespace)
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

	if err := k.Namespace.Remove(ctx, msg.Namespace); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to remove namespace")
	}

	return &types.MsgDeleteNamespaceResponse{}, nil
}
