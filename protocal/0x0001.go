package protocal

import (
	"github.com/dushu1105/jt808/common"
)

const TCommonResponse = 0x0001
const PCommonResponse = 0x8001

const (
	Succ = 0
	Failed = 1
	BadMsg = 2
	NotSupport = 3
	Warning = 4
)

type CommonRespHandler struct {
	BaseHandler
	Seq 	uint16
	Id 		uint16
	Result 	byte
}

func (c *CommonRespHandler) Parse(data []byte) error {
	//用于终端通用应答
	return common.Parse(data, c)
}

func (c *CommonRespHandler) Packet() ([]byte, error) {
	//用于平台通用应答
	c.Id = PCommonResponse
	return common.Packet(c)
}


func (c *CommonRespHandler) Do(msg *JT808Msg) (*JT808Msg, error) {
	err := c.Parse(msg.Body)
	if err != nil{
		return nil, err
	}

	//todo common response
	switch c.Result {
	case Failed:
		msg.Printf("Seq %d Command 0x%x failed\n", c.Seq, c.Id)
		break
	case BadMsg:
		msg.Printf("Seq %d Command 0x%x bad message\n", c.Seq, c.Id)
		break
	case NotSupport:
		msg.Printf("Seq %d Command 0x%x not support\n", c.Seq, c.Id)
		break
	case Warning:
		msg.Printf("Seq %d Command 0x%x warning\n", c.Seq, c.Id)
		break
	default:
	}

	return nil, err
}