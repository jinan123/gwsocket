package gws

import (
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/jinan123/gwsocket/gwio/dto"
	"github.com/zishang520/socket.io/v2/socket"
	"golang.org/x/net/context"
	"sync"
)

var GwsEventCtrl = NewGwsEvent()

// Middleware 定义中间件类型
type Middleware func(ctx context.Context, next func())

// Handler 处理函数
type Handler func(ctx context.Context, conn *socket.Socket, msg dto.WsMsg) error

// Event 事件结构体，包含中间件和最终处理函数
type Event struct {
	middlewares []Middleware
	handler     Handler
}

// GwsEvent 结构体，管理多个事件
type GwsEvent struct {
	mu     sync.RWMutex
	events map[string]*Event
}

func NewGwsEvent() *GwsEvent {
	return &GwsEvent{
		events: make(map[string]*Event),
	}
}

// Register 注册事件及处理函数
func (this *GwsEvent) Register(eventName string, handler Handler, middlewares ...Middleware) {
	this.mu.Lock()
	defer this.mu.Unlock()
	_, exists := this.events[eventName]
	if exists {
		g.Log().Warningf(gctx.New(), "The event already exists. Do not re-register event:%s", eventName)
		return
	}
	this.events[eventName] = &Event{
		handler:     handler,
		middlewares: middlewares,
	}
}

// Trigger 触发事件，执行中间件和处理函数（洋葱模型）
func (this *GwsEvent) Trigger(ctx context.Context, conn *socket.Socket, msg dto.WsMsg) (err error) {
	this.mu.RLock()
	event, exists := this.events[msg.GetEventName()]
	this.mu.RUnlock()

	if !exists {
		return gerror.New(fmt.Sprintf("Event not found: %s", msg.GetEventName()))
	}

	// 洋葱模型执行中间件
	var exec func(index int)
	exec = func(index int) {
		if index < len(event.middlewares) {
			event.middlewares[index](ctx, func() { exec(index + 1) })
		} else {
			// 所有中间件执行完，调用最终的 Handler
			err = event.handler(ctx, conn, msg)
		}
	}
	exec(0)
	return
}

func (this *GwsEvent) loadEvents() map[string]*Event {
	this.mu.RLock()
	defer this.mu.RUnlock()
	return this.events
}

// RunAll 主动触发所有事件
func (this *GwsEvent) RunAll(ctx context.Context, conn *socket.Socket, msg dto.WsMsg) {
	evs := this.loadEvents()
	for eventName := range evs {
		go func(evName string) {
			err := this.Trigger(ctx, conn, msg)
			if err != nil {
				g.Log().Warning(ctx, err)
			}
		}(eventName)
	}
	g.Log().Infof(ctx, "%s All events started")
}
