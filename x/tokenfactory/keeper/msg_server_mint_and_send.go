package keeper

import (
	"context"
	"errors"

	"spider/x/tokenfactory/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) MintAndSend(ctx context.Context, msg *types.MsgMintAndSend) (*types.MsgMintAndSendResponse, error) {
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

	if msg.Owner != val.Owner {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	if val.Supply+msg.Amount > val.MaxSupply {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Cannot mint more than Max Supply")
	}
	moduleAcct := k.authKeeper.GetModuleAddress(types.ModuleName)

	recipientAddress, err := sdk.AccAddressFromBech32(msg.Recipient)
	if err != nil {
		return nil, err
	}

	var mintCoins sdk.Coins

	mintCoins = mintCoins.Add(sdk.NewCoin(msg.Denom, math.NewInt(int64(msg.Amount))))
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, mintCoins); err != nil {
		return nil, err
	}
	if err := k.bankKeeper.SendCoins(ctx, moduleAcct, recipientAddress, mintCoins); err != nil {
		return nil, err
	}

	var denom = types.Denom{
		Owner:              val.Owner,
		Denom:              val.Denom,
		Description:        val.Description,
		MaxSupply:          val.MaxSupply,
		Supply:             val.Supply + msg.Amount,
		Precision:          val.Precision,
		Ticker:             val.Ticker,
		Url:                val.Url,
		CanChangeMaxSupply: val.CanChangeMaxSupply,
	}

	if err := k.Denom.Set(ctx, denom.Denom, denom); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to mint and send denom")
	}

	return &types.MsgMintAndSendResponse{}, nil
}
