package keeper

import (
	"context"
	"errors"

	"spider/x/tokenfactory/types"

	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) ListNamespace(ctx context.Context, req *types.QueryAllNamespaceRequest) (*types.QueryAllNamespaceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	namespaces, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.Namespace,
		req.Pagination,
		func(_ string, value types.Namespace) (types.Namespace, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllNamespaceResponse{Namespace: namespaces, Pagination: pageRes}, nil
}

func (q queryServer) GetNamespace(ctx context.Context, req *types.QueryGetNamespaceRequest) (*types.QueryGetNamespaceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, err := q.k.Namespace.Get(ctx, req.Namespace)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetNamespaceResponse{Namespace: val}, nil
}
