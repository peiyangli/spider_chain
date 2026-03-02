package keeper

import (
	"context"
	"errors"

	"spider/x/snft/types"

	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) ListClassOwner(ctx context.Context, req *types.QueryAllClassOwnerRequest) (*types.QueryAllClassOwnerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	classOwners, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.ClassOwner,
		req.Pagination,
		func(_ string, value types.ClassOwner) (types.ClassOwner, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllClassOwnerResponse{ClassOwner: classOwners, Pagination: pageRes}, nil
}

func (q queryServer) GetClassOwner(ctx context.Context, req *types.QueryGetClassOwnerRequest) (*types.QueryGetClassOwnerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, err := q.k.ClassOwner.Get(ctx, req.ClassId)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetClassOwnerResponse{ClassOwner: val}, nil
}
