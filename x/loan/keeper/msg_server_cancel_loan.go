package keeper

import (
	"context"

	"spider/x/loan/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CancelLoan(ctx context.Context, msg *types.MsgCancelLoan) (*types.MsgCancelLoanResponse, error) {
	borrowerAddr, err := k.addressCodec.StringToBytes(msg.Creator)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	key := collections.Join3(msg.Creator, uint64(LoanStatusRequested), msg.Seq)

	loan, err := k.Loan.Get(ctx, key)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrKeyNotFound, "key %d doesn't exist", key)
	}
	if loan.Borrower != msg.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "Cannot cancel: not the borrower")
	}

	if loan.CollateralType == CollateralTypeCoin {
		collateral, err := sdk.ParseCoinsNormalized(loan.CollateralCoin)
		if collateral.IsValid() && !collateral.Empty() {
			err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, borrowerAddr, collateral)
			if err != nil {
				return nil, err
			}
		}
	} else {
		//todo nft
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "collateral not support")
	}
	err = k.Loan.Remove(ctx, key)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}
	return &types.MsgCancelLoanResponse{}, nil
}
