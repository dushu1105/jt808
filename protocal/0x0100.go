package protocal

import (
	"github.com/dushu1105/jt808/common"
	"bytes"
	"fmt"
)

const TRegistRequest2013 = 0x0100
const TRegistRequest2019 = 0x10100
const PRegistResponse = 0x8100

type TRegistReqHandler2013 struct {
	BaseHandler
	Provice uint16
	City    uint16
	Manu	string `len:"5"`
	TType   string `len:"20"`
	TId 	string `len:"7"`
	Color   byte
	Licence string
}

type TRegistReqHandler2019 struct {
	BaseHandler
	Provice uint16
	City    uint16
	Manu	string `len:"11"`
	TType   string `len:"30"`
	TId 	string `len:"30"`
	Color   byte
	Licence string
}

type PRegistResp struct {
	Seq 	uint16
	Result 	byte
	Token 	string
}

func (t *TRegistReqHandler2019) Parse(data []byte)  error{
	err := common.ReadStruct(data, common.BigEndian, t)
	return err
}

func (t *TRegistReqHandler2013) Parse(data []byte)  error{
	err := common.ReadStruct(data, common.BigEndian, t)
	return err
}

func (t *PRegistResp) Packet() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := common.WriteStruct(buf, common.BigEndian, t)
	if err != nil{
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t *TRegistReqHandler2019) Do(msg *JT808Msg) (*Jt808ResultMsg, error) {
	err := t.Parse(msg.Body)
	if err != nil{
		return nil, err
	}
	r, token, err := t.SaveInfo()
	if err != nil{
		return nil, err
	}

	token = "abcd"
	v := PRegistResp{Seq:msg.Header.Seq, Result:r, Token:token}
	ret, err := v.Packet()
	if err != nil{
		return nil, err
	}

	return &Jt808ResultMsg{Msg:msg.CopyAndSet(PRegistResponse, ret), NeedFeedBack:true}, err
}

func (t *TRegistReqHandler2013) Do(msg *JT808Msg) (*Jt808ResultMsg, error) {
	err := t.Parse(msg.Body)
	if err != nil{
		return nil, err
	}
	r, token, err := t.SaveInfo()
	if err != nil{
		return nil, err
	}

	token = "abcd"
	v := PRegistResp{Seq:msg.Header.Seq, Result:r, Token:token}
	ret, err := v.Packet()
	if err != nil{
		return nil, err
	}

	return &Jt808ResultMsg{Msg:msg.CopyAndSet(PRegistResponse, ret), NeedFeedBack:true}, err
}


func (t *TRegistReqHandler2019) SaveInfo() (byte, string, error){
	//todo
	fmt.Printf("provice=%d, city=%d, licence=%s, manu=%s, type=%s, id=%s\n", t.Provice, t.City, string(t.Licence), string(t.Manu[:]), string(t.TType[:]), string(t.TId[:]))
	return 0, "nil", nil
}

func (t *TRegistReqHandler2013) SaveInfo() (byte, string, error){
	//todo
	fmt.Printf("provice=%d, city=%d, licence=%s, manu=%s, type=%s, id=%s\n", t.Provice, t.City, string(t.Licence), string(t.Manu[:]), string(t.TType[:]), string(t.TId[:]))
	return 0, "nil", nil
}
