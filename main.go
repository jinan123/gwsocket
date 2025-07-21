package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/jinan123/gwsocket/gwio/dto"
	"github.com/jinan123/gwsocket/gwio/gws"
	"github.com/zishang520/socket.io/v2/socket"
	"golang.org/x/net/context"
)

func main() {
	var ctx = gctx.New()
	s := g.Server()
	var TestHandle = func(ctx context.Context, conn *socket.Socket, msg dto.WsMsg) error {
		g.Dump("获取的ID：" + conn.Id())
		g.Dump("获得数据" + gconv.String(msg.GetData()))
		return nil
	}
	gws.BindOn("gwsTest", TestHandle)
	// /socket.io/
	socketIoHandle := gws.NewSocketIoHandle(ctx)
	s.BindHandler("/socket.io/", func(r *ghttp.Request) {
		socketIoHandle.ServeHTTP(r.Response.Writer, r.Request)
	})
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Write("哈喽世界！")
	})
	s.SetPort(1996)
	s.Run()
}
