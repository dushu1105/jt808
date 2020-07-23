package protocal

import (
	"github.com/dushu1105/jt808/common"
	"bytes"
	"fmt"
	"github.com/pkg/errors"
)

const (
	TPositionRequest = 0x10200
	TPositionRequest1 = 0x0200
)

type TPositionHandler struct{
	BaseHandler
	WarningFlag 	uint32
	Status		 	uint32
	Latitude		 	uint32
	Longitude		 	uint32
	Height		 	uint16
	Speed		 	uint16
	Direction		 	uint16 // 方向 0-359，正北为 0，顺时针
	Time 			string  `len:"6" type:"time"`
}

func (t *TPositionHandler) Len() int{
	return 4 + 4 + 4 + 4 + 2 + 2 + 2 + 6
}

func (t *TPositionHandler) Parse(data []byte) error {
	err := common.ReadStruct(data, common.BigEndian, t)
	return err
}

func (t *TPositionHandler) Show(){
	var s TPosStatus
	s.Parse(t.Status)
	s.Show()
	var f TPosWarningFlag
	f.Parse(t.WarningFlag)
	f.Show()
}

type TPosStatus struct {
	ACC               uint8
	Location          uint8
	South             uint8
	West              uint8
	Offline            uint8
	PosEncypt         uint8
	WarningForBrake   uint8
	WarningForDeflect uint8
	Load              uint8
	Oil 			  uint8
	Circuit 			  uint8
	Lock 			  uint8
	FrontDoor 			  uint8
	MiddleDoor 			  uint8
	BackDoor 			  uint8
	DriverDoor 			  uint8
	Door 			  uint8
	GPS 			  uint8
	BeiDou 			  uint8
	Glonass			  uint8
	Galileo			  uint8
	Running			  uint8
}

func (t *TPosStatus) Parse(status uint32){
	t.ACC = uint8(status >> 31)
	t.Location = uint8((status >> 30) & 0x1)
	t.South = uint8((status >> 29) & 0x1)
	t.West = uint8((status >> 28) & 0x1)
	t.Offline = uint8((status >> 27) & 0x1)
	t.PosEncypt = uint8((status >> 26) & 0x1)
	t.WarningForBrake = uint8((status >> 25) & 0x1)
	t.WarningForDeflect = uint8((status >> 24) & 0x1)
	t.Load = uint8((status >> 23) & 0x3)
	t.Oil = uint8((status >> 21) & 0x1)
	t.Circuit = uint8((status >> 20) & 0x1)
	t.Lock = uint8((status >> 19) & 0x1)
	t.FrontDoor = uint8((status >> 18) & 0x1)
	t.MiddleDoor = uint8((status >> 17) & 0x1)
	t.BackDoor = uint8((status >> 16) & 0x1)
	t.DriverDoor = uint8((status >> 15) & 0x1)
	t.Door = uint8((status >> 14) & 0x1)
	t.GPS = uint8((status >> 13) & 0x1)
	t.BeiDou = uint8((status >> 12) & 0x1)
	t.Glonass = uint8((status >> 11) & 0x1)
	t.Galileo = uint8((status >> 10) & 0x1)
	t.Running = uint8((status >> 9) & 0x1)
}

func (t *TPosStatus) Show(){
	var s = make([]string, 0)
	if t.ACC == 1{
		s = append(s, "ACC 开")
	} else{
		s = append(s, "ACC 关")
	}
	if t.Location == 1{
		s = append(s, "定位")
	} else{
		s = append(s, "未定位")
	}
	if t.South == 1{
		s = append(s, "南纬")
	} else{
		s = append(s, "北纬")
	}
	if t.West == 1{
		s = append(s, "西经")
	} else{
		s = append(s, "东经")
	}
	if t.Offline == 1{
		s = append(s, "停运")
	} else{
		s = append(s, "运营")
	}
	if t.PosEncypt == 1{
		s = append(s, "经纬度保密插件加密")
	} else{
		s = append(s, "经纬度未保密插件加密")
	}
	if t.WarningForBrake == 1{
		s = append(s, "紧急刹车系统采集的前撞预警")
	} else{
		s = append(s, "非紧急刹车系统采集的前撞预警")
	}
	if t.WarningForDeflect == 1{
		s = append(s, "车道偏移预警")
	} else{
		s = append(s, "非车道偏移预警")
	}
	if t.Load == 0{
		s = append(s, "空车")
	} else if t.Load == 1{
		s = append(s, "半载")
	}else if t.Load == 3{
		s = append(s, "满载")
	}
	if t.Oil == 1{
		s = append(s, "车辆油路断开")
	} else{
		s = append(s, "车辆油路正常")
	}
	if t.Circuit == 1{
		s = append(s, "车辆电路断开")
	} else{
		s = append(s, "车辆电路正常")
	}
	if t.Lock == 1{
		s = append(s, "车门解锁")
	} else{
		s = append(s, "车门加锁")
	}
	if t.FrontDoor == 1{
		s = append(s, "前门开")
	} else{
		s = append(s, "前门关")
	}
	if t.MiddleDoor == 1{
		s = append(s, "中门开")
	} else{
		s = append(s, "中门关")
	}
	if t.BackDoor == 1{
		s = append(s, "后门开")
	} else{
		s = append(s, "后门关")
	}
	if t.DriverDoor == 1{
		s = append(s, "驾驶门开")
	} else{
		s = append(s, "驾驶门关")
	}
	if t.Door == 1{
		s = append(s, "自定义门开")
	} else{
		s = append(s, "自定义门关")
	}
	if t.GPS == 1{
		s = append(s, "使用GPS卫星定位")
	} else{
		s = append(s, "未使用GPS卫星定位")
	}
	if t.BeiDou == 1{
		s = append(s, "使用北斗卫星定位")
	} else{
		s = append(s, "未使用北斗卫星定位")
	}
	if t.Glonass == 1{
		s = append(s, "使用Glonass卫星定位")
	} else{
		s = append(s, "未使用Glonass卫星定位")
	}
	if t.Galileo == 1{
		s = append(s, "使用Galileo卫星定位")
	} else{
		s = append(s, "未使用Galileo卫星定位")
	}
	if t.Running == 1{
		s = append(s, "车辆行驶状态")
	} else{
		s = append(s, "车辆停止状态")
	}
	fmt.Println("状态：", s)
}

type TPosWarningFlag struct {
	Urgency               uint8
	Overspeed          uint8
	Fatigue             uint8
	Danger              uint8
	GNSSFail            uint8
	GNSSDisconnect         uint8
	GNSSShortcut   uint8
	Undervoltage uint8
	PowerDown uint8
	LCDFailed              uint8
	TTSFailed 			  uint8
	CameraFailed 			  uint8
	ICFailed 			  uint8
	EOverspeed 			  uint8
	EFatigue 			  uint8
	Irregularity 			  uint8
	ETire 			  uint8
	RightBlind 			  uint8
	OvertimeDriver 			  uint8
	OverTimePark 			  uint8
	Area			  uint8
	Line			  uint8
	DrivingTime			  uint8
	Deflect			  uint8
	VSS			  uint8
	Oil			  uint8
	Steal			  uint8
	IllegalIgnite			  uint8
	Illegalshift			  uint8
	Rollover			  uint8
	ERollover			  uint8
}

func (t *TPosWarningFlag) Parse(status uint32){
	i := 31
	t.Urgency = uint8(status >> i)
	i -= 1
	t.Overspeed = uint8((status >> i) & 0x1)
	i -= 1
	t.Fatigue = uint8((status >> i) & 0x1)
	i -= 1
	t.Danger = uint8((status >> i) & 0x1)
	i -= 1
	t.GNSSFail = uint8((status >> i) & 0x1)
	i -= 1
	t.GNSSDisconnect = uint8((status >> i) & 0x1)
	i -= 1
	t.GNSSShortcut = uint8((status >> i) & 0x1)
	i -= 1
	t.Undervoltage = uint8((status >> i) & 0x1)
	i -= 1
	t.PowerDown = uint8((status >> i) & 0x1)
	i -= 1
	t.LCDFailed = uint8((status >> i) & 0x3)
	i -= 1
	t.TTSFailed = uint8((status >> i) & 0x3)
	i -= 1
	t.CameraFailed = uint8((status >> i) & 0x1)
	i -= 1
	t.ICFailed = uint8((status >> i) & 0x1)
	i -= 1
	t.EOverspeed = uint8((status >> i) & 0x1)
	i -= 1
	t.EFatigue = uint8((status >> i) & 0x1)
	i -= 1
	t.Irregularity = uint8((status >> i) & 0x1)
	i -= 1
	t.ETire = uint8((status >> i) & 0x1)
	i -= 1
	t.RightBlind = uint8((status >> i) & 0x1)
	i -= 1
	t.OvertimeDriver = uint8((status >> i) & 0x1)
	i -= 1
	t.OverTimePark = uint8((status >> i) & 0x1)
	i -= 1
	t.Area = uint8((status >> i) & 0x1)
	i -= 1
	t.Line = uint8((status >> i) & 0x1)
	i -= 1
	t.DrivingTime = uint8((status >> i) & 0x1)
	i -= 1
	t.Deflect = uint8((status >> i) & 0x1)
	i -= 1
	t.VSS = uint8((status >> i) & 0x1)
	i -= 1
	t.Oil = uint8((status >> i) & 0x1)
	i -= 1
	t.Steal = uint8((status >> i) & 0x1)
	i -= 1
	t.IllegalIgnite = uint8((status >> i) & 0x1)
	i -= 1
	t.Illegalshift = uint8((status >> i) & 0x1)
	i -= 1
	t.Rollover = uint8((status >> i) & 0x1)
	i -= 1
	t.ERollover = uint8((status >> i) & 0x1)
}

func (t *TPosWarningFlag) Show(){
	var s = make([]string, 0)
	if t.Urgency == 1{
		s = append(s, "紧急报警，触动报警开关后")
	}
	if t.Overspeed == 1{
		s = append(s, "超速报警")
	}
	if t.Fatigue == 1{
		s = append(s, "疲劳驾驶报警")
	}
	if t.Danger == 1{
		s = append(s, "危险行为驾驶报警")
	}
	if t.GNSSFail == 1{
		s = append(s, "GNSS模块故障报警")
	}
	if t.GNSSDisconnect == 1{
		s = append(s, "GNSS天线未接或被剪断报警")
	}
	if t.GNSSShortcut == 1{
		s = append(s, "GNSS天线短路报警")
	}
	if t.Undervoltage == 1{
		s = append(s, "终端主电源欠压报警")
	}
	if t.PowerDown == 1{
		s = append(s, "终端主电源掉电报警")
	}
	if t.LCDFailed == 1{
		s = append(s, "终端LCD或显示器故障报警")
	}
	if t.TTSFailed == 1{
		s = append(s, "TTS模块故障报警")
	}
	if t.CameraFailed == 1{
		s = append(s, "摄像头故障报警")
	}
	if t.ICFailed == 1{
		s = append(s, "道路运输证IC卡模块故障报警")
	}
	if t.EOverspeed == 1{
		s = append(s, "超速预警")
	}
	if t.EFatigue == 1{
		s = append(s, "疲劳驾驶预警")
	}
	if t.Irregularity == 1{
		s = append(s, "违规行驶报警")
	}
	if t.ETire == 1{
		s = append(s, "胎压预警")
	}
	if t.RightBlind == 1{
		s = append(s, "右转盲区异常报警")
	}
	if t.OvertimeDriver == 1{
		s = append(s, "当天累计驾驶超时报警")
	}
	if t.OverTimePark == 1{
		s = append(s, "超时停车报警")
	}
	if t.Area == 1{
		s = append(s, "进出区域报警")
	}
	if t.Line == 1{
		s = append(s, "进出路线报警")
	}
	if t.DrivingTime == 1{
		s = append(s, "路段行驶时间不足或过长报警")
	}
	if t.Deflect == 1{
		s = append(s, "路线偏离报警")
	}
	if t.VSS == 1{
		s = append(s, "车辆VSS故障")
	}
	if t.Oil == 1{
		s = append(s, "车辆油量异常报警")
	}
	if t.Steal == 1{
		s = append(s, "车辆被盗报警")
	}
	if t.IllegalIgnite == 1{
		s = append(s, "车辆非法点火报警")
	}
	if t.Illegalshift == 1{
		s = append(s, "车辆非法位移报警")
	}
	if t.Rollover == 1{
		s = append(s, "碰撞侧翻报警")
	}
	if t.ERollover == 1{
		s = append(s, "侧翻预警")
	}
	fmt.Println("报警信息：", s)
}

func (t *TPositionHandler) Do(msg *JT808Msg) (*JT808Msg, error) {
	err := t.Parse(msg.Body)
	if err != nil{
		return nil, err
	}
	t.Show()

	l := uint16(t.Len())
	for ;l < msg.Header.Attr.BodyLen;{
		var p PosAddtion
		offset, err := p.Parse(msg.Body[l:])
		if err != nil{
			return nil, err
		}
		l += offset
	}


	return nil, err
}

type PosAddtion struct{
	Id  byte
	Length byte
}

func readAddition(buf *bytes.Buffer, v interface{}, k string) error {
	err := common.BRead(buf, common.BigEndian, v)
	if err != nil {
		return err
	}
	fmt.Println(k, v)
	return nil
}

type PosAreaAddition struct{
	Type byte
	Id   uint32
	Direction byte
}
type PosDriveTimeAddition struct{
	Id uint32
	Time   uint16
	Result byte
}
type PosExtStatusAddition struct{
	V uint32
}

func (p *PosExtStatusAddition) Parse(v uint32) {
	p.V = v
}

func (p *PosExtStatusAddition) Show() {
	var extMap = map[int]string{
		0:"近光灯信号",
		1:"远光灯信号",
		2:"右转向灯信号",
		3:"左转向灯信号",
		4:"制动信号",
		5:"倒档信号",
		6:"雾灯信号",
		7:"示廓灯信号",
		8:"喇叭信号",
		9:"空调状态",
		10:"空挡信号",
		11:"缓速器工作",
		12:"ABS工作",
		13:"加热器工作",
		14:"离合器状态",
	}
	var s = make([]string, 0)
	for i:=0;i < 15;i+=1{
		if (p.V >> i) & 0x01 == 1{
			s = append(s, extMap[i])
		}
	}
	fmt.Println("扩展车辆信号状态:", s)

}

type PosIOStatusAddition struct{
	DeepSleep uint8
	Sleep uint8
}

func (p *PosIOStatusAddition) Parse(v uint16) {
	p.DeepSleep = uint8(v & 0x0001)
	p.Sleep = uint8(v & 0x0002)
}
func (p *PosIOStatusAddition) Show() {
	if p.DeepSleep == 1 {
		fmt.Println("IO深度睡眠")
	}else if p.Sleep == 1{
		fmt.Println("IO睡眠")
	}

}

type PosAnalogAddition struct{
	AD0 	uint16
	AD1     uint16
}

func (p *PosAnalogAddition) Parse(v uint32) {
	p.AD0 = uint16(v & 0x0011)
	p.AD1 = uint16((v >> 16) & 0x0011)
}

func (p *PosAnalogAddition) Show() {
	fmt.Println("模拟量", p.AD0, p.AD1)
}

var posTypeMap = map[byte]string{0:"无特定位置",1:"圆形区域",2:"矩形区域",3:"多边形区域",4:"路段"}

func (p *PosAddtion) Parse(data []byte) (uint16, error) {
	buf := bytes.NewBuffer(data)


	err := common.BRead(buf, common.BigEndian, p)
	if err != nil {
		return 0, err
	}
	switch p.Id {
	case 0x01:
		//里程，单位1/10km，车上的里程表读数
		var v uint32
		err = readAddition(buf, &v, "里程:")
		break
	case 0x02:
		//油量，单位1/10L，车上的油量读数
		var v uint16
		err = readAddition(buf, &v, "油量:")
		break
	case 0x03:
		//速度，单位1/10km/hour
		var v uint16
		err = readAddition(buf, &v, "速度:")
		break
	case 0x04:
		//需要人工确认的报警id
		var v uint16
		err = readAddition(buf, &v, "人工确认ID:")
		break
	case 0x05:
		//胎压，单位pa，轮子从车头开始从左到右，前左1，前左2，前右1，前右2， 中左1 。。。
		var v [30]byte
		err = readAddition(buf, &v, "胎压:")
		break
	case 0x06:
		//温度，摄氏度，最高位1表示负
		var v int16
		err = readAddition(buf, &v, "温度:")
		break
	case 0x11:
		//超速信息
		var posType byte
		err = common.BRead(buf, common.BigEndian, &posType)
		if err != nil {
			return 0, err
		}
		if posType > 4{
			return 0, errors.Errorf("超速信息，位置值错误:%d", posType)
		}

		if posType == 0{
			fmt.Println("超速信息：", posTypeMap[posType])
		} else {
			var id uint32
			err = common.BRead(buf, common.BigEndian, &id)
			if err != nil {
				return 0, err
			}
			fmt.Println("超速信息：", posTypeMap[posType], id)
		}
		break
	case 0x12:
		//进出区域或路线报警信息
		var v PosAreaAddition
		err = common.BRead(buf, common.BigEndian, &v)
		if err != nil {
			return 0, err
		}
		if v.Type > 4 || v.Type < 1{
			return 0, errors.Errorf("进出区域或路线，位置值错误:%d", v.Type)
		}
		if v.Direction > 1{
			return 0, errors.Errorf("进出区域或路线，进出值错误:%d", v.Direction)
		}
		var directionMap = map[byte]string{0:"进",1:"出"}
		fmt.Println("进出区域或路线：", posTypeMap[v.Type], v.Id, directionMap[v.Direction])
		break
	case 0x13:
		//路段行驶时间不足或太长报警信息
		var v PosDriveTimeAddition
		err = common.BRead(buf, common.BigEndian, &v)
		if err != nil {
			return 0, err
		}
		if v.Result > 1{
			return 0, errors.Errorf("路段行驶时间不足或太长报警信息，结果值错误:%d", v.Result)
		}
		var resMap = map[byte]string{0:"不足",1:"太长"}
		fmt.Println("路段行驶时间不足或太长：", v.Id, v.Time, resMap[v.Result])
		break
	case 0x25:
		//扩展车辆信号状态
		var v uint32
		err = common.BRead(buf, common.BigEndian, &v)
		if err != nil {
			return 0, err
		}
		var ext PosExtStatusAddition
		ext.Parse(v)
		ext.Show()
		break
	case 0x2A:
		//IO状态
		var v uint16
		err = common.BRead(buf, common.BigEndian, &v)
		if err != nil {
			return 0, err
		}
		var ios PosIOStatusAddition
		ios.Parse(v)
		ios.Show()
		break
	case 0x2B:
		//模拟量
		var v uint32
		err = common.BRead(buf, common.BigEndian, &v)
		if err != nil {
			return 0, err
		}
		var a PosAnalogAddition
		a.Parse(v)
		a.Show()
		break
	case 0x30:
		//无线通信网络信号强度
		var v byte
		err = readAddition(buf, &v, "无线通信网络信号强度:")
		break
	case 0x31:
		//GNSS定位卫星数
		var v byte
		err = readAddition(buf, &v, "GNSS定位卫星数:")
		break
	}
	if err != nil {
		return 0, err
	}
	l := 2 + p.Length
	return uint16(l), nil
}