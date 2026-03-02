package keeper_test

import (
	"strconv"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"spider/x/snft/keeper"
	"spider/x/snft/types"
)

func TestClassOwnerMsgServerCreate(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)
	owner, err := f.addressCodec.BytesToString([]byte("signerAddr__________________"))
	require.NoError(t, err)

	for i := 0; i < 5; i++ {
		expected := &types.MsgCreateClassOwner{Owner: owner,
			ClassId: strconv.Itoa(i),
		}
		_, err := srv.CreateClassOwner(f.ctx, expected)
		require.NoError(t, err)
		rst, err := f.keeper.ClassOwner.Get(f.ctx, expected.ClassId)
		require.NoError(t, err)
		require.Equal(t, expected.Owner, rst.Owner)
	}
}

func TestClassOwnerMsgServerUpdate(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)

	owner, err := f.addressCodec.BytesToString([]byte("signerAddr__________________"))
	require.NoError(t, err)

	unauthorizedAddr, err := f.addressCodec.BytesToString([]byte("unauthorizedAddr___________"))
	require.NoError(t, err)

	expected := &types.MsgCreateClassOwner{Owner: owner,
		ClassId: strconv.Itoa(0),
	}
	_, err = srv.CreateClassOwner(f.ctx, expected)
	require.NoError(t, err)

	tests := []struct {
		desc    string
		request *types.MsgUpdateClassOwner
		err     error
	}{
		{
			desc: "invalid address",
			request: &types.MsgUpdateClassOwner{Owner: "invalid",
				ClassId: strconv.Itoa(0),
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			desc: "unauthorized",
			request: &types.MsgUpdateClassOwner{Owner: unauthorizedAddr,
				ClassId: strconv.Itoa(0),
			},
			err: sdkerrors.ErrUnauthorized,
		},
		{
			desc: "key not found",
			request: &types.MsgUpdateClassOwner{Owner: owner,
				ClassId: strconv.Itoa(100000),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
		{
			desc: "completed",
			request: &types.MsgUpdateClassOwner{Owner: owner,
				ClassId: strconv.Itoa(0),
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			_, err = srv.UpdateClassOwner(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				rst, err := f.keeper.ClassOwner.Get(f.ctx, expected.ClassId)
				require.NoError(t, err)
				require.Equal(t, expected.Owner, rst.Owner)
			}
		})
	}
}

func TestClassOwnerMsgServerDelete(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)

	owner, err := f.addressCodec.BytesToString([]byte("signerAddr__________________"))
	require.NoError(t, err)

	unauthorizedAddr, err := f.addressCodec.BytesToString([]byte("unauthorizedAddr___________"))
	require.NoError(t, err)

	_, err = srv.CreateClassOwner(f.ctx, &types.MsgCreateClassOwner{Owner: owner,
		ClassId: strconv.Itoa(0),
	})
	require.NoError(t, err)

	tests := []struct {
		desc    string
		request *types.MsgDeleteClassOwner
		err     error
	}{
		{
			desc: "invalid address",
			request: &types.MsgDeleteClassOwner{Owner: "invalid",
				ClassId: strconv.Itoa(0),
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			desc: "unauthorized",
			request: &types.MsgDeleteClassOwner{Owner: unauthorizedAddr,
				ClassId: strconv.Itoa(0),
			},
			err: sdkerrors.ErrUnauthorized,
		},
		{
			desc: "key not found",
			request: &types.MsgDeleteClassOwner{Owner: owner,
				ClassId: strconv.Itoa(100000),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
		{
			desc: "completed",
			request: &types.MsgDeleteClassOwner{Owner: owner,
				ClassId: strconv.Itoa(0),
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			_, err = srv.DeleteClassOwner(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				found, err := f.keeper.ClassOwner.Has(f.ctx, tc.request.ClassId)
				require.NoError(t, err)
				require.False(t, found)
			}
		})
	}
}
