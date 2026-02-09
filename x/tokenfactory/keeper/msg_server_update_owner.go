package keeper

import (
	"context"
	"errors"

	"spider/x/tokenfactory/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) UpdateOwner(ctx context.Context, msg *types.MsgUpdateOwner) (*types.MsgUpdateOwnerResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Owner); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	val, err := k.Denom.Get(ctx, msg.Denom)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "Denom not found")
		}

		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	// Checks if the the msg owner is the same as the current owner
	if msg.Owner != val.Owner {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	var denom = types.Denom{
		Owner:              msg.NewOwner,
		Denom:              msg.Denom,
		Description:        val.Description,
		MaxSupply:          val.MaxSupply,
		Supply:             val.Supply,
		Precision:          val.Precision,
		Ticker:             val.Ticker,
		Url:                val.Url,
		CanChangeMaxSupply: val.CanChangeMaxSupply,
	}

	if err := k.Denom.Set(ctx, denom.Denom, denom); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to update denom owner")
	}

	return &types.MsgUpdateOwnerResponse{}, nil
}
