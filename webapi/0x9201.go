package webapi

import (
	"github.com/dushu1105/jt808/common"
	"github.com/dushu1105/jt808/jtnet"
	"github.com/dushu1105/jt808/protocal"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Jt9201Req struct {
	Sim string
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

func Jt808_9201(c *gin.Context) {
	fmt.Println("Jt808_9101")
	s, ok := c.Get(SERVER_KEY)
	if !ok{
		fmt.Println(ERR_PROC_INITING)
		c.String(http.StatusInternalServerError, ERR_DEVICE_CONN_NOT_READY)
		return
	}
	var req Jt9201Req
	err := c.BindJSON(&req)
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	var jt protocal.PReplayHandler
	jt.Channel = req.Channel
	jt.DataType = req.DataType
	jt.TcpPort = req.TcpPort
	jt.UdpPort = req.UdpPort
	jt.Ip = req.Ip
	jt.StreamType = req.StreamType
	d, err := common.Utf8ToGbk([]byte(req.Ip))
	if err != nil{
		fmt.Println(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	jt.IPLen = uint8(len(d))
	jt.ForwardMutiple = req.ForwardMutiple
	jt.ReplayType = req.ReplayType
	jt.StartTime = req.StartTime
	jt.StorageType = req.StorageType
	jt.EndTime = req.EndTime

	buf, err := jt.Packet()
	if err != nil{
		fmt.Println(err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	err = SendDevice(s.(*jtnet.Server), req.Sim, protocal.PReplayRequest, 0, buf)
	if err != nil{
		fmt.Println(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"status":     "posted",
	})
}

