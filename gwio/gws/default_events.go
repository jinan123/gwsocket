package gws

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/jinan123/gwsocket/gwio/consts"
	"github.com/zishang520/socket.io/v2/socket"
)

func DefaultWsEventHandles() []WsEventHandle {
	return []WsEventHandle{
		Ping,
	}
}

// Ping 用户心跳
func Ping(ctx context.Context, conn *socket.Socket) (err error) {
	err = conn.On(consts.Ping, func(dataArray ...any) {
		ack, err := loadAckHandle(dataArray...)
		if err != nil {
			g.Log().Error(ctx, err)
			return
		}
		wsInfo, err := GetWxCtx(conn)
		if err != nil {
			AckError(ack, err)
			return
		}
		AckSuccess(ack, wsInfo)
	})
	return
}
