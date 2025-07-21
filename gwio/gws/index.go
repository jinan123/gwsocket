package gws

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/jinan123/gwsocket/gwio/consts"
	"github.com/zishang520/engine.io/v2/types"
	"github.com/zishang520/socket.io/v2/socket"
	"golang.org/x/net/context"
	"net/http"
	"time"
)

type WsConf struct {
	Port              int   `json:"port" dc:"服务端口"`
	PingInterval      int   `json:"ping_interval" dc:"ping间隔 单位毫秒"`
	PingTimeout       int   `json:"ping_timeout" dc:"每次ping超时 单位毫秒"`
	MaxHttpBufferSize int64 `json:"max_http_buffer_size" dc:"http Buffer最大值"`
	ConnectTimeout    int   `json:"connect_timeout" dc:"连接超时 单位毫秒"`
}

func newDefaultConf() WsConf {
	return WsConf{
		Port:              3400,
		PingInterval:      3000,
		PingTimeout:       5000,
		MaxHttpBufferSize: 1000000,
		ConnectTimeout:    5000,
	}
}

// NewSocketIoHandle websocket服务
func NewSocketIoHandle(ctx context.Context, conf ...WsConf) http.Handler {
	var err error
	cnf := newDefaultConf()
	if len(conf) > 0 {
		cnf = conf[0]
	}
	c := socket.DefaultServerOptions()
	c.SetServeClient(true)
	c.SetPingInterval(time.Duration(cnf.PingInterval) * time.Millisecond)
	c.SetPingTimeout(time.Duration(cnf.PingTimeout) * time.Millisecond)
	c.SetMaxHttpBufferSize(cnf.MaxHttpBufferSize)
	c.SetConnectTimeout(time.Duration(cnf.ConnectTimeout) * time.Millisecond)
	c.SetCors(&types.Cors{Origin: "*", Credentials: true})
	socketIo := socket.NewServer(nil, nil)
	err = socketIo.On(consts.WsConnection, func(clients ...interface{}) {
		if len(clients) > 0 {
			ctx = gctx.New()
			conn := clients[0].(*socket.Socket)
			g.Dump("连接成功....")
			h := DefaultWsEventHandles()
			for _, f := range h {
				err = f(ctx, conn)
				if err != nil {
					g.Log().Error(ctx, err)
				}
			}
		}
	})
	if err != nil {
		g.Log().Error(gctx.New(), err)
	}

	// 返回socketio服务
	return socketIo.ServeHandler(c)
}
