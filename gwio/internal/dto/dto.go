package dto

import (
	"github.com/zishang520/socket.io/v2/socket"
)

// 消息体
type WsMsg interface {
	GetWsMsgId() int
	GetEventName() string
	GetFormId() int
	GetToId() int
	GetData() interface{}
}

// 用户上下文
type UserCtx interface {
	GetSocketConnId() socket.SocketId
	GetUserId() int
}

type GwsUserCtx struct {
	SocketConnId socket.SocketId `json:"socket_conn_id"`
	OrgId        int             `json:"org_id"`
	UserId       int             `json:"user_id"`
}

func (this *GwsUserCtx) GetSocketConnId() socket.SocketId {
	return this.SocketConnId
}
func (this *GwsUserCtx) GetUserId() int {
	return this.UserId
}

type GwsWsMsg struct {
	EventName string `json:"event_name" dc:"事件名称"`
	FromId    int    `json:"from_id" dc:"发出者ID"`
	ToId      int    `json:"to_terminal" dc:"接收者"`
	Content   string `json:"content" dc:"消息内容"`
	WsMsgId   int    `json:"ws_msg_id" dc:"消息Id"`
}

func (this *GwsWsMsg) GetWsMsgId() int {
	return this.WsMsgId
}
func (this *GwsWsMsg) GetEventName() string {
	return this.EventName
}
func (this *GwsWsMsg) GetToId() int {
	return this.ToId
}
func (this *GwsWsMsg) GetContent() string {
	return this.Content
}
func (this *GwsWsMsg) GetFromId() int {
	return this.FromId
}
