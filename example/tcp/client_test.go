package main

import (
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"gopros/handlers/rpc"
	"gopros/proto/pb"
	"gopros/share/config"
	"gopros/share/logger"
	"gopros/sockets"
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
	socketConfig := sockets.GetSocketConfig()
	for _, host := range socketConfig.Im.Hosts {
		t.Log(host)
	}
	host, port, err := net.SplitHostPort("127.0.0.1:18002")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(host,port)
}

func TestEnv(t *testing.T)  {
	data := os.Getenv(config.ServerName)
	t.Log(data)
}

func TestRpc(t *testing.T)  {
	srv := sockets.GetRpcServer()
	srv.Register([]sockets.RpcHandle{{Hid: pb.HandleId_HI_UserLoginQuery,Fun:rpc.UserLoginQuery()}}...)
	cfg := sockets.GetSocketConfig()
	node := cfg.CalcNode(3)
	client, err := sockets.NewRpcClient(node)
	if err != nil {
		t.Fatal(err)
	}
	res, err := client.Handle(3, &pb.UserLoginQuery{Uid: 3, Token: "5141654"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res.Code)
	var tres pb.UserLoginQueryResponse
	err = proto.Unmarshal(res.Data, &tres)
	t.Log(err,tres.N)
}
