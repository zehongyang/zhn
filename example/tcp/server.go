package main

import (
	"go.uber.org/zap"
	"gopros/handlers"
	"gopros/handlers/rpc"
	"gopros/proto/pb"
	"gopros/share/logger"
	"gopros/sockets"
	"log"
)

func main() {
	//tcp server
	srv := sockets.GetSocketServer()
	srv.Register([]sockets.RegisterHandle{
		{pb.HandleId_HI_UserLoginQuery,handlers.UserLoginQuery()},
	}...)
	go func() {
		//rpc server
		srv := sockets.GetRpcServer()
		srv.Register([]sockets.RpcHandle{{Hid: pb.HandleId_HI_UserLoginQuery,Fun:rpc.UserLoginQuery()}}...)
		err := srv.Serve()
		if err != nil {
			logger.F(zap.Error(err))
		}
	}()
	log.Println(srv.Serve())
}
