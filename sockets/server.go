package sockets

import (
	"go.uber.org/zap"
	"gopros/share/config"
	"gopros/share/errs"
	"gopros/share/logger"
	"net"
	"os"
	"sync"
)


const (
	TcpMode = iota
	WebSocketMode
)


type IServer interface {
	IHandle
	Serve() error
}


type TcpServer struct {
	Config *TcpConfig
	*HandlerManager
	sm *SessionManager
}

func (s *TcpServer) Serve () error  {
	if len(s.Config.Addr) < 1 {
		logger.Error("TcpServer serve",zap.Any("cfg",s.Config))
		return errs.ErrAddr
	}
	if s.HandlerManager == nil || s.sm == nil {
		logger.Error("TcpServer serve",zap.Any("s",s))
		return errs.ErrInvalid
	}
	ln, err := net.Listen("tcp", s.Config.Addr)
	if err != nil {
		return err
	}
	logger.Info("listening on tcp server",zap.Any("port",s.Config.Addr))
	for  {
		conn, err := ln.Accept()
		if err != nil {
			logger.E(zap.Error(err))
			continue
		}
		logger.Debug("accept client",zap.Any("remote addr",conn.RemoteAddr()))
		sess := &Session{Conn:conn,hm: s.HandlerManager,sm: s.sm,data: make(chan []byte,ChanSize)}
		go sess.Handle()
	}
}


type WebsocketServer struct {
	Config *WebsocketConfig
	*HandlerManager
	sm *SessionManager
}

func (s *WebsocketServer) Serve () error {
	return nil
}

var GetSocketServer = func() func(modes ...int) IServer {
	scfg := GetSocketConfig()
	return func(modes ...int) IServer {
		var mode = TcpMode
		if len(modes) > 0 {
			mode = modes[0]
		}
		switch mode {
		case WebSocketMode:
			return &WebsocketServer{Config: scfg.Im.Ws,HandlerManager: GetHandlerManager(),sm: GetSessionManager()}
		default:
			return &TcpServer{Config:scfg.Im.Tcp,HandlerManager: GetHandlerManager(),sm: GetSessionManager()}
		}
	}
}()

type SocketConfig struct {
	Im struct{
		Ws  *WebsocketConfig
		Tcp  *TcpConfig
		Hosts []*HostsConfig
	}
}


var GetSocketConfig = func() func() *SocketConfig {
	var (
		once sync.Once
		s SocketConfig
	)
	return func() *SocketConfig {
		once.Do(func() {
			err := config.Load(&s)
			if err != nil {
				logger.F(zap.Error(err))
			}
		})
		return &s
	}
}()

type TcpConfig struct {
	Addr string
}

type WebsocketConfig struct {
	Addr string
}

type HostsConfig struct {
	ServerName string `yaml:"serverName"`
	Host string `yaml:"host"`
	Rpc string `yaml:"rpc"`
}

func (s *HostsConfig) isLocal () bool {
	return s.ServerName == os.Getenv(config.ServerName)
}

func (s *SocketConfig) CalcNode (uid int64) *HostsConfig {
	if len(s.Im.Hosts) < 1 {
		return nil
	}
	return s.Im.Hosts[uid%int64(len(s.Im.Hosts))]
}