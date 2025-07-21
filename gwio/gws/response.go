package gws

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/jinan123/gwsocket/gwio/consts"
)

type Response struct {
	Code int         `json:"code" dc:"状态码"`
	Msg  string      `json:"msg" dc:"提示信息"`
	Data interface{} `json:"data" dc:"返回数据"`
}

func AckSuccess(ackHandle AckHandle, data ...interface{}) {
	resData := interface{}(g.Map{})
	if len(data) > 0 {
		resData = data[0]
	}
	resp := Response{Code: consts.Success, Msg: "操作成功", Data: resData}
	ackHandle([]any{(resp)}, nil)
}

func AckError(ackHandle AckHandle, err error) {
	g.Log().Warning(gctx.New(), err)
	resp := Response{Code: consts.Error, Msg: "操作失败"}
	ackHandle([]any{resp}, nil)
}

func AckUnauthorized(ackHandle AckHandle) {
	resp := Response{Code: consts.Unauthorized, Msg: "未授权"}
	ackHandle([]any{resp}, nil)
}
