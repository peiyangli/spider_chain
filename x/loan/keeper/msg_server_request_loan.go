package keeper

import (
	"bytes"
	"context"

	"spider/x/loan/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	LoanStatusNone      = 0
	LoanStatusRequested = 1
	LoanStatusAprroved  = 2

	CollateralTypeCoin = "coin"
	CollateralTypeNft  = "nft"
)

func (k msgServer) RequestLoan(ctx context.Context, msg *types.MsgRequestLoan) (*types.MsgRequestLoanResponse, error) {
	borrowerAddr, err := k.addressCodec.StringToBytes(msg.Creator)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// TODO: Handle the message
	amount, err := sdk.ParseCoinsNormalized(msg.Amount)
	if err != nil || !amount.IsValid() {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "amount is not a valid Coins object")
	}
	if amount.Empty() {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "amount is empty")
	}

	fee, err := sdk.ParseCoinsNormalized(msg.Fee)
	if err != nil || !fee.IsValid() {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "fee is not a valid Coins object")
	}
	if msg.Deadline < 1 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "deadline should be a positive integer")
	}

	if msg.CollateralType == CollateralTypeCoin {
		collateral, err := sdk.ParseCoinsNormalized(msg.CollateralCoin)
		if err != nil || !collateral.IsValid() {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "collateral is not a valid Coins object")
		}
		if collateral.Empty() {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "collateral is empty")
		}
		sdkError := k.bankKeeper.SendCoinsFromAccountToModule(ctx, borrowerAddr, types.ModuleName, collateral)
		if sdkError != nil {
			return nil, sdkError
		}
	} else if msg.CollateralType == CollateralTypeNft {
		//todo nft
		nftOwner := k.nftKeeper.GetOwner(ctx, msg.CollateralNftClass, msg.CollateralNftId)
		// if !nftOwner.Equals(sdk.AccAddress(borrowerAddr)) {
		// 	return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "not nft owner")
		// }
		if !bytes.Equal(nftOwner, borrowerAddr) {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "not nft owner")
		}
		// k.nftKeeper.Send(ctx, &nft.MsgSend{})
		moduleAddr := k.authKeeper.GetModuleAddress(types.ModuleName)
		if moduleAddr == nil {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "module account not set")
		}
		err = k.nftKeeper.Transfer(ctx, msg.CollateralNftClass, msg.CollateralNftId, moduleAddr)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "collateral not support")
	}

	seq, err := k.LoanSeq.Next(ctx)
	if err != nil {
		return nil, err
	}

	var loan = types.Loan{
		Borrower:           msg.Creator,
		Status:             LoanStatusRequested,
		Seq:                seq,
		Term:               msg.Term,
		ApproveDeadline:    msg.Deadline,
		Amount:             msg.Amount,
		Fee:                msg.Fee,
		Lender:             "",
		CollateralType:     msg.CollateralType,
		CollateralCoin:     msg.CollateralCoin,
		CollateralNftClass: msg.CollateralNftClass,
		CollateralNftId:    msg.CollateralNftId,
	}

	key := collections.Join3(loan.Borrower, loan.Status, loan.Seq)

	err = k.Loan.Set(ctx, key, loan)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	return &types.MsgRequestLoanResponse{}, nil
}
