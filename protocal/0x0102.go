package protocal

import (
	"github.com/dushu1105/jt808/common"
	"fmt"
)

const TVerifyRequest2013 = 0x0102
const TVerifyRequest2019 = 0x10102

const (
	CarRegisted = 1
	NotFindCarInDB = 2
	TerminalRegisted = 3
	NotFindTerminalInDB = 4
)

type TVerifyHandler2013 struct {
	BaseHandler
	Token 	string
}

func (t *TVerifyHandler2013) Parse(data []byte)  error{
	err := common.ReadStruct(data, common.BigEndian, t)
	return err
}

func (t *TVerifyHandler2013) JT808Msg() byte{
	//todo
	fmt.Println(string(t.Token))
	return 0
}

func (t *TVerifyHandler2013) Do(msg *JT808Msg) (*JT808Msg, error) {
	if msg.Header.Ver == 1 {

	}
	err := t.Parse(msg.Body)
	if err != nil{
		return nil, err
	}

	r := t.JT808Msg()

	v := CommonRespHandler{Seq:msg.Header.Seq, Result:r}
	ret, err := v.Packet()
	if err != nil{
		return nil, err
	}

	return msg.CopyAndSet(PCommonResponse, ret), err
}

type TVerifyHandler2019 struct {
	BaseHandler
	Len 	byte
	Token 	string `ref:"Len"`
	IMEI    string `len:"15"`
	Ver 	string `len:"20"`
}

func (t *TVerifyHandler2019) Parse(data []byte)  error{
	err := common.ReadStruct(data, common.BigEndian, &t.Len)
	return err
}

func (t *TVerifyHandler2019) JT808Msg() byte{
	//todo
	fmt.Println(string(t.Token), t)
	return 0
}


func (t *TVerifyHandler2019) Do(msg *JT808Msg) (*JT808Msg, error) {
	if msg.Header.Ver == 1 {

	}
	err := t.Parse(msg.Body)
	if err != nil{
		return nil, err
	}

	r := t.JT808Msg()

	v := CommonRespHandler{Seq:msg.Header.Seq, Result:r}
	ret, err := v.Packet()
	if err != nil{
		return nil, err
	}

	return msg.CopyAndSet(PCommonResponse, ret), err
}