package knot

import (
	"context"
	"fmt"
	"sync"

	"github.com/cosmos/gogoproto/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// Handler：处理一个 Any。实现里通常会解包成具体类型再处理。
type Handler interface {
	TypeURL() string
	Handle(ctx context.Context, a *anypb.Any) error
}

type Registry struct {
	mu sync.RWMutex
	m  map[string]Handler
}

var Handlers = NewRegistry()

func NewRegistry() *Registry {
	return &Registry{m: make(map[string]Handler)}
}

func (r *Registry) Register(h Handler) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	u := h.TypeURL()
	if u == "" {
		return fmt.Errorf("empty typeUrl")
	}
	if _, ok := r.m[u]; ok {
		return fmt.Errorf("handler already registered for %s", u)
	}
	r.m[u] = h
	return nil
}

func (r *Registry) Dispatch(ctx context.Context, a *anypb.Any) error {
	if a == nil {
		return fmt.Errorf("nil any")
	}

	r.mu.RLock()
	h := r.m[a.GetTypeUrl()]
	r.mu.RUnlock()

	if h == nil {
		return fmt.Errorf("no handler for typeUrl=%s", a.GetTypeUrl())
	}
	return h.Handle(ctx, a)
}

// PB 约束：*T 必须实现 proto.Message（即 T 是某个生成的消息 struct 类型）
type PB[T any] interface {
	*T
	proto.Message
}

type genericHandler[T any, PT PB[T]] struct {
	typeURL string
	fn      func(context.Context, *T) error
}

func (h genericHandler[T, PT]) TypeURL() string { return h.typeURL }

func (h genericHandler[T, PT]) Handle(ctx context.Context, a *anypb.Any) error {
	msg := new(T) // *T
	// 通过约束把 *T 视为 proto.Message
	var pm PT = msg
	err := proto.Unmarshal(a.Value, pm)
	if err != nil {
		return err
	}
	return h.fn(ctx, msg)
}

// 只传 func(ctx, *T) error，自动生成 typeUrl + handler
func NewGenericHandler[T any, PT PB[T]](fn func(context.Context, *T) error) Handler {
	msg := new(T)
	var pm PT = msg

	typeURL := "/" + proto.MessageName(pm)
	return genericHandler[T, PT]{typeURL: typeURL, fn: fn}
}
