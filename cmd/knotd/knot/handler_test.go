package knot

import (
	"context"
	"spider/x/tokenfactory/types"
	"testing"

	"google.golang.org/protobuf/types/known/anypb"
)

func TestHandler(t *testing.T) {
	txEventsHandlers.Register(NewGenericHandler(func(context.Context, *types.MsgCreateDenom) error { return nil }))

	txEventsHandlers.Dispatch(context.Background(), &anypb.Any{})
}
