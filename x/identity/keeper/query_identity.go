package keeper

import (
	"context"
	"errors"

	"spider/x/identity/types"

	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) ListIdentity(ctx context.Context, req *types.QueryAllIdentityRequest) (*types.QueryAllIdentityResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	identitys, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.Identity,
		req.Pagination,
		func(_ string, value types.Identity) (types.Identity, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllIdentityResponse{Identity: identitys, Pagination: pageRes}, nil
}

func (q queryServer) GetIdentity(ctx context.Context, req *types.QueryGetIdentityRequest) (*types.QueryGetIdentityResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, err := q.k.Identity.Get(ctx, req.Uid)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetIdentityResponse{Identity: val}, nil
}
