package keeper

import (
	"context"

	"spider/x/loan/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) LiquidateLoan(ctx context.Context, msg *types.MsgLiquidateLoan) (*types.MsgLiquidateLoanResponse, error) {
	liquidatorAddr, err := k.addressCodec.StringToBytes(msg.Creator)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	key := collections.Join3(msg.Borrower, uint64(LoanStatusAprroved), msg.Seq)
	loan, err := k.Loan.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	blockHeight := uint64(sdkCtx.BlockHeight())
	if blockHeight < loan.RepayDeadline {
		return nil, errorsmod.Wrap(types.ErrDeadline, "Cannot liquidate before deadline")
	}
	if msg.Creator != loan.Lender {
		//public liquidate
		if loan.RepayDeadline-blockHeight < loan.PublicLiquidationDelay {
			return nil, errorsmod.Wrap(types.ErrDeadline, "Cannot liquidate before deadline")
		}
	}
	lenderAddr := liquidatorAddr
	isPublic := loan.Lender != msg.Creator
	if isPublic {
		lenderAddr, err = k.addressCodec.StringToBytes(loan.Lender)
		if err != nil {
			return nil, errorsmod.Wrap(err, "invalid Lender address")
		}
	}

	if loan.CollateralType == CollateralTypeCoin {
		collateral, _ := sdk.ParseCoinsNormalized(loan.CollateralCoin)
		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, lenderAddr, collateral)
		if err != nil {
			return nil, err
		}
	} else {
		//todo nft
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "collateral not support")
	}

	if isPublic {
		//奖励
		if len(loan.PublicLiquidationReward) > 0 {
			amountReward, _ := sdk.ParseCoinsNormalized(loan.PublicLiquidationReward)
			if amountReward.IsValid() && !amountReward.Empty() {
				err = k.bankKeeper.SendCoins(ctx, lenderAddr, liquidatorAddr, amountReward)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	err = k.Loan.Remove(ctx, key)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	return &types.MsgLiquidateLoanResponse{}, nil
}
