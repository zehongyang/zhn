package errs

import "errors"

var (
	ErrConfig = errors.New("load config error")
	ErrSession = errors.New("session error")
	ErrNil  = errors.New("nil error")
	ErrHandle = errors.New("no handler to do")
	ErrAddr = errors.New("addr error")
	ErrBody= errors.New("request body err")
	ErrDataType = errors.New("unkonwn data type")
	ErrReq = errors.New("request err")
	ErrInvalid = errors.New("invalid err")
	ErrDataLength = errors.New("data length err")
)
