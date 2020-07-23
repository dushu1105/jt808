package protocal

import (
	"808/common"
	"bytes"
)

const PStreamRequest = 0x9101

type PStreamReqHandler struct {
	IPLen byte
	Ip string `ref:"IPLen"`
	TcpPort uint16
	UdpPort uint16
	Channel byte
	DataType byte //0音频， 1视频，2双向对讲，3监听， 4中心广播， 5透出
	StreamType byte //0主码流，1子码流
}

func (p *PStreamReqHandler) Packet() ([]byte, error){
	buf := new(bytes.Buffer)
	err := common.WriteStruct(buf, common.BigEndian, p)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}