package main

import (
	"gopros/handlers"
	"gopros/proto/pb"
	"gopros/sockets"
	"log"
)

func main()  {
	srv := sockets.GetSocketServer()
	srv.Register([]sockets.RegisterHandle{
		{pb.HandleId_HI_UserLoginQuery,handlers.UserLoginQuery()},
	}...)
	log.Println(srv.Serve())
}