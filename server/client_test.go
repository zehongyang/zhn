package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopros/proto/pb"
	"gopros/share/config"
	"gopros/share/logger"
	"net"
	"os"
	"testing"
)

func TestClient(t *testing.T)  {
	conn, err := net.Dial("tcp", ":8000")
	if err != nil {
		logger.F(zap.Error(err))
	}
	n, err := conn.Write([]byte("nihao"))
	fmt.Println(err,n)
}

func TestChan(t *testing.T)  {
	var ch = make(chan int)
	close(ch)
}

func TestEnv(t *testing.T)  {
	data := os.Getenv(config.YmlSource)
	t.Log(data)
}

func TestRpc(t *testing.T)  {
	conn, err := grpc.Dial(":8001",grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	cli := pb.NewGreeterClient(conn)
	res, err := cli.SayHello(context.Background(), &pb.HelloRequest{Name: "zhangsan"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res.Message)
}