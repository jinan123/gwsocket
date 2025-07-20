package core

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/zishang520/socket.io/v2/socket"
	"golang.org/x/net/context"
	"runtime"
)

var defaultHook = websocketin.WsHook{
	BeforeHandle: func(ctx context.Context, eventName string, msg websocketin.WsNotificationMsg) (int, bool) {
		isPush := 0
		if msg.JobId == 0 {
			isPush = 1
		}
		if eventName == sysws.AdminMemberOnline || eventName == sysws.AdminEmJobStatus {
			isPush = 1
			return 0, true
		}
		id := wsmsgrepo.WriteWsMsg(ctx, adminin.WsMsg{
			OrgId:      msg.OrgId,
			EmId:       msg.EmId,
			JobId:      msg.JobId,
			EventName:  eventName,
			ToTerminal: msg.ToTerminal,
			Contents:   msg.Content,
			IsPush:     isPush,
		})
		return id, true
	},
	AfterHandle: func(ctx context.Context, msg websocketin.WsNotificationMsg) {
		if msg.WsMsgId > 0 {
			wsmsgrepo.SetMsgIsPull(ctx, msg.WsMsgId)
		}
	},
}

func doSendBefore(ctx context.Context, event string, msg websocketin.WsNotificationMsg, wh ...websocketin.WsHook) bool {
	success := true
	if defaultHook.BeforeHandle != nil && msg.WsMsgId == 0 {
		msg.WsMsgId, success = defaultHook.BeforeHandle(ctx, event, msg)
	}
	if success == false {
		return success
	}
	if len(wh) > 0 {
		for _, hook := range wh {
			if hook.BeforeHandle != nil {
				_, success = hook.BeforeHandle(ctx, event, msg)
				if success == false {
					return success
				}
			}
		}
	}
	return success
}

func doSendAfter(ctx context.Context, msg websocketin.WsNotificationMsg, wh ...websocketin.WsHook) {
	//只要有一次成功，则认为是消息已经正确投递
	if defaultHook.AfterHandle != nil && msg.WsMsgId > 0 {
		defaultHook.AfterHandle(gctx.New(), msg)
	}
	if len(wh) > 0 {
		for _, hook := range wh {
			if hook.BeforeHandle != nil {
				hook.AfterHandle(ctx, msg)
			}
		}
	}
}

func Send(ctx context.Context, event string, msg websocketin.WsNotificationMsg) {
	doSendBefore(ctx, event, msg)
	var conn = &socket.Socket{}
	var err error
	if msg.ToTerminal == 1 {
		conn, err = GetEquipmentSocketConn(msg.EmId)
		if err != nil {
			return
		}
	} else {
		conn, err = GetAdminSocketConn(msg.OrgId)
		if err != nil {
			return
		}
	}
	//printCallStack()
	err = conn.Emit(event, gconv.String(msg), func(d []any, ackErr error) {
		if ackErr != nil {
			g.Log().Warning(ctx, ackErr)
			Offline(ctx, conn)
		} else {
			//只要有一次成功，则认为是消息已经正确投递
			doSendAfter(ctx, msg)
		}
	})

	if err != nil {
		g.Log().Error(ctx, err)
	}

}

// 堆栈打印，用于debug
func printCallStack() {
	pcs := make([]uintptr, 10)
	n := runtime.Callers(0, pcs)
	frames := runtime.CallersFrames(pcs[:n])

	for {
		frame, more := frames.Next()
		fmt.Printf("- 函数: %s\n  文件: %s:%d\n", frame.Function, frame.File, frame.Line)
		if !more {
			break
		}
	}
}
