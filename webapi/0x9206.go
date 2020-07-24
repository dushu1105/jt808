package webapi

import (
	"github.com/dushu1105/jt808/common"
	"github.com/dushu1105/jt808/jtnet"
	"github.com/dushu1105/jt808/protocal"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Jt9206Req struct {
	Sim string
	Ip string
	Port uint16
	User string `ref:"UserLen"`
	Passwd string `ref:"PasswdLen"`
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

func Jt808_9206(c *gin.Context) {
	fmt.Println("Jt808_9101")
	s, ok := c.Get(SERVER_KEY)
	if !ok{
		fmt.Println(ERR_PROC_INITING)
		c.String(http.StatusInternalServerError, ERR_DEVICE_CONN_NOT_READY)
		return
	}
	var req Jt9206Req
	err := c.BindJSON(&req)
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	var jt protocal.PUploadHandler
	jt.Channel = req.Channel
	jt.DataType = req.DataType
	jt.Flag = req.Flag
	jt.ExcuteCondition = req.ExcuteCondition
	jt.Ip = req.Ip
	jt.StreamType = req.StreamType
	jt.StorageType = req.StorageType

	d, err := common.Utf8ToGbk([]byte(req.Ip))
	if err != nil{
		fmt.Println(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	jt.IPLen = uint8(len(d))

	d, err = common.Utf8ToGbk([]byte(req.User))
	if err != nil{
		fmt.Println(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	jt.UserLen = uint8(len(d))

	d, err = common.Utf8ToGbk([]byte(req.Passwd))
	if err != nil{
		fmt.Println(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	jt.PasswdLen = uint8(len(d))

	d, err = common.Utf8ToGbk([]byte(req.Path))
	if err != nil{
		fmt.Println(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	jt.PathLen = uint8(len(d))


	buf, err := jt.Packet()
	if err != nil{
		fmt.Println(err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	err = SendDevice(s.(*jtnet.Server), req.Sim, protocal.PUploadRequest, 0, buf)
	if err != nil{
		fmt.Println(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"status":     "posted",
	})
}

