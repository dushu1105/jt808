package webapi

import (
"github.com/dushu1105/jt808/jtnet"
"github.com/dushu1105/jt808/protocal"
"fmt"
"github.com/gin-gonic/gin"
"net/http"
)

func Jt808_8702(c *gin.Context) {
	fmt.Println("Jt808_8702")
	s, ok := c.Get(SERVER_KEY)
	if !ok{
		fmt.Println(ERR_PROC_INITING)
		c.String(http.StatusInternalServerError, ERR_DEVICE_CONN_NOT_READY)
		return
	}
	var req JtCommonRequet
	err := c.BindJSON(&req)
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	var jt8702 protocal.PDriverInfoHandler
	buf, err := jt8702.Packet()
	if err != nil{
		fmt.Println(err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	err = SendDevice(s.(*jtnet.Server), req.Sim, protocal.PDriverInfoRequest, 0, buf)
	if err != nil{
		fmt.Println(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"status":     "posted",
	})
}

