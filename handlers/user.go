package handlers

import (
	"go.uber.org/zap"
	"gopros/proto/pb"
	"gopros/share/logger"
	"gopros/sockets"
	"net/http"
)

func UserLoginQuery() sockets.HandleFunc {
	return func(ctx *sockets.SessionContext) {
		var q pb.UserLoginQuery
		var res pb.UserLoginQueryResponse
		err := ctx.ShouldBind(&q)
		if err != nil {
			logger.E(zap.Error(err))
			ctx.Response(http.StatusInternalServerError)
			return
		}
		if q.Uid < 1 || len(q.Token) < 1 || q.Token != "123456" {
			logger.E(zap.Any("q",q))
			ctx.Response(http.StatusBadRequest)
			return
		}
		err = ctx.StoreSession(q.Uid)
		if err != nil {
			logger.E(zap.Error(err),zap.Any("q",q))
			ctx.Response(http.StatusInternalServerError)
			return
		}
		res.N = 1
		ctx.ResponseMsg(http.StatusOK,&res)
		return
	}
}