package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"spider/cmd/knot/knot"
	identitytypes "spider/x/identity/types"
	loantypes "spider/x/loan/types"
	officialtypes "spider/x/official/types"
	tokenfactorytypes "spider/x/tokenfactory/types"

	txv1beta1 "cosmossdk.io/api/cosmos/tx/v1beta1"
	abci "github.com/cometbft/cometbft/abci/types"
	logger "github.com/cometbft/cometbft/libs/log"
	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	coretypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/gogoproto/proto"
)

const (
	wsEndpoint = "tcp://127.0.0.1:26657" // 注意：这里是 RPC 地址（http/tcp），WS 由客户端内部走 /websocket
	subscriber = "test-listener"
	eventType  = "my_module.my_event"
)

// ==============================
// handlers
func MsgCreateOperator(ctx context.Context, evt *officialtypes.MsgCreateOperator) error {
	log.Println(evt)
	return nil
}
func MsgCreateDenom(ctx context.Context, evt *tokenfactorytypes.MsgCreateDenom) error {
	log.Println(evt)
	return nil
}
func MsgCreateIdentity(ctx context.Context, evt *identitytypes.MsgCreateIdentity) error {
	log.Println(evt)
	return nil
}
func MsgRequestLoan(ctx context.Context, evt *loantypes.MsgRequestLoan) error {
	log.Println(evt)
	return nil
}

// ------------------------------

func RegisterHandlers() {
	knot.TxEventsRegister(knot.NewGenericHandler(MsgCreateDenom))
	knot.TxEventsRegister(knot.NewGenericHandler(MsgRequestLoan))
	knot.TxEventsRegister(knot.NewGenericHandler(MsgCreateIdentity))
	knot.TxEventsRegister(knot.NewGenericHandler(MsgCreateOperator))
}

// ==============================

func main() {
	log.SetFlags(11)

	ctx := context.Background()

	//注册处理函数
	RegisterHandlers()

	//客户端
	conn, err := rpchttp.New(wsEndpoint, "/websocket")
	if err != nil {
		log.Fatal(err)
	}
	if err := conn.Start(); err != nil {
		log.Fatal(err)
	}
	defer conn.Stop()

	lw := logger.NewTMLogger(os.Stdout)
	conn.SetLogger(lw)

	st, _ := conn.Status(ctx)
	log.Printf("node=%s height=%d", st.NodeInfo.Moniker, st.SyncInfo.LatestBlockHeight)

	// 恢复测试
	catchUpByHeight(ctx, conn, 18839, st.SyncInfo.LatestBlockHeight)

	// 订阅所有 Tx；然后在客户端侧过滤自定义事件（最稳）
	query := "tm.event='Tx'"

	out, err := conn.Subscribe(ctx, subscriber, query)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("subscribed:", query)

	//订阅事件
	for {
		select {
		case msg := <-out:
			ev, ok := msg.Data.(coretypes.EventDataTx)
			if !ok {
				log.Println("event failed?")
				continue
			}
			// v0.50+：TxResult 里有 ABCI events
			handleTx(ev.TxResult)
		case <-ctx.Done():
			return
		}
	}
}

// ------------------------------------------------

//------------------------------------------------

func handleTx(ev abci.TxResult) error {
	if ev.Result.Code != 0 {
		log.Println(ev.Result.Log)
		return nil
	}
	log.Printf("=======tx> height:%v, index:%v", ev.Height, ev.Index)
	var tx txv1beta1.Tx
	if err := proto.Unmarshal(ev.Tx, &tx); err != nil {
		return err
	}
	if tx.Body == nil || len(tx.Body.Messages) < 1 {
		return nil
	}
	ctx := context.Background()
	// log.Println(tx)
	for _, mi := range tx.Body.Messages {
		err := knot.TxEventsDispatch(ctx, mi)
		if err != nil {
			log.Println(err, mi)
			continue
		}
	}
	return nil
}

func attrsToMap(kvs []abci.EventAttribute) map[string]string {
	m := make(map[string]string, len(kvs))
	for _, kv := range kvs {
		// kv.Key/Value 是 []byte
		m[string(kv.Key)] = string(kv.Value)
	}
	return m
}

func catchUpByHeight(ctx context.Context, c *rpchttp.HTTP, from, to int64) error {
	const perPage = 100
	log.Println("catchUpByHeight")
	for h := from; h <= to; h++ {
		page := 1
		for {
			q := fmt.Sprintf("tx.height=%d", h)
			res, err := c.TxSearch(ctx, q, false, &page, ptrInt(perPage), "")
			if err != nil {
				log.Println(err)
				return err
			}

			// log.Println("catchUpByHeight", h, len(res.Txs), perPage)
			for _, ev := range res.Txs {
				// tx 是 *coretypes.ResultTx，里面有 TxResult（含 events）
				handleTx(abci.TxResult{
					Height: ev.Height,
					Index:  ev.Index,
					Tx:     ev.Tx,
					Result: ev.TxResult,
				})
			}

			if len(res.Txs) < perPage {
				break
			}
			page++
		}
	}
	return nil
}

func ptrInt(v int) *int { return &v }
