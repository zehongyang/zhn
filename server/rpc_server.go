package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"gopros/proto/pb"
	"log"
	"net"
)

type Server struct {
	
}

func (s *Server) SayHello (ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error)  {
	fmt.Println("name",req.Name)
	return &pb.HelloReply{Message: req.Name+"nihao"},nil
}

func main2()  {
	ln, err := net.Listen("tcp", ":8001")
	if err != nil {
		log.Fatal(err)
	}
	srv := grpc.NewServer()
	pb.RegisterGreeterServer(srv,&Server{})
	srv.Serve(ln)
}