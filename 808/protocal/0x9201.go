package protocal

import (
	"808/common"
	"bytes"
)

const PReplayRequest = 0x9201

type PReplayHandler struct {
	IPLen byte
	Ip string `ref:"IPLen"`
	TcpPort uint16
	UdpPort uint16
	Channel byte
	DataType byte //0音视频，1音频， 2视频，3 视频或音视频
	StreamType byte //0 主或子，只音频， 1主码流，2子码流
	StorageType byte //0主或备 1 主存 2备存
	ReplayType byte //0正常，1快进 2关键帧快退回放 3关键帧播放 4单帧上传
	ForwardMutiple byte //2的减一幂次， 0无效，最大5， ReplayType=1或2才有效，否则0
	StartTime  string `len:"6" type:"time"`
	EndTime    string `len:"6" type:"time"`
}

func (p *PReplayHandler) Packet() ([]byte, error){
	buf := new(bytes.Buffer)
	err := common.WriteStruct(buf, common.BigEndian, p)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
