package gws

import (
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/zishang520/socket.io/v2/socket"
	"golang.org/x/net/context"
)

// 事件处理
type WsEventHandle func(ctx context.Context, conn *socket.Socket) error

// ack处理
type AckHandle func([]any, error)

// 获取ack方法，并校验数据有效性
func loadAckHandle(dataArray ...any) (AckHandle, error) {
	if len(dataArray) == 0 {
		return nil, gerror.New("ack is nil!")
	}
	ackInterface, ok := dataArray[len(dataArray)-1].(func([]any, error))
	if ackInterface == nil || !ok {
		return nil, gerror.New("ack handle is nil")
	}
	return ackInterface, nil
}

func BindOn() {

}
