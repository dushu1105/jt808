package webapi

import (
	"github.com/dushu1105/jt808/jtnet"
	"github.com/dushu1105/jt808/protocal"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type jt8300Request struct {
	Sim string
	Type uint8
	Flag uint8
	Txt string
	Ver uint8
}

func Jt808_8300(c *gin.Context) {
	fmt.Println("Jt808_8300")
	s, ok := c.Get(SERVER_KEY)
	if !ok{
		fmt.Println(ERR_PROC_INITING)
		c.String(http.StatusInternalServerError, ERR_DEVICE_CONN_NOT_READY)
		return
	}
	var req jt8300Request
	err := c.BindJSON(&req)
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	var buf []byte
	if req.Ver == 1 {
		var jt protocal.PTextSendReqHandler2019
		jt.Txt = req.Txt
		jt.Flag = jt.Flag
		jt.Type = req.Type
		buf, err = jt.Packet()
		if err != nil{
			fmt.Println(err)
			c.String(http.StatusBadRequest, err.Error())
			return
		}
	} else {
		var jt protocal.PTextSendReqHandler2013
		jt.Txt = req.Txt
		jt.Flag = jt.Flag
		buf, err = jt.Packet()
		if err != nil{
			fmt.Println(err)
			c.String(http.StatusBadRequest, err.Error())
			return
		}
	}


	err = SendDevice(s.(*jtnet.Server), req.Sim, protocal.PTextSendRequest2013, 0, buf)
	if err != nil{
		fmt.Println(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"status":     "posted",
	})
}

