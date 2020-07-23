package protocal

import (
	"808/common"
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"io"
)

type Handler interface {
	Do(msg *JT808Msg) (*JT808Msg, error)
	Send([]byte) error
}

type BaseHandler struct{}

func (bh *BaseHandler) Do(msg *JT808Msg) (*JT808Msg, error) {
	return nil, nil
}

func (bh *BaseHandler) Send([] byte) error {
	return nil
}

const JT808Sign = 0x7e

type JT808Msg struct {
	ConnId uint32
	Header *JT808Header
	IsCompleted bool
	Body   []byte
}

func (j *JT808Msg) Printf(format string, a ...interface{}){
	fmt.Printf(fmt.Sprintf("ConnID = %d ", j.ConnId)+format, a...)
}

func (j *JT808Msg) Print(a ...interface{}){
	fmt.Println("ConnID =", j.ConnId, a)
}

func checkCode(data []byte) byte {
	var code byte
	for _, d := range data{
		code ^= d
	}
	return code
}

func (j *JT808Msg) Packet() ([]byte, error){
	buf := new(bytes.Buffer)
	hData, err := j.Header.Packet()
	if err != nil {
		return nil, err
	}

	err = common.BWrite(buf, common.BigEndian, hData)
	if err != nil {
		return nil, err
	}

	err = common.BWrite(buf, common.BigEndian, j.Body)
	if err != nil {
		return nil, err
	}

	code := checkCode(buf.Bytes())
	err = buf.WriteByte(code)
	if err != nil {
		return nil, err
	}

	return common.Encode(buf.Bytes())
}

func (j *JT808Msg) Parse(data []byte) error{
	j.Printf("recv msg<< % x\n", data)
	d, err := common.Decode(data)
	if err != nil {
		return err
	}

	code := checkCode(d)
	if code != 0 {
		j.Print("错误校验码")
		return errors.Errorf("错误的校验码")
	}

	d = d[:len(d) - 1]

	var h JT808Header
	err = h.Parse(d)
	if err != nil {
		if err == io.ErrUnexpectedEOF{
			return errors.Errorf("无效数据，JT808Header长度不够 ")
		}
		return err
	}
	j.Printf("msg header: id=%x, seq=%x, ver=%d, attr=%b, sim=%s\n", h.Id, h.Seq, h.Ver, h.Attr, string(h.Sim[:]))
	a := h.Attr
	j.Printf("msg attr: frag=%d, encrypt=%d, len=%d\n", a.FragFlag, a.EncryptType, a.BodyLen)

	var body []byte

	var dataOffset uint16

	if a.FragFlag == 1{
		if h.Frag.Seq == h.Frag.TotalNum{
			j.IsCompleted = true
		} else {
			j.IsCompleted = false
		}

		if a.Ver == 1 {
			dataOffset = JT808HeaderAdditionLen
		} else {
			dataOffset = JT808HeaderAdditionLen2011
		}

		if dataOffset+a.BodyLen > uint16(len(d)){
			return errors.Errorf("数据长度与数据头内记录不符")
		}
		body = d[dataOffset:dataOffset+a.BodyLen]
	} else {
		if a.Ver == 1 {
			dataOffset = JT808HeaderCommonLen
		} else {
			dataOffset = JT808HeaderCommonLen2011
		}
		if dataOffset+a.BodyLen > uint16(len(d)){
			return errors.Errorf("数据长度与数据头内记录不符")
		}
		body = d[dataOffset:dataOffset+a.BodyLen]
		j.IsCompleted = true
	}
	j.Body = append(j.Body, body...)
	j.Header = &h

	return nil
}

func (j *JT808Msg) CopyAndSet(Id uint16, body []byte) *JT808Msg{
	var r JT808Msg
	r.Body = body
	r.Header = j.Header
	r.Header.Id = Id
	v := JT808HeaderAttr{BodyLen: uint16(len(body))}
	r.Header.Attr = v
	return &r
}

func (j *JT808Msg) IsVerifyMsg() bool {
	return j.Header.Id == TVerifyRequest2013
}

func (j *JT808Msg) Is2019Ver() bool {
	return j.Header.Ver == 1
}

type JT808HeaderAttr struct {
	Ver 		uint8
	FragFlag    uint8
	EncryptType uint8
	BodyLen     uint16
}

func (j *JT808HeaderAttr) Parse(attr uint16) {
	j.Ver = uint8((attr>>14) & 0b1)
	j.FragFlag = uint8((attr >> 13) & 0b1)
	j.EncryptType = uint8((attr >> 10) & 0b111)
	j.BodyLen = attr & 0b1111111111
}

func (j *JT808HeaderAttr) Packet() uint16 {
	var attr uint16
	attr |= uint16(j.Ver) << 14
	attr |= uint16(j.FragFlag) << 13
	attr |= uint16(j.EncryptType) << 10
	attr |= j.BodyLen
	return attr
}

type JT808Header struct {
	Id		uint16 `json:"id"`
	Attr	JT808HeaderAttr
	Ver 	byte
	Sim 	string
	Seq 	uint16 `json:"seq"`
	Frag 	JT808HeaderFrag `json:"frag"`
}

type JT808HeaderCommon struct {
	Id   uint16 `json:"id"`
	Attr uint16
}

type JT808Header2011 struct {
	Sim  [JT808HeaderSimLen2011]byte
	Seq 	uint16 `json:"seq"`
}
type JT808Header2019 struct {
	Ver 	byte
	Sim  [JT808HeaderSimLen2019]byte
	Seq 	uint16 `json:"seq"`
}

type JT808HeaderFrag struct {
	TotalNum  uint16
	Seq 	  uint16
}

const JT808HeaderCommonLen2011 = 12
const JT808HeaderAdditionLen2011 = 16
const JT808HeaderSimLen2011 = 6
const JT808HeaderSimLen2019 = 10

const JT808HeaderCommonLen = 17
const JT808HeaderAdditionLen = 21



func (j *JT808Header) Parse(data []byte) error {
	buf := bytes.NewBuffer(data)
	var c JT808HeaderCommon
	err := common.BRead(buf, common.BigEndian, &c)
	if err != nil {
		return err
	}

	var a JT808HeaderAttr
	a.Parse(c.Attr)
	j.Id = c.Id
	j.Attr = a
	j.Ver = a.Ver
	if a.Ver == 1 {
		var j1 JT808Header2019
		err = common.BRead(buf, common.BigEndian, &j1)
		if err != nil {
			return err
		}

		j.Sim = common.BCD2DEC(j1.Sim[:])
		j.Seq = j1.Seq
	} else {
		j.Ver = 0
		var j2 JT808Header2011
		err = common.BRead(buf, common.BigEndian, &j2)
		if err != nil {
			return err
		}
		j.Sim = common.BCD2DEC(j2.Sim[:])
		j.Seq = j2.Seq
	}

	if a.FragFlag == 1{
		var f JT808HeaderFrag
		err = common.BRead(buf, common.BigEndian, &f)
		if err != nil {
			return err
		}
		j.Frag = f
	}

	return nil
}

func (j *JT808Header) Packet() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := common.BWrite(buf, common.BigEndian, j.Id)
	if err != nil {
		return nil, err
	}

	err = common.BWrite(buf, common.BigEndian, j.Attr.Packet())
	if err != nil {
		return nil, err
	}

	s := common.Str2BCD(j.Sim)
	l := len(s)
	if j.Attr.Ver == 1{
		err = common.BWrite(buf, common.BigEndian, j.Ver)
		if err != nil {
			return nil, err
		}
		if l < JT808HeaderSimLen2019{
			sim := make([]byte, JT808HeaderSimLen2019 - l)
			s = append(sim, s...)
		} else {
			s = s[:JT808HeaderSimLen2019]
		}
	} else {
		if l < JT808HeaderSimLen2011{
			sim := make([]byte, JT808HeaderSimLen2011 - l)
			s = append(sim, s...)
		} else {
			s = s[:JT808HeaderSimLen2011]
		}
	}

	err = common.BWrite(buf, common.BigEndian, s)
	if err != nil {
		return nil, err
	}

	err = common.BWrite(buf, common.BigEndian, j.Seq)
	if err != nil {
		return nil, err
	}

	if j.Attr.FragFlag == 1{
		err = common.BWrite(buf, common.BigEndian, j.Frag)
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func MakeJt808Msg(id, seq uint16, ver byte, fragflag uint8, sim string, data []byte) ([]byte, error) {
	var h JT808Header
	h.Id = id
	h.Ver = ver
	h.Seq = seq
	h.Sim = sim

	var a JT808HeaderAttr
	a.Ver = ver
	a.FragFlag = fragflag
	a.BodyLen = uint16(len(data))
	h.Attr = a

	var m JT808Msg
	m.Body = data
	m.Header = &h

	return m.Packet()
}