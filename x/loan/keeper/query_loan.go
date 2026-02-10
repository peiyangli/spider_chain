package keeper

import (
	"context"

	"spider/x/loan/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) ListLoan(ctx context.Context, req *types.QueryAllLoanRequest) (*types.QueryAllLoanResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	loans, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.Loan,
		req.Pagination,
		func(_ collections.Triple[string, uint64, uint64], value types.Loan) (types.Loan, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllLoanResponse{Loan: loans, Pagination: pageRes}, nil
}

func (q queryServer) GetLoan(ctx context.Context, req *types.QueryGetLoanRequest) (*types.QueryGetLoanResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	rng := collections.NewPrefixUntilTripleRange[string, uint64, uint64](req.Borrower)
	itr, err := q.k.Loan.Iterate(ctx, rng)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}
	vals, err := itr.Values()
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}
	return &types.QueryGetLoanResponse{Loan: vals}, nil
}
