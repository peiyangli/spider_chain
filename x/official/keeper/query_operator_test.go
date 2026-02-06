package keeper_test

import (
	"context"
	"strconv"
	"testing"

	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"spider/x/official/keeper"
	"spider/x/official/types"
)

func createNOperator(keeper keeper.Keeper, ctx context.Context, n int) []types.Operator {
	items := make([]types.Operator, n)
	for i := range items {
		items[i].Address = strconv.Itoa(i)
		items[i].Module = strconv.Itoa(i)
		items[i].Name = strconv.Itoa(i)
		items[i].Role = uint64(i)
		items[i].Permissions = uint64(i)
		_ = keeper.Operator.Set(ctx, collections.Join(items[i].Address, items[i].Module), items[i])
	}
	return items
}

func TestOperatorQuerySingle(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNOperator(f.keeper, f.ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetOperatorRequest
		response *types.QueryGetOperatorResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetOperatorRequest{
				Address: msgs[0].Address,
			},
			response: &types.QueryGetOperatorResponse{Operator: msgs[0:1]},
		},
		{
			desc: "Second",
			request: &types.QueryGetOperatorRequest{
				Address: msgs[1].Address,
			},
			response: &types.QueryGetOperatorResponse{Operator: msgs[1:2]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetOperatorRequest{
				Address: strconv.Itoa(100000),
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
			response, err := qs.GetOperator(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.EqualExportedValues(t, tc.response, response)
			}
		})
	}
}

func TestOperatorQueryPaginated(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNOperator(f.keeper, f.ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllOperatorRequest {
		return &types.QueryAllOperatorRequest{
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
			resp, err := qs.ListOperator(f.ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Operator), step)
			require.Subset(t, msgs, resp.Operator)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListOperator(f.ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Operator), step)
			require.Subset(t, msgs, resp.Operator)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.ListOperator(f.ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.EqualExportedValues(t, msgs, resp.Operator)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := qs.ListOperator(f.ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
