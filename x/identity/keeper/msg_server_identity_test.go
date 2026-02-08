package keeper_test

import (
	"strconv"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"spider/x/identity/keeper"
	"spider/x/identity/types"
)

func TestIdentityMsgServerCreate(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)
	creator, err := f.addressCodec.BytesToString([]byte("signerAddr__________________"))
	require.NoError(t, err)

	for i := 0; i < 5; i++ {
		expected := &types.MsgCreateIdentity{Creator: creator,
			Uid: strconv.Itoa(i),
		}
		_, err := srv.CreateIdentity(f.ctx, expected)
		require.NoError(t, err)
		rst, err := f.keeper.Identity.Get(f.ctx, expected.Uid)
		require.NoError(t, err)
		require.Equal(t, expected.Creator, rst.Creator)
	}
}

func TestIdentityMsgServerUpdate(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)

	creator, err := f.addressCodec.BytesToString([]byte("signerAddr__________________"))
	require.NoError(t, err)

	unauthorizedAddr, err := f.addressCodec.BytesToString([]byte("unauthorizedAddr___________"))
	require.NoError(t, err)

	expected := &types.MsgCreateIdentity{Creator: creator,
		Uid: strconv.Itoa(0),
	}
	_, err = srv.CreateIdentity(f.ctx, expected)
	require.NoError(t, err)

	tests := []struct {
		desc    string
		request *types.MsgUpdateIdentity
		err     error
	}{
		{
			desc: "invalid address",
			request: &types.MsgUpdateIdentity{Creator: "invalid",
				Uid: strconv.Itoa(0),
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			desc: "unauthorized",
			request: &types.MsgUpdateIdentity{Creator: unauthorizedAddr,
				Uid: strconv.Itoa(0),
			},
			err: sdkerrors.ErrUnauthorized,
		},
		{
			desc: "key not found",
			request: &types.MsgUpdateIdentity{Creator: creator,
				Uid: strconv.Itoa(100000),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
		{
			desc: "completed",
			request: &types.MsgUpdateIdentity{Creator: creator,
				Uid: strconv.Itoa(0),
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			_, err = srv.UpdateIdentity(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				rst, err := f.keeper.Identity.Get(f.ctx, expected.Uid)
				require.NoError(t, err)
				require.Equal(t, expected.Creator, rst.Creator)
			}
		})
	}
}

func TestIdentityMsgServerDelete(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)

	creator, err := f.addressCodec.BytesToString([]byte("signerAddr__________________"))
	require.NoError(t, err)

	unauthorizedAddr, err := f.addressCodec.BytesToString([]byte("unauthorizedAddr___________"))
	require.NoError(t, err)

	_, err = srv.CreateIdentity(f.ctx, &types.MsgCreateIdentity{Creator: creator,
		Uid: strconv.Itoa(0),
	})
	require.NoError(t, err)

	tests := []struct {
		desc    string
		request *types.MsgDeleteIdentity
		err     error
	}{
		{
			desc: "invalid address",
			request: &types.MsgDeleteIdentity{Creator: "invalid",
				Uid: strconv.Itoa(0),
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			desc: "unauthorized",
			request: &types.MsgDeleteIdentity{Creator: unauthorizedAddr,
				Uid: strconv.Itoa(0),
			},
			err: sdkerrors.ErrUnauthorized,
		},
		{
			desc: "key not found",
			request: &types.MsgDeleteIdentity{Creator: creator,
				Uid: strconv.Itoa(100000),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
		{
			desc: "completed",
			request: &types.MsgDeleteIdentity{Creator: creator,
				Uid: strconv.Itoa(0),
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			_, err = srv.DeleteIdentity(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				found, err := f.keeper.Identity.Has(f.ctx, tc.request.Uid)
				require.NoError(t, err)
				require.False(t, found)
			}
		})
	}
}
