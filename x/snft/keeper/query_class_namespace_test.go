package keeper_test

import (
	"context"
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"spider/x/snft/keeper"
	"spider/x/snft/types"
)

func createNClassNamespace(keeper keeper.Keeper, ctx context.Context, n int) []types.ClassNamespace {
	items := make([]types.ClassNamespace, n)
	for i := range items {
		items[i].Namespace = strconv.Itoa(i)
		items[i].CreationFee = sdk.NewInt64Coin(`token`, int64(i+100))
		_ = keeper.ClassNamespace.Set(ctx, items[i].Namespace, items[i])
	}
	return items
}

func TestClassNamespaceQuerySingle(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNClassNamespace(f.keeper, f.ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetClassNamespaceRequest
		response *types.QueryGetClassNamespaceResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetClassNamespaceRequest{
				Namespace: msgs[0].Namespace,
			},
			response: &types.QueryGetClassNamespaceResponse{ClassNamespace: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetClassNamespaceRequest{
				Namespace: msgs[1].Namespace,
			},
			response: &types.QueryGetClassNamespaceResponse{ClassNamespace: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetClassNamespaceRequest{
				Namespace: strconv.Itoa(100000),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := qs.GetClassNamespace(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.EqualExportedValues(t, tc.response, response)
			}
		})
	}
}

func TestClassNamespaceQueryPaginated(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNClassNamespace(f.keeper, f.ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllClassNamespaceRequest {
		return &types.QueryAllClassNamespaceRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListClassNamespace(f.ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.ClassNamespace), step)
			require.Subset(t, msgs, resp.ClassNamespace)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListClassNamespace(f.ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.ClassNamespace), step)
			require.Subset(t, msgs, resp.ClassNamespace)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.ListClassNamespace(f.ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.EqualExportedValues(t, msgs, resp.ClassNamespace)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := qs.ListClassNamespace(f.ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
