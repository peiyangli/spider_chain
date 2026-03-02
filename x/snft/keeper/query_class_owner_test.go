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

func createNClassOwner(keeper keeper.Keeper, ctx context.Context, n int) []types.ClassOwner {
	items := make([]types.ClassOwner, n)
	for i := range items {
		items[i].ClassId = strconv.Itoa(i)
		items[i].PendingOwner = strconv.Itoa(i)
		items[i].CreationFee = sdk.NewInt64Coin(`token`, int64(i+100))
		_ = keeper.ClassOwner.Set(ctx, items[i].ClassId, items[i])
	}
	return items
}

func TestClassOwnerQuerySingle(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNClassOwner(f.keeper, f.ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetClassOwnerRequest
		response *types.QueryGetClassOwnerResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetClassOwnerRequest{
				ClassId: msgs[0].ClassId,
			},
			response: &types.QueryGetClassOwnerResponse{ClassOwner: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetClassOwnerRequest{
				ClassId: msgs[1].ClassId,
			},
			response: &types.QueryGetClassOwnerResponse{ClassOwner: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetClassOwnerRequest{
				ClassId: strconv.Itoa(100000),
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
			response, err := qs.GetClassOwner(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.EqualExportedValues(t, tc.response, response)
			}
		})
	}
}

func TestClassOwnerQueryPaginated(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNClassOwner(f.keeper, f.ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllClassOwnerRequest {
		return &types.QueryAllClassOwnerRequest{
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
			resp, err := qs.ListClassOwner(f.ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.ClassOwner), step)
			require.Subset(t, msgs, resp.ClassOwner)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListClassOwner(f.ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.ClassOwner), step)
			require.Subset(t, msgs, resp.ClassOwner)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.ListClassOwner(f.ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.EqualExportedValues(t, msgs, resp.ClassOwner)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := qs.ListClassOwner(f.ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
