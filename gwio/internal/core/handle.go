package core

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/jinan123/gwsocket/gwio/internal/consts"
	"github.com/zishang520/socket.io/v2/socket"
	"golang.org/x/net/context"
)

// 事件处理
type WsEventHandle func(ctx context.Context, conn *socket.Socket) error

// ack处理
type AckHandle func([]any, error)

func DefaultWsEventHandles() []WsEventHandle {
	return []WsEventHandle{}
}

// Ping 用户心跳(主动)
func Ping(ctx context.Context, conn *socket.Socket) error {
	err := conn.On(consts.Ping, func(dataArray ...any) {
		if len(dataArray) == 0 {
			return
		}
		ackInterface := dataArray[len(dataArray)-1]
		if ackInterface == nil {
			return
		}
		ack := ackInterface.(AckHandle)

		wsInfo, err := GetWxCtx(conn)
		if err != nil {
			ack([]any{gconv.String(websocketin.UserCtx{})}, err)
		}
		ack([]any{gconv.String(wsInfo)}, nil)
	})
	if err != nil {
		g.Log().Warning(ctx, err)
	}
	return err
}
