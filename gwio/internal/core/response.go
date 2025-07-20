package core

import "github.com/jinan123/gwsocket/gwio/internal/consts"

func AckSuccess(ackHandle AckHandle, data ...interface{}) {
	ackHandle([]any{consts.Success, "操作成功", data}, nil)
}

func AckError(ackHandle AckHandle, err error) {
	ackHandle([]any{consts.Error, "操作失败"}, err)
}

func AckUnauthorized(ackHandle AckHandle) {
	ackHandle([]any{consts.Unauthorized, "token失效"}, nil)
}
