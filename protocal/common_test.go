package protocal

import (
	"808/common"
	"bytes"
	"fmt"
	"testing"
)
var testString = []byte{0x30, 0x7e, 0x08, 0x7d, 0x55}

//func TestEncodeDecode(t *testing.T) {
//	t.Logf("% x", testString)
//	r, err := Encode(testString)
//	if err != nil{
//		t.Errorf("encode error, %s", err)
//	} else {
//		t.Logf("% x", r)
//		r1, err := Decode(r)
//		if err != nil{
//			t.Errorf("decode error, %s", err)
//		} else {
//			t.Logf("% x", r1)
//		}
//	}
//}


func TestJT808HeaderAttr_Parse(t *testing.T) {
	var v uint16
	v = 0b0010010000001011
	j := JT808HeaderAttr{}
	j.Parse(v)
	t.Logf("%d, %d, %d", j.FragFlag, j.EncryptType, j.BodyLen)
	v1 := j.Packet()
	t.Logf("%b", v1)
}


func TestJT808Msg_Parse(t *testing.T) {
	var data = "7e 01 00 00 31 01 89 19 88 98 01 00 01 00 00 00 00 4b 6f 69 6b 65 30 30 31 4b 6f 69 6b 65 30 30 31 30 30 30 30 30 30 30 30 30 30 30 30 4b 6f 69 6b 65 30 30 31 02 d4 c1 41 31 32 33 34 35 a4 7e"

	var j JT808Msg
	err := j.Parse(s2b(data))
	if err != nil{
		t.Log(err)
		return
	}
	t.Log(j.Header, len(j.Header.Sim))
	//j.Header.Sim = "18919889801"
	d, err := j.Packet()
	if err != nil{
		t.Log(err)
		return
	}
	t.Logf("% x", d)
}


type Child struct{
	X byte
	_Y byte
	_Z byte
}

type JT808Test struct {
	BaseHandler
	Id		uint16 `json:"id"`
	Attr	uint16
	Ver 	byte
	OldSim  string `len:"6" type:"time"`
	_OO     byte
	Seq 	uint16 `json:"seq"`
	Len     byte
	//Ad      string `ref:"Len"`
	Ad      []byte `ref:"Len"`
	Last    uint8
}

//func TestByte2Struct(t *testing.T) {
//	var data = "01 00 00 31 01 11 11 11 11 11 00 01 01 02 03 04"
//	var j JT808Test
//	//j.Ad = make([]byte, 4)
//	buf := bytes.NewBuffer(s2b(data))
//	err := BRead(buf, BigEndian, &j)
//	if err != nil{
//		t.Error(err)
//		return
//	}
//	t.Logf("% x", j)
//
//	buf = new(bytes.Buffer)
//	err = BWrite(buf, BigEndian, j)
//	if err != nil{
//		t.Error(err)
//		return
//	}
//	t.Logf("% x", buf.Bytes())
//}

func TestReadStruct(t *testing.T) {
	var data = "01 00 00 31 01 20 11 11 11 11 12 00 01 03 01 02 03 04"
	var j JT808Test
	//j.Ad = make([]byte, 4)
	//buf := bytes.NewBuffer(s2b(data))
	err := common.ReadStruct(s2b(data), common.BigEndian, &j)
	if err != nil{
		t.Error(err)
		return
	}
	t.Logf("%x", j)
	t.Logf("%v, %v", j.Ad,j.OldSim)

	d := new(bytes.Buffer)
	err = common.WriteStruct(d, common.BigEndian, &j)
	if err != nil{
		t.Error(err)
		return
	}
	t.Logf("% x", d.Bytes())

	var x = ^uint16(0) - 1
	fmt.Println(x+1)
}