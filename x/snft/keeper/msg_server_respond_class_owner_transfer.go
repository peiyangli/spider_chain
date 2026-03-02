package keeper

import (
	"context"
	"errors"

	"spider/x/snft/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	RespondClassOwnerTransferDecisionAccept = 1
	RespondClassOwnerTransferDecisionReject = 2
)

func (k msgServer) RespondClassOwnerTransfer(ctx context.Context, msg *types.MsgRespondClassOwnerTransfer) (*types.MsgRespondClassOwnerTransferResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Owner); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	switch msg.Decision {
	case RespondClassOwnerTransferDecisionAccept:
	case RespondClassOwnerTransferDecisionReject:
	default:
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect Decision")
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
	if msg.Owner != val.PendingOwner {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	var classOwner = val
	classOwner.PendingOwner = ""

	switch msg.Decision {
	case RespondClassOwnerTransferDecisionAccept:
		classOwner.Owner = msg.Owner
	case RespondClassOwnerTransferDecisionReject:
	default:
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect Decision")
	}

	if err := k.ClassOwner.Set(ctx, classOwner.ClassId, classOwner); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to update classOwner")
	}
	return &types.MsgRespondClassOwnerTransferResponse{}, nil
}
