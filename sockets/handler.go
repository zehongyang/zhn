package sockets

import (
	"go.uber.org/zap"
	"gopros/proto/pb"
	"gopros/share/logger"
)

type RegisterHandle struct {
	Hid pb.HandleId
	Fun HandleFunc
}

type IHandle interface {
	Register(hds ...RegisterHandle)
	Handle(hid pb.HandleId) HandleFunc
}


type HandleFunc func(ctx *SessionContext)


type HandlerManager struct {
	mp map[pb.HandleId]HandleFunc
}

func (s *HandlerManager) Register (hds ...RegisterHandle)  {
	for _, hd := range hds {
		if _,ok := s.mp[hd.Hid];!ok{
			logger.Info("Register handler",zap.Any("hid",hd.Hid))
			s.mp[hd.Hid] = hd.Fun
		}
	}
}

func (s *HandlerManager) Handle (hid pb.HandleId) HandleFunc {
	return s.mp[hid]
}

var GetHandlerManager = func() func() *HandlerManager {
	var (
		s *HandlerManager
	)
	return func() *HandlerManager {
		s = &HandlerManager{mp: make(map[pb.HandleId]HandleFunc)}
		return s
	}
}()
