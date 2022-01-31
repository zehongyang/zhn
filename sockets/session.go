package sockets

import (
	"context"
	"encoding/binary"
	"go.uber.org/zap"
	"gopros/proto/pb"
	"gopros/share/errs"
	"gopros/share/logger"
	"io"
	"net"
	"sync"
	"time"
)

type SessionManager struct {
	cm *ConcurrentMap
}


var GetSessionManager = func() func(shardsNums ...int) *SessionManager {
	var (
		//once sync.Once
		s *SessionManager
	)
	return func(shardsNums ...int) *SessionManager {
		//once.Do(func() {
			s = newSessionManager(shardsNums...)
		//})
		return s
	}
}()

func newSessionManager(shardsNums ...int) *SessionManager {
	if len(shardsNums) < 1 {
		shardsNums = append(shardsNums,32)
	}
	cm := &ConcurrentMap{}
	for i := 0; i < shardsNums[0]; i++ {
		cm.shards = append(cm.shards,&ItemMap{
			item: make(map[int64]*Session),
			mu:   sync.RWMutex{},
		})
	}
	return &SessionManager{cm: cm}
}

func (s *SessionManager) StoreSession (sess *Session) error {
	if sess == nil {
		return errs.ErrSession
	}
	if sess.uid < 1 {
		return nil
	}
	return s.cm.storeSession(sess)
}

func (s *SessionManager) GetSession (uid int64) (*Session,error) {
	if uid < 1 {
		return nil,nil
	}
	return s.cm.getSession(uid)
}

func (s *SessionManager) RemoveSession (uid int64) error  {
	if uid < 1 {
		return nil
	}
	return s.cm.removeSession(uid)
}

const (
	ChanSize = 100
	DataMaxLength = 4096 //默认最大长度4k
	DataLength = 4 //数据长度字节数
)

var (
	ReadTimeOut = 50 * time.Second
	WriteTimeout = 10 * time.Second
)

type Session struct {
	net.Conn
	uid  int64
	isLogin bool
	hm *HandlerManager
	sm *SessionManager
	data chan []byte
	closed bool
}

func (s *Session) Login (uid int64) error {
	if uid < 1 {
		return errs.ErrSession
	}
	s.uid = uid
	s.isLogin = true
	return s.sm.StoreSession(s)
}


func (s *Session) GetUid () int64 {
	return s.uid
}

func (s *Session) IsLogin () bool {
	return s.isLogin
}

func (s *Session) Handle ()  {
	go s.read()
	go s.write()
}

func (s *Session) read ()  {
	//todo deal packet data
	defer s.close()
	for  {
		err := s.SetReadDeadline(time.Now().Add(ReadTimeOut))
		if err != nil {
			logger.E(zap.Error(err),zap.Any("session",s))
			return
		}
		var dataLen [DataLength]byte
		_, err = io.ReadFull(s.Conn, dataLen[:])
		if err != nil {
			logger.E(zap.Error(err),zap.Any("session",s))
			return
		}
		length := binary.BigEndian.Uint32(dataLen[:])
		dataLength, dataType := parseDataLen(length)
		var data = make([]byte,dataLength)
		n, err := io.ReadFull(s.Conn, data)
		if err != nil {
			logger.E(zap.Error(err),zap.Any("session",s))
			return
		}
		if len(data) < 1 || len(data) > DataMaxLength {
			logger.Error("read data length err",zap.Any("length",len(data)),zap.Any("n",n))
			return
		}
		err = s.dealMsg(data,dataType)
		if err != nil {
			logger.E(zap.Error(err),zap.Any("session",s))
			return
		}
	}
}

func (s *Session) dealMsg (data []byte,dt pb.DataType) error {
	ctx := &SessionContext{Context:context.Background(),Request: data,Dt: dt,session: s}
	var req pb.MsgRequestQuery
	err := ctx.parseRequest(&req)
	if err != nil {
		return err
	}
	if req.Hid != pb.HandleId_HI_UserLoginQuery && !s.isLogin {
		return errs.ErrSession
	}
	ctx.Body = req.Body
	ctx.sequence = req.Sequence
	ctx.hid = req.Hid
	handleFunc,ok := s.hm.mp[req.Hid]
	if !ok {
		logger.Error("dealMsg",zap.Any("req",req))
		return errs.ErrHandle
	}
	handleFunc(ctx)
	return nil
}


func (s *Session) write ()  {
	for  {
		select {
		case data,ok := <- s.data:
			if !ok {
				return
			}
			err := s.SetWriteDeadline(time.Now().Add(WriteTimeout))
			if err != nil {
				logger.E(zap.Error(err))
				return
			}
			n, err := s.Write(data)
			if err != nil {
				logger.E(zap.Error(err),zap.Any("n",n))
				return
			}
		}
	}
}

func (s *Session) close () {
	close(s.data)
	if !s.closed {
		err := s.sm.RemoveSession(s.uid)
		logger.Error("Session close",zap.Error(err))
		s.Close()
	}
}

type ConcurrentMap struct {
	shards []*ItemMap
}

func (s *ConcurrentMap) storeSession (sess *Session) error {
	index := sess.uid % int64(len(s.shards))
	itemMap := s.shards[index]
	if itemMap == nil {
		return errs.ErrNil
	}
	itemMap.mu.Lock()
	if itemMap.item == nil {
		itemMap.item = make(map[int64]*Session)
	}
	tsess,ok := itemMap.item[sess.uid]
	if ok {
		tsess.closed = true
		tsess.Close()
	}
	itemMap.item[sess.uid] = sess
	itemMap.mu.Unlock()
	return nil
}

func (s *ConcurrentMap) getSession (uid int64) (*Session,error) {
	index := uid % int64(len(s.shards))
	itemMap := s.shards[index]
	if itemMap == nil {
		return nil, errs.ErrNil
	}
	itemMap.mu.RLock()
	sess := itemMap.item[uid]
	itemMap.mu.RUnlock()
	return sess,nil
}

func (s *ConcurrentMap) removeSession (uid int64) error {
	index := uid % int64(len(s.shards))
	itemMap := s.shards[index]
	if itemMap == nil {
		return errs.ErrNil
	}
	itemMap.mu.Lock()
	delete(itemMap.item,uid)
	itemMap.mu.Unlock()
	return nil
}

type ItemMap struct {
	item map[int64]*Session
	mu sync.RWMutex
}

