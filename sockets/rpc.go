package sockets

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"gopros/proto/pb"
	"gopros/share/errs"
)

type RpcServer struct {

}


type RpcClient struct {
	client pb.ServerClient
	conn *grpc.ClientConn
}

func NewRpcClient(config *HostsConfig)(*RpcClient,error) {
	if config == nil {
		return nil,errs.ErrConfig
	}
	if len(config.Rpc) < 1 {
		return nil,errs.ErrRpcHost
	}
	conn, err := grpc.Dial(config.Rpc,grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil,err
	}
	cli := pb.NewServerClient(conn)
	return &RpcClient{client: cli,conn: conn},nil
}

func (s *RpcClient) Handle (fid int64,req *pb.ServerRequest,msg proto.Message) (*pb.ServerResponse,error) {
	defer s.conn.Close()
	req.Uid = fid
	md, err := proto.Marshal(msg)
	if err != nil {
		return nil,err
	}
	req.Body = md
	req.Dt = pb.DataType_DT_Proto
	return s.client.Serve(context.Background(),req)
}


type RpcContext struct {
	sm *SessionManager
}