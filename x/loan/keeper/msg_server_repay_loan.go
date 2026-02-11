package keeper

import (
	"context"

	"spider/x/loan/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) RepayLoan(ctx context.Context, msg *types.MsgRepayLoan) (*types.MsgRepayLoanResponse, error) {
	repayorAddr, err := k.addressCodec.StringToBytes(msg.Creator)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}
	borrowerAddr := repayorAddr
	if msg.Creator != msg.Borrower {
		borrowerAddr, err = k.addressCodec.StringToBytes(msg.Borrower)
		if err != nil {
			return nil, errorsmod.Wrap(err, "invalid authority address")
		}
	}

	key := collections.Join3(msg.Borrower, uint64(LoanStatusAprroved), msg.Seq)
	loan, err := k.Loan.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	lenderAddr, err := k.addressCodec.StringToBytes(loan.Lender)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	amount, _ := sdk.ParseCoinsNormalized(loan.Amount)
	fee, _ := sdk.ParseCoinsNormalized(loan.Fee)
	//本金
	err = k.bankKeeper.SendCoins(ctx, repayorAddr, lenderAddr, amount)
	if err != nil {
		return nil, err
	}
	//利息
	err = k.bankKeeper.SendCoins(ctx, repayorAddr, lenderAddr, fee)
	if err != nil {
		return nil, err
	}
	if loan.CollateralType == CollateralTypeCoin {
		collateral, _ := sdk.ParseCoinsNormalized(loan.CollateralCoin)
		//取回质押
		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, borrowerAddr, collateral)
		if err != nil {
			return nil, err
		}
	} else {
		//todo nft
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "collateral not support")
	}

	err = k.Loan.Remove(ctx, key)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	return &types.MsgRepayLoanResponse{}, nil
}
