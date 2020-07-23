package protocal

import (
	"bytes"
	"encoding/binary"
)

const (
	TBatchUploadRequest = 0x0107
)

type PosInfo struct {
	Len uint16
	Pos TPositionHandler
}

type TBatchUploadHandler struct{
	BaseHandler
	Num 	uint16
	Type 	byte
	PosList []*PosInfo
}

func (t *TBatchUploadHandler) Parse(data []byte) error {
	buf := bytes.NewBuffer(data)
	err := binary.Read(buf, binary.BigEndian, &t.Num)
	if err != nil{
		return err
	}
	err = binary.Read(buf, binary.BigEndian, &t.Type)
	if err != nil{
		return err
	}

	for i:=0;i<int(t.Num);i+=1{
		var p PosInfo
		err = binary.Read(buf, binary.BigEndian, &p.Len)
		if err != nil{
			return err
		}
		d := make([]byte, p.Len)
		err = binary.Read(buf, binary.BigEndian, &d)
		if err != nil{
			return err
		}

		err = p.Pos.Parse(d)
		if err != nil{
			return err
		}
	}

	return nil
}