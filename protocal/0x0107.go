package protocal

import (
	"github.com/dushu1105/jt808/common"
	"fmt"
)

const (
	TQueryAttrResponse2019 = 0x10107
	TQueryAttrResponse2013 = 0x0107
)

type TQueryAttrHandler2019 struct{
	BaseHandler
	Type 	uint16
	Manu 	string `len:"5"`
	Model   string `len:"30"`
	Id 		string `len:"30"`
	Sim 	string `len:"10" type:"bcd"`
	HWVerLen  byte
	HWVer   string `ref:"HWVerLen"`
	SWVerLen byte
	SWVer   string `ref:"SWVerLen"`
	GNSS    byte
	Com     byte
}

func (t *TQueryAttrHandler2019) Parse(data []byte) error {
	err := common.ReadStruct(data, common.BigEndian, t)
	return err
}

type TQueryAttrHandler2013 struct{
	BaseHandler
	Type 	uint16
	Manu 	string `len:"5"`
	Model   string `len:"20"`
	Id 		string `len:"7"`
	Sim 	string `len:"10" type:"bcd"`
	HWVerLen  byte
	HWVer   string `ref:"HWVerLen"`
	SWVerLen byte
	SWVer   string `ref:"SWVerLen"`
	GNSS    byte
	Com     byte
}

func (t *TQueryAttrHandler2013) Parse(data []byte) error {
	err := common.ReadStruct(data, common.BigEndian, t)
	return err
}

var terminalTypeMap = map[uint16]string{
	0:"客运车辆",
	1:"危险品车辆",
	2:"普通货运车辆",
	3:"出租车辆",
	6:"硬盘录像",
	7:"分体机",
	8:"挂车",
}

var terminalGNSSMap = map[uint8]string{
	0:"GPS",
	1:"北斗",
	2:"GLONASS",
	3:"GALILEO",
}

var terminalComMap = map[uint8]string{
	0:"GPRS",
	1:"CDMA",
	2:"TD-SCDMA",
	3:"WCDMA",
	4:"CDMA2000",
	5:"TD-LTE",
	6:"其它通信",
}

func (t *TQueryAttrHandler2013) Show(){
	fmt.Println(t)
	for i:=0;i<8;i+=1{
		if v, ok := terminalTypeMap[(t.Type >> i) & 0x01];ok{
			fmt.Println("支持", v)
		}
	}

	for i:=0;i<8;i+=1{
		if v, ok := terminalGNSSMap[(t.GNSS >> i) & 0x01];ok{
			fmt.Println("支持", v)
		}
	}

	for i:=0;i<8;i+=1{
		if v, ok := terminalComMap[(t.Com >> i) & 0x01];ok{
			fmt.Println("支持", v)
		}
	}
	fmt.Println(t)
	//may save in db
}


func (t *TQueryAttrHandler2019) Show(){
	fmt.Println(t)
	for i:=0;i<8;i+=1{
		if v, ok := terminalTypeMap[(t.Type >> i) & 0x01];ok{
			fmt.Println("支持", v)
		}
	}

	for i:=0;i<8;i+=1{
		if v, ok := terminalGNSSMap[(t.GNSS >> i) & 0x01];ok{
			fmt.Println("支持", v)
		}
	}

	for i:=0;i<8;i+=1{
		if v, ok := terminalComMap[(t.Com >> i) & 0x01];ok{
			fmt.Println("支持", v)
		}
	}
	fmt.Println(t)
	//may save in db
}