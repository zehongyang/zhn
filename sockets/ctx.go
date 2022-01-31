package sockets

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"gopros/proto/pb"
	"gopros/share/errs"
	"gopros/share/logger"
)

type SessionContext struct {
	context.Context
	Request []byte
	Body []byte
	Dt pb.DataType
	session *Session
	sequence int64
	hid pb.HandleId
}

func newSessionContext() *SessionContext {
	return &SessionContext{Context:context.Background()}
}

func (s *SessionContext) parseRequest (req proto.Message) error {
	if len(s.Request) < 1 {
		return errs.ErrReq
	}
	switch s.Dt {
	case pb.DataType_DT_Proto:
		return proto.Unmarshal(s.Request,req)
	case pb.DataType_DT_Json:
		return json.Unmarshal(s.Request,req)
	default:
		logger.Error("parseRequest",zap.Any("dt",s.Dt))
		return errs.ErrDataType
	}
}

func (s *SessionContext) ShouldBind (out proto.Message) error {
	if len(s.Body) < 1 {
		return errs.ErrBody
	}
	switch s.Dt {
	case pb.DataType_DT_Proto:
		return proto.Unmarshal(s.Body,out)
	case pb.DataType_DT_Json:
		return json.Unmarshal(s.Body,out)
	default:
		logger.Error("ShouldBind",zap.Any("dt",s.Dt))
		return errs.ErrDataType
	}
}


func (s *SessionContext) GetUid () int64 {
	return s.session.GetUid()
}

func (s *SessionContext) Response (code int,errs ...string)  {
	var res pb.MsgRequestQueryResponse
	res.Code = pb.HandleCode(code)
	if len(errs) > 0 {
		res.Error = errs[0]
	}
	res.Sequence = s.sequence
	res.Hid = s.hid
	marshal, err := s.marshal(&res)
	if err != nil {
		logger.E(zap.Error(err))
		return
	}
	s.session.data <- marshal
}

func (s *SessionContext) ResponseMsg (code int,msg proto.Message)  {
	var res pb.MsgRequestQueryResponse
	res.Code = pb.HandleCode(code)
	res.Hid = s.hid
	res.Sequence = s.sequence
	bts, err := s.marshal(msg)
	if err != nil {
		logger.E(zap.Error(err))
		return
	}
	res.Data = bts
	bts, err = s.marshal(&res)
	if err != nil {
		logger.E(zap.Error(err))
		return
	}
	s.session.data <- bts
}


func parseDataLen(length uint32) (uint32,pb.DataType) {
	return length & 0x00FFFFFF,pb.DataType(length & 0xFF000000)
}

func (s *SessionContext) marshal (msg proto.Message) ([]byte,error) {
	switch s.Dt {
	case pb.DataType_DT_Proto:
		return proto.Marshal(msg)
	case pb.DataType_DT_Json:
		return json.Marshal(msg)
	default:
		return nil,errs.ErrDataType
	}
}

func (s *SessionContext) StoreSession (uid int64) error {
	return s.session.Login(uid)
}
