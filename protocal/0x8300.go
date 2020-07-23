package protocal

import (
	"808/common"
	"bytes"
)

const (
	PTextSendRequest2013 = 0x8300
	PTextSendRequest2019 = 0x18300
)

type TextFlag struct {
	Urgent bool
	Serve  bool
	Notify bool
	TerminalMonitor bool
	TTS bool
	Ad bool
	ErrCode bool
}

type PTextSendReqHandler2013 struct {
	Flag byte
	Txt  string
}

type PTextSendReqHandler2019 struct {
	Flag byte
	Type byte  //1 通知，2服务
	Txt  string
}

func (p *PTextSendReqHandler2013) Packet() ([]byte, error){
	buf := new(bytes.Buffer)
	err := common.WriteStruct(buf, common.BigEndian, p)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (p *PTextSendReqHandler2013) MakeFlag(t *TextFlag) {
	var f byte
	if t.Urgent{
		f |= 1
	}
	if t.TerminalMonitor{
		f |= 0b00000100
	}
	if t.TTS{
		f |= 0b00001000
	}
	if t.Ad {
		f |= 0b00010000
	}
	if t.ErrCode {
		f |= 0b00100000
	}
	p.Flag = f
}

func (p *PTextSendReqHandler2019) Packet() ([]byte, error){
	buf := new(bytes.Buffer)
	err := common.WriteStruct(buf, common.BigEndian, p)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (p *PTextSendReqHandler2019) MakeFlag(t *TextFlag) {
	var f byte
	if t.Urgent{
		f |= 0b00000010
	}
	if t.Serve{
		f |= 0b00000001
	}
	if t.Notify{
		f |= 0b00000011
	}
	if t.TerminalMonitor{
		f |= 0b00000100
	}
	if t.TTS{
		f |= 0b00001000
	}
	if t.ErrCode {
		f |= 0b00100000
	}
	p.Flag = f
	if t.Serve{
		p.Type = 2
	} else if t.Notify{
		p.Type = 1
	}
}