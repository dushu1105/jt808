package protocal

import (
	"bytes"
	"github.com/dushu1105/jt808/common"
)

const PUploadRequest = 0x9206


type PUploadHandler struct {
	IPLen byte
	Ip string `ref:"IPLen"`
	Port uint16
	UserLen byte
	User string `ref:"UserLen"`
	PasswdLen byte
	Passwd string `ref:"PasswdLen"`
	PathLen byte
	Path string `ref:"PathLen"`
	Channel byte
	StartTime  string `len:"6" type:"time"`
	EndTime    string `len:"6" type:"time"`
	Flag uint64 //todo
	DataType byte //0音视频，1音频， 2视频，3 视频或音视频
	StreamType byte //0 主或子，只音频， 1主码流，2子码流
	StorageType byte //0主或备 1 主存 2备存
	ExcuteCondition byte //todo
}

func (p *PUploadHandler) Packet() ([]byte, error){
	buf := new(bytes.Buffer)
	err := common.WriteStruct(buf, common.BigEndian, p)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}