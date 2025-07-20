package gsmiddleware

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/jinan123/gwsocket/gwio/internal/dto"
	"github.com/zishang520/socket.io/v2/socket"
	"sync"
)

// Middleware 定义中间件类型
type Middleware func(ctx context.Context, next func())

// Handler 处理函数
type Handler func(ctx context.Context, s socket.SocketId, msg dto.WsMsg)

// Event 事件结构体，包含中间件和最终处理函数
type Event struct {
	middlewares []Middleware
	handler     Handler
}

// Lifecycle 结构体，管理多个事件
type Lifecycle struct {
	mu     sync.RWMutex
	name   string
	events map[string]*Event
}

// NewLifecycle 创建生命周期管理器
func NewLifecycle(name string) *Lifecycle {
	return &Lifecycle{
		events: make(map[string]*Event),
		name:   name,
	}
}

// Register 注册事件及处理函数
func (this *Lifecycle) Register(eventName string, handler Handler, middlewares ...Middleware) {
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
func (this *Lifecycle) Trigger(ctx context.Context, s socket.SocketId, msg dto.WsMsg) {
	this.mu.RLock()
	event, exists := this.events[msg.GetEventName()]
	this.mu.RUnlock()

	if !exists {
		g.Log().Warning(ctx, fmt.Sprintf("Event not found: %s", msg.GetEventName()))
		return
	}

	err := g.Try(ctx, func(ctx context.Context) {
		// 洋葱模型执行中间件
		var exec func(index int)
		exec = func(index int) {
			if index < len(event.middlewares) {
				event.middlewares[index](ctx, func() { exec(index + 1) })
			} else {
				// 所有中间件执行完，调用最终的 Handler
				event.handler(ctx, s, msg)
			}
		}
		exec(0)
	})

	if err != nil {
		g.Log().Warning(ctx, err)
	}
}

func (this *Lifecycle) loadEvents() map[string]*Event {
	this.mu.RLock()
	defer this.mu.RUnlock()
	return this.events
}

// RunAll 主动触发所有事件
func (this *Lifecycle) RunAll(ctx context.Context, s socket.SocketId, msg dto.WsMsg) {
	evs := this.loadEvents()
	for eventName := range evs {
		go func(evName string) {
			this.Trigger(ctx, s, msg)
		}(eventName)
	}
	g.Log().Infof(ctx, "%s All events started", this.GetName())
}

func (this *Lifecycle) GetName() string {
	return this.name
}
