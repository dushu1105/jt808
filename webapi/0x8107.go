package webapi

import (
	"github.com/dushu1105/jt808/jtnet"
	"github.com/dushu1105/jt808/protocal"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Jt808_8107(c *gin.Context) {
	fmt.Println("Jt808_8107")
	s, ok := c.Get(SERVER_KEY)
	if !ok{
		fmt.Println("没有发送方式")
		c.String(http.StatusInternalServerError, "no send method")
		return
	}
	var req JtCommonRequet
	err := c.BindJSON(&req)
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	var jt8107 protocal.PQueryAttrHandler
	buf, err := jt8107.Packet()
	if err != nil{
		fmt.Println(err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	err = SendDevice(s.(*jtnet.Server), req.Sim, protocal.PQueryAttrRequest, 0, buf)
	if err != nil{
		fmt.Println(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"status":     "posted",
	})
}
