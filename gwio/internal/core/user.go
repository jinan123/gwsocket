package core

import (
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/jinan123/gwsocket/gwio/internal/dto"
	"github.com/zishang520/socket.io/v2/socket"
)

// GetCtx 获取连接上下文
func GetWxCtx(conn *socket.Socket) (dto.UserCtx, error) {
	wsCtx := conn.Data()
	if wsCtx != nil {
		return wsCtx.(*dto.GwsUserCtx), nil
	}
	return nil, gerror.New("读取连接上下文失败")
}
