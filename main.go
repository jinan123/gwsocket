package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/jinan123/gwsocket/gwio/gws"
)

func main() {
	var ctx = gctx.New()
	s := g.Server()
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
