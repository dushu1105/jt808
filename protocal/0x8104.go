package protocal

import (
	"github.com/dushu1105/jt808/common"
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	PQueryParamRequest = 0x8104
	PQuerySpecifyParamRequest = 0x8106
	TQueryParamResponse = 0x0104
)

type PQuerySpecifyParamHandler struct{
	Num 	byte
	List 	[]uint32
}

func (p *PQuerySpecifyParamHandler) Packet() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, p.Num)
	if err != nil{
		return nil, err
	}

	for _, v := range p.List{
		err = binary.Write(buf, binary.BigEndian, v)
		if err != nil{
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

type PQueryParamHandler struct{}

func (p *PQueryParamHandler) Packet() ([]byte, error) {
	return nil, nil
}

var paramMap = map[uint32]string{
	0x0001:"uint32", //终端心跳间隔 单位s
	0x0002:"uint32",//TCP应答超时时间，单位s
	0x0003:"uint32",//tcp重传次数
	0x0004:"uint32",//udp应答超时时间， 单位s
	0x0005:"uint32",//udp重传次数
	0x0006:"uint32",//sms应答超时时间，单位s
	0x0007:"uint32",//sms重传次数
	0x0010:"string",//主服务器APN，无线通信拨号访问点，若网络制式cdma，则此处是ppp拨号号码
	0x0011:"string",//主服务器无线通信拨号用户名
	0x0012:"string",//主服务器无线通信拨号密码
	0x0013:"string",//主服务器备份地址，ip或域名，冒号分割主机和端口，多个用分号分隔
	0x0014:"string",//备用服务器APN
	0x0015:"string",//备用服务器无线通信拨号用户名
	0x0016:"string",//备用服务器无线通信拨号密码
	0x0017:"string",//备用服务器备份地址，ip或域名，冒号分割主机和端口，多个用分号分隔
	0x001A:"string",//道路运输证IC卡认证主服务器ip或域名
	0x001B:"uint32",//道路运输证IC卡认证主服务器tcp端口
	0x001C:"uint32",//道路运输证IC卡认证主服务器udp端口
	0x001D:"string",//道路运输证IC卡认证备份服务器ip或域名，端口同主的
	0x0020:"uint32", //位置汇报策略，0：定时汇报；1：定距汇报；2：定时和定距汇报
	0x0021:"uint32", //位置汇报方案，0：根据 ACC 状态； 1：根据登录状态和 ACC 状态，先判断登录状态，若登录再根据 ACC 状态
	0x0022:"uint32", //驾驶员未登录汇报时间间隔，单位为秒（s），>0
	0x0023:"string", //从服务器APN。该值为空时，终端应使用主服务器相同配置
	0x0024:"string", //从服务器无线通信拨号用户名。该值为空时，终端应使用主服务器相同配置
	0x0025:"string", //从服务器无线通信拨号密码。该值为空，终端应使用主服务器相同配置
	0x0026:"string", //从服务器备份地址IP。该值为空，终端应使用主服务器相同配置
	0x0027:"uint32", //休眠时汇报时间间隔，单位为秒（s），>0
	0x0028:"uint32", //紧急报警时汇报时间间隔，单位为秒（s），>0
	0x0029:"uint32", //缺省时间汇报间隔，单位为秒（s），>0
}
type QueryParam struct {
	Id 	uint32
	Len byte
	Value []byte `ref:"Len"`
}

type TQueryParamHandler struct{
	BaseHandler
	Seq 	uint16
	Num 	byte
	List 	[]*QueryParam `ref:"Num"`
}

func (t *TQueryParamHandler) Parse(data []byte) error {
	err := common.ReadStruct(data, common.BigEndian, t)
	if err != nil{
		return err
	}
	return nil
}

func (q *QueryParam) Show(){
	if JT808Msg, ok := paramMap[q.Id]; !ok{
		fmt.Printf("暂不支持参数:%x\n", q.Id)
	} else {
		if JT808Msg == "uint32"{
			b := bytes.NewBuffer(q.Value)
			var v uint32
			err := binary.Read(b, binary.BigEndian, &v)
			if err != nil{
				fmt.Printf("参数%x值错误:%s", q.Id, err.Error())
			} else {
				fmt.Printf("参数%x：%d", q.Id, v)
			}
		} else if JT808Msg == "string"{
			fmt.Printf("参数%x：%s", q.Id, string(q.Value))
		}
	}
}

func (t *TQueryParamHandler) Show(){
	fmt.Printf("%d 查询结果共有 %d 个参数\n", t.Seq, t.Num)
	for i:=0;i<int(t.Num);i+=1{
		t.List[i].Show()
	}
}

func (t *TQueryParamHandler) Do(msg *JT808Msg) (*JT808Msg, error) {
	err := t.Parse(msg.Body)
	if err != nil{
		return nil, err
	}

	//todo param
	t.Show()

	return nil, nil
}
