package rpc

import (
	"go.uber.org/zap"
	"gopros/proto/pb"
	"gopros/share/logger"
	"gopros/sockets"
	"net/http"
)

func UserLoginQuery() sockets.RpcHandleFunc {
	return func(ctx *sockets.RpcContext) {
		var req pb.UserLoginQuery
		err := ctx.ShouldBind(&req)
		if err != nil {
			logger.E(zap.Error(err))
			ctx.ResponseCode(http.StatusInternalServerError,err)
			return
		}
		logger.I(zap.Any("req",req))
		ctx.ResponseMsg(http.StatusOK,&pb.UserLoginQueryResponse{N: 1})
		return
	}
}