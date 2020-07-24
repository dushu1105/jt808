package webapi

import (
	"github.com/dushu1105/jt808/common"
	"github.com/dushu1105/jt808/jtnet"
	"github.com/dushu1105/jt808/protocal"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Jt9101Req struct {
	Sim string
	Ip string
	TcpPort uint16
	UdpPort uint16
	Channel byte
	DataType byte //0音频， 1视频，2双向对讲，3监听， 4中心广播， 5透出
	StreamType byte //0主码流，1子码流
}

func Jt808_9101(c *gin.Context) {
	fmt.Println("Jt808_9101")
	s, ok := c.Get(SERVER_KEY)
	if !ok{
		fmt.Println(ERR_PROC_INITING)
		c.String(http.StatusInternalServerError, ERR_DEVICE_CONN_NOT_READY)
		return
	}
	var req Jt9101Req
	err := c.BindJSON(&req)
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	var jt protocal.PStreamReq
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
	buf, err := jt.Packet()
	if err != nil{
		fmt.Println(err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	err = SendDevice(s.(*jtnet.Server), req.Sim, protocal.PStreamRequest, 0, buf)
	if err != nil{
		fmt.Println(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"status":     "posted",
	})
}

