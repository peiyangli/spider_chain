package keeper

import (
	"context"
	"errors"
	"fmt"

	"spider/x/official/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	PermCreate = 0x0001
	PermDelete = 0x0002
	RoleSuper  = 0x0100
)

func (k msgServer) CreateOperator(ctx context.Context, msg *types.MsgCreateOperator) (*types.MsgCreateOperatorResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	// if len(msg.Module) < 1 {
	// 	return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "you must be a manager")
	// }

	// check if the creator is manager
	opt, err := k.GetOperator(ctx, msg.Creator, types.ModuleName) //Has(ctx, collections.Join(msg.Creator, ""))
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}
	if PermCreate&opt.GetPermissions() == 0 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "no perm")
	}

	if msg.Module == types.ModuleName {
		if opt.GetRole() < RoleSuper {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "need super role")
		}
	}

	// Check if the value already exists
	var key = collections.Join(msg.Address, msg.Module)
	ok, err := k.Operator.Has(ctx, key)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	} else if ok {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "address-module already set, use UpdateOperator to modify")
	}

	var operator = types.Operator{
		Creator:     msg.Creator,
		Address:     msg.Address,
		Module:      msg.Module,
		Name:        msg.Name,
		Role:        msg.Role,
		Permissions: msg.Permissions,
	}

	if err := k.Operator.Set(ctx, key, operator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	return &types.MsgCreateOperatorResponse{}, nil
}

func (k msgServer) UpdateOperator(ctx context.Context, msg *types.MsgUpdateOperator) (*types.MsgUpdateOperatorResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid signer address: %s", err))
	}
	// check if the creator is manager
	opt, err := k.GetOperator(ctx, msg.Creator, types.ModuleName) //Has(ctx, collections.Join(msg.Creator, ""))
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}
	if msg.Module == types.ModuleName {
		if opt.GetRole() < RoleSuper {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "need super role")
		}
	}

	// Check if the value exists
	var key = collections.Join(msg.Address, msg.Module)
	val, err := k.Operator.Get(ctx, key)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "address-module not set")
		}

		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != val.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	var operator = types.Operator{
		Creator:     msg.Creator,
		Address:     msg.Address,
		Module:      msg.Module,
		Name:        msg.Name,
		Role:        msg.Role,
		Permissions: msg.Permissions,
	}

	if err := k.Operator.Set(ctx, key, operator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to update operator")
	}

	return &types.MsgUpdateOperatorResponse{}, nil
}

func (k msgServer) DeleteOperator(ctx context.Context, msg *types.MsgDeleteOperator) (*types.MsgDeleteOperatorResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid signer address: %s", err))
	}

	var role uint64 = 0
	var isSuper = false
	// check if the creator is manager
	if msg.Creator != msg.Address {
		//not self
		opt, err := k.GetOperator(ctx, msg.Creator, types.ModuleName) //Has(ctx, collections.Join(msg.Creator, ""))
		if err != nil {
			return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
		}
		role = opt.GetRole()
		isSuper = role >= RoleSuper
	}

	rng := collections.NewPrefixUntilPairRange[string, string](msg.Address)
	itr, err := k.Operator.Iterate(ctx, rng)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}
	defer itr.Close()
	var n int
	for ; itr.Valid(); itr.Next() {
		key, err := itr.Key()
		if err != nil {
			return nil, err
		}

		if key.K2() == types.ModuleName {
			//canot remove official manager
			if !isSuper {
				continue
			}
			targetOpter, err := itr.Value()
			if err != nil {
				return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
			}
			if role <= targetOpter.Role {
				continue
			}
		}

		if err := k.Operator.Remove(ctx, key); err != nil {
			return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to remove operator")
		}
		n++
	}
	if n < 1 {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "ok but operator removed")
	}
	return &types.MsgDeleteOperatorResponse{}, nil
}
