package sockets

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"gopros/proto/pb"
	"gopros/share/errs"
	"gopros/share/logger"
	"net"
	"reflect"
	"sync"
)

type RpcServer struct {
	hmp map[pb.HandleId]RpcHandleFunc
	config *SocketConfig
}


type RpcHandle struct {
	Hid pb.HandleId
	Fun RpcHandleFunc
}

func (s *RpcServer) Handle (ctx context.Context,req  *pb.ServerRequest) (*pb.ServerResponse,error)  {
	handleFunc,ok := s.hmp[req.Hid]
	if !ok {
		logger.Error("RpcHandle",zap.Any("req",req))
		return nil,errs.ErrHandle
	}
	rctx := &RpcContext{Uid: req.Uid,Hid: req.Hid,Dt: req.Dt,res: &pb.ServerResponse{},body: req.Body}
	handleFunc(rctx)
	return rctx.response()
}


func (s *RpcServer) Register (hds ...RpcHandle)  {
	for _, hd := range hds {
		if _,ok := s.hmp[hd.Hid];!ok{
			s.hmp[hd.Hid] = hd.Fun
		}
	}
}

func (s *RpcServer) Serve () error {
	var th *HostsConfig
	for _, host := range s.config.Im.Hosts {
		if host.isLocal() {
			th = host
		}
	}
	if th == nil || len(th.Rpc) < 1 {
		logger.Error("RpcServer",zap.Any("host",th))
		return errs.ErrRpcHost
	}
	ln, err := net.Listen("tcp", th.Rpc)
	if err != nil {
		logger.Error("RpcServer",zap.Any("host",th),zap.Error(err))
		return err
	}
	srv := grpc.NewServer()
	pb.RegisterServerServer(srv,s)
	logger.Info("rpc listening on:",zap.Any("addr",th.Rpc))
	return srv.Serve(ln)
}

var GetRpcServer = func() func() *RpcServer {
	var (
		once sync.Once
		s *RpcServer
	)
	return func() *RpcServer {
		once.Do(func() {
			s = &RpcServer{hmp: make(map[pb.HandleId]RpcHandleFunc),config: GetSocketConfig()}
		})
		return s
	}
}()

type RpcHandleFunc func(ctx *RpcContext)


type RpcClient struct {
	client pb.ServerClient
	conn *grpc.ClientConn
	rpcServer *RpcServer
	isLocal bool
}

func NewRpcClient(config *HostsConfig)(*RpcClient,error) {
	if config == nil {
		return nil,errs.ErrConfig
	}
	if len(config.Rpc) < 1 {
		return nil,errs.ErrRpcHost
	}
	var rc RpcClient
	if !config.isLocal() {
		conn, err := grpc.Dial(config.Rpc,grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil,err
		}
		cli := pb.NewServerClient(conn)
		rc.client = cli
		rc.conn = conn
	}
	rc.rpcServer = GetRpcServer()
	rc.isLocal = config.isLocal()
	return &rc,nil
}

func (s *RpcClient) Handle (fid int64,msg proto.Message) (*pb.ServerResponse,error) {
	var req pb.ServerRequest
	rt := reflect.TypeOf(msg)
	var name string
	if rt.Kind() == reflect.Ptr {
		name = rt.Elem().Name()
	}else{
		name = rt.Name()
	}
	hid,ok := pb.HandleId_value[fmt.Sprintf("HI_%s",name)]
	if !ok {
		return nil,errs.ErrHandle
	}
	req.Hid = pb.HandleId(hid)
	req.Uid = fid
	md, err := proto.Marshal(msg)
	if err != nil {
		return nil,err
	}
	req.Body = md
	req.Dt = pb.DataType_DT_Proto
	if !s.isLocal {
		defer s.conn.Close()
		return s.client.Handle(context.Background(),&req)
	}
	return s.rpcServer.Handle(context.Background(),&req)
}


type RpcContext struct {
	res *pb.ServerResponse
	Uid int64
	Hid pb.HandleId
	Dt pb.DataType
	err error
	body []byte
}

func (s *RpcContext) ShouldBind (out proto.Message) error {
	if len(s.body) < 1 {
		return errs.ErrBody
	}
	switch s.Dt {
	case pb.DataType_DT_Proto:
		return proto.Unmarshal(s.body,out)
	case pb.DataType_DT_Json:
		return json.Unmarshal(s.body,out)
	default:
		return errs.ErrDataType
	}
}

func (s *RpcContext) ResponseCode (code int,errs ...error)  {
	s.res.Code = int32(code)
	if len(errs) > 0 {
		s.err = errs[0]
	}
}

func (s *RpcContext) ResponseMsg (code int,msg proto.Message)  {
	s.res.Code = int32(code)
	switch s.Dt {
	case pb.DataType_DT_Proto:
		s.res.Data, s.err = proto.Marshal(msg)
	case pb.DataType_DT_Json:
		s.res.Data, s.err = json.Marshal(msg)
	}
}

func (s *RpcContext) response () (*pb.ServerResponse,error) {
	return s.res,s.err
}