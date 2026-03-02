package keeper_test

import (
	"strconv"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"spider/x/snft/keeper"
	"spider/x/snft/types"
)

func TestClassNamespaceMsgServerCreate(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)
	creator, err := f.addressCodec.BytesToString([]byte("signerAddr__________________"))
	require.NoError(t, err)

	for i := 0; i < 5; i++ {
		expected := &types.MsgCreateClassNamespace{Creator: creator,
			Namespace: strconv.Itoa(i),
		}
		_, err := srv.CreateClassNamespace(f.ctx, expected)
		require.NoError(t, err)
		rst, err := f.keeper.ClassNamespace.Get(f.ctx, expected.Namespace)
		require.NoError(t, err)
		require.Equal(t, expected.Creator, rst.Creator)
	}
}

func TestClassNamespaceMsgServerUpdate(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)

	creator, err := f.addressCodec.BytesToString([]byte("signerAddr__________________"))
	require.NoError(t, err)

	unauthorizedAddr, err := f.addressCodec.BytesToString([]byte("unauthorizedAddr___________"))
	require.NoError(t, err)

	expected := &types.MsgCreateClassNamespace{Creator: creator,
		Namespace: strconv.Itoa(0),
	}
	_, err = srv.CreateClassNamespace(f.ctx, expected)
	require.NoError(t, err)

	tests := []struct {
		desc    string
		request *types.MsgUpdateClassNamespace
		err     error
	}{
		{
			desc: "invalid address",
			request: &types.MsgUpdateClassNamespace{Creator: "invalid",
				Namespace: strconv.Itoa(0),
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			desc: "unauthorized",
			request: &types.MsgUpdateClassNamespace{Creator: unauthorizedAddr,
				Namespace: strconv.Itoa(0),
			},
			err: sdkerrors.ErrUnauthorized,
		},
		{
			desc: "key not found",
			request: &types.MsgUpdateClassNamespace{Creator: creator,
				Namespace: strconv.Itoa(100000),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
		{
			desc: "completed",
			request: &types.MsgUpdateClassNamespace{Creator: creator,
				Namespace: strconv.Itoa(0),
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			_, err = srv.UpdateClassNamespace(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				rst, err := f.keeper.ClassNamespace.Get(f.ctx, expected.Namespace)
				require.NoError(t, err)
				require.Equal(t, expected.Creator, rst.Creator)
			}
		})
	}
}

func TestClassNamespaceMsgServerDelete(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)

	creator, err := f.addressCodec.BytesToString([]byte("signerAddr__________________"))
	require.NoError(t, err)

	unauthorizedAddr, err := f.addressCodec.BytesToString([]byte("unauthorizedAddr___________"))
	require.NoError(t, err)

	_, err = srv.CreateClassNamespace(f.ctx, &types.MsgCreateClassNamespace{Creator: creator,
		Namespace: strconv.Itoa(0),
	})
	require.NoError(t, err)

	tests := []struct {
		desc    string
		request *types.MsgDeleteClassNamespace
		err     error
	}{
		{
			desc: "invalid address",
			request: &types.MsgDeleteClassNamespace{Creator: "invalid",
				Namespace: strconv.Itoa(0),
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			desc: "unauthorized",
			request: &types.MsgDeleteClassNamespace{Creator: unauthorizedAddr,
				Namespace: strconv.Itoa(0),
			},
			err: sdkerrors.ErrUnauthorized,
		},
		{
			desc: "key not found",
			request: &types.MsgDeleteClassNamespace{Creator: creator,
				Namespace: strconv.Itoa(100000),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
		{
			desc: "completed",
			request: &types.MsgDeleteClassNamespace{Creator: creator,
				Namespace: strconv.Itoa(0),
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			_, err = srv.DeleteClassNamespace(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				found, err := f.keeper.ClassNamespace.Has(f.ctx, tc.request.Namespace)
				require.NoError(t, err)
				require.False(t, found)
			}
		})
	}
}
