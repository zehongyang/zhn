package main

import (
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"gopros/proto/pb"
	"gopros/share/logger"
	"gopros/share/utils"
	"net"
	"time"
)

func main1()  {
	conn, err := net.Dial("tcp", ":8000")
	if err != nil {
		logger.F(zap.Error(err))
	}
	var req = pb.UserLoginQuery{Uid: 1,Token: "1234567"}
	body, err := proto.Marshal(&req)
	if err != nil {
		logger.F(zap.Error(err))
	}
	var msgReq = pb.MsgRequestQuery{Hid: pb.HandleId_HI_UserLoginQuery,Body: body,Sequence: time.Now().UnixNano()}
	data, err := utils.Pack(&msgReq, pb.DataType_DT_Proto)
	if err != nil {
		logger.F(zap.Error(err))
	}
	n, err := conn.Write(data)
	fmt.Println("write",n,err)
	var buf [1024]byte
	n, err = conn.Read(buf[:])
	fmt.Println("read",n,err)
	var res pb.MsgRequestQueryResponse
	err = proto.Unmarshal(buf[:n], &res)
	if err != nil {
		logger.F(zap.Error(err))
	}
	var tres pb.UserLoginQueryResponse
	fmt.Println("res",res)
	err = proto.Unmarshal(res.Data, &tres)
	if err != nil {
		logger.F(zap.Error(err))
	}
	fmt.Println("tres n",tres.N)
}
