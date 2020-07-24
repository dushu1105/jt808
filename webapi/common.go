package webapi

import (
	"fmt"
	"github.com/dushu1105/jt808/jtnet"
	"github.com/pkg/errors"
	"net/http"
	"github.com/gin-gonic/gin"
)

const ERR_PROC_INITING = "程序初始化未完成"
const ERR_RESULT_PENDING = "结果还未返回，请等待"
const ERR_DEVICE_CONN_NOT_READY = "与设备的连接未准备好"
type JtCommonRequet struct {
	Sim string
}

func SendDevice(s *jtnet.Server, sim string, id uint16, fragFlag uint8, data []byte) error {
	return s.SendBySim(sim, id, fragFlag, data)
}

type WaitResponse struct{
	Sim string
	Id  uint16
	Seq uint16
}

func (w *WaitResponse)getCommonResp(s *jtnet.Server) (byte, error) {
	c, err := s.ConnMgr.GetBySim(w.Sim)
	if err != nil{
		return 0, err
	}
	r, ok := c.GetRespResult(w.Id, w.Seq)
	if !ok{
		return 0, errors.Errorf(ERR_RESULT_PENDING)
	}

	return r, nil
}

func WaitCommonResp(c *gin.Context) {
	fmt.Println("WaitCommonResp")
	s, ok := c.Get(SERVER_KEY)
	if !ok {
		fmt.Println(ERR_PROC_INITING)
		c.String(http.StatusInternalServerError, ERR_DEVICE_CONN_NOT_READY)
		return
	}
	var req WaitResponse
	err := c.BindJSON(&req)
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	r, err := req.getCommonResp(s.(*jtnet.Server))
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"result":     r,
	})
}