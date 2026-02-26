package knot

import (
	"context"
	"spider/x/tokenfactory/types"
	"testing"

	"google.golang.org/protobuf/types/known/anypb"
)

func TestHandler(t *testing.T) {
	Handlers.Register(NewGenericHandler(func(context.Context, *types.MsgCreateDenom) error { return nil }))

	Handlers.Dispatch(context.Background(), &anypb.Any{})
}
