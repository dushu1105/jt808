package protocal

import (
	"github.com/dushu1105/jt808/common"
	"fmt"
)

const (
	PDriverInfoRequest = 0x8702
	TDriverInfoResponse = 0x0702
)

type PDriverInfoHandler struct{}

func (p *PDriverInfoHandler) Packet() ([]byte, error) {
	return nil, nil
}

type TDriverInfoHandler struct{
	BaseHandler
	Status 	byte
	Time 	string `len:"6" type:"time"`
	Result  byte
	NameLen    byte
	Name 	string
	Cert  string `len:"20"`
	InstituteLen byte
	Institue string
	ValidDate string `len:"6" type:"time"`
}

func (t *TDriverInfoHandler) Parse(data []byte) error {
	err := common.ReadStruct(data, common.BigEndian, t)
	return err
}

var driverStatusMap = map[uint8]string{
	1:"从业资格证ic卡插入",
	2:"从业资格证ic卡拔出",
}

var icResultMap = map[uint8]string{
	0:"IC卡读卡成功",
	1:"读卡失败，卡片密匙认证未成功",
	2:"读卡失败，卡片已被锁定",
	3:"读卡失败，卡片已被拔出",
	6:"读卡失败，数据校验错误",
}

func (t *TDriverInfoHandler) Show(){
	fmt.Println(t)
	for i:=0;i<8;i+=1{
		if v, ok := driverStatusMap[(t.Status >> i) & 0x01];ok{
			fmt.Println(v)
		}
	}

	for i:=0;i<8;i+=1{
		if v, ok := icResultMap[(t.Result >> i) & 0x01];ok{
			fmt.Println(v)
		}
	}

	if t.Result == 0 {
		fmt.Println(t)
	}
	//save to somewhere
}