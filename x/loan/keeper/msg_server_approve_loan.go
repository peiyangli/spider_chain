package keeper

import (
	"context"

	"spider/x/loan/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) ApproveLoan(ctx context.Context, msg *types.MsgApproveLoan) (*types.MsgApproveLoanResponse, error) {
	lenderAddr, err := k.addressCodec.StringToBytes(msg.Creator)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}
	if len(msg.PublicLiquidationReward) > 0 {
		amountReward, err := sdk.ParseCoinsNormalized(msg.PublicLiquidationReward)
		if err != nil || !amountReward.IsValid() {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "reward is not a valid Coins object")
		}
	}

	key := collections.Join3(msg.Borrower, uint64(LoanStatusRequested), msg.Seq)
	loan, err := k.Loan.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if loan.Borrower != msg.Borrower {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "not borrower")
	}
	borrowerAddr, err := k.addressCodec.StringToBytes(loan.Borrower)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid borrower address")
	}
	if loan.CollateralType == CollateralTypeCoin {
		amount, err := sdk.ParseCoinsNormalized(loan.Amount)
		if err != nil || !amount.IsValid() {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "amount is not a valid Coins object")
		}
		if amount.Empty() {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "amount is empty")
		}
		err = k.bankKeeper.SendCoins(ctx, lenderAddr, borrowerAddr, amount)
		if err != nil {
			return nil, err
		}
	} else {
		//todo nft
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "collateral not support")
	}
	err = k.Loan.Remove(ctx, key)
	if err != nil {
		return nil, err
	}
	loan.Lender = msg.Creator
	loan.Status = uint64(LoanStatusAprroved)
	loan.PublicLiquidationDelay = msg.PublicLiquidationDelay
	loan.PublicLiquidationReward = msg.PublicLiquidationReward

	key2 := collections.Join3(loan.Borrower, loan.Status, msg.Seq)
	err = k.Loan.Set(ctx, key2, loan)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	return &types.MsgApproveLoanResponse{}, nil
}
