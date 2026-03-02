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

func (q queryServer) ListClassNamespace(ctx context.Context, req *types.QueryAllClassNamespaceRequest) (*types.QueryAllClassNamespaceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	classNamespaces, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.ClassNamespace,
		req.Pagination,
		func(_ string, value types.ClassNamespace) (types.ClassNamespace, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllClassNamespaceResponse{ClassNamespace: classNamespaces, Pagination: pageRes}, nil
}

func (q queryServer) GetClassNamespace(ctx context.Context, req *types.QueryGetClassNamespaceRequest) (*types.QueryGetClassNamespaceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, err := q.k.ClassNamespace.Get(ctx, req.Namespace)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetClassNamespaceResponse{ClassNamespace: val}, nil
}
