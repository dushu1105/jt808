package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/dushu1105/jt808/protocal"
	"net"
	"os"
	"time"
)

func BufioRead(name string, conn net.Conn) {

}

func downStream(name string, conn net.Conn) {
	var err error
	p := protocal.TRegistReqHandler2013{
		Provice:     1,
		City:        2,
		Color:       3,
	}
	p.Manu = "chinaa"
	p.TType = "abc1234567890"
	p.TId = "987654321"
	p.Licence = "京p12345"

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, p.Provice)
	binary.Write(buf, binary.BigEndian, p.City)
	binary.Write(buf, binary.BigEndian, p.Manu)
	binary.Write(buf, binary.BigEndian, p.TType)
	binary.Write(buf, binary.BigEndian, p.TId)
	binary.Write(buf, binary.BigEndian, p.Color)
	binary.Write(buf, binary.BigEndian, p.Licence)
	data := buf.Bytes()

	a := protocal.JT808HeaderAttr{
		FragFlag:    0,
		EncryptType: 0,
		BodyLen:     uint16(len(data)),
	}

	h := protocal.JT808Header{
		Id:   0x0004,
		Attr: a,
		Seq:  1,
	}

	h.Sim = "13911111111"

	j := protocal.JT808Msg{
		Header:      &h,
		IsCompleted: false,
		Body:        data,
	}
	sendData, err := j.Packet()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("send data % x\n", sendData)
	n1, err := conn.Write(sendData)
	if err != nil {
		fmt.Println("send data error,", err)
		return
	} else {
		fmt.Println("send data len,", n1)
	}

	readData := make([]byte, 1024)
	n, err := conn.Read(readData)
	if err != nil{
		fmt.Println(err)
		return
	}
	fmt.Printf("readData % x\n", readData[:n])
}

func str2BCD(s string) []uint8 {
	n := 12
	numZero := n - len(s)
	if numZero > 0 {
		for i := 0; i < numZero; i++ {
			s = "0" + s
		}
	}
	fmt.Println("str2bcd, s=", s)
	ret := make([]uint8, n/2)
	for i := 0; i < n; i += 2 {
		a := int(s[i]) - int('0')
		b := int(s[i+1]) - int('0')
		ret[i/2] = uint8(a<<4 | b)
	}
	return ret
}
func BCD2DEC(bcd uint8) uint8 {
	return (bcd - (bcd>>4)*6)
}

/*
	模拟客户端
*/
func main() {
	fmt.Println("Client Test ... start")
	//3秒之后发起测试请求，给服务端开启服务的机会
	//time.Sleep(3 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:808")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}
	time.Sleep(3000 * time.Second)
	return
	fname := os.Args[2]
	down := os.Args[1]
	fmt.Println("arg,", fname)
	//time.Sleep(time.Second * 3)
	//BufioRead("./tcpdump.bin", conn)
	if down == "1"{
		BufioRead(fname, conn)
	} else {
		downStream(fname, conn)
	}
}
