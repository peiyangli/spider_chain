package keeper

import (
	"context"

	"spider/x/official/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) ListOperator(ctx context.Context, req *types.QueryAllOperatorRequest) (*types.QueryAllOperatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	operators, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.Operator,
		req.Pagination,
		func(_ collections.Pair[string, string], value types.Operator) (types.Operator, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllOperatorResponse{Operator: operators, Pagination: pageRes}, nil
}

func (q queryServer) GetOperator(ctx context.Context, req *types.QueryGetOperatorRequest) (*types.QueryGetOperatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	rng := collections.NewPrefixUntilPairRange[string, string](req.Address)
	itr, err := q.k.Operator.Iterate(ctx, rng)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}
	vals, err := itr.Values()
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}
	return &types.QueryGetOperatorResponse{Operator: vals}, nil
}
