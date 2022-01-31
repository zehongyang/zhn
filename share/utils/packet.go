package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"google.golang.org/protobuf/proto"
	"gopros/proto/pb"
	"gopros/share/errs"
	"gopros/sockets"
)

func Pack(message proto.Message,dt pb.DataType) ([]byte,error) {
	var (
		bts []byte
		err error
	)
	switch dt {
	case pb.DataType_DT_Proto:
		bts, err = proto.Marshal(message)
	case pb.DataType_DT_Json:
		bts, err = json.Marshal(message)
	default:
		return nil,errs.ErrDataType
	}
	if err != nil {
		return nil,err
	}
	if len(bts) > sockets.DataMaxLength || len(bts) < 1 {
		return nil,errs.ErrDataLength
	}
	var buf bytes.Buffer
	ulen := uint32(dt) << 24 | uint32(len(bts))
	var slen [sockets.DataLength]byte
	binary.BigEndian.PutUint32(slen[:],ulen)
	err = binary.Write(&buf, binary.BigEndian, slen[:])
	if err != nil {
		return nil,err
	}
	err = binary.Write(&buf, binary.BigEndian, bts)
	if err != nil {
		return nil,err
	}
	return buf.Bytes(),nil
}
