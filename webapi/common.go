package webapi

import (
	"github.com/dushu1105/jt808/jtnet"
)

type JtCommonRequet struct {
	Sim string
}

func SendDevice(s *jtnet.Server, sim string, id uint16, fragFlag uint8, data []byte) error {
	return s.SendBySim(sim, id, fragFlag, data)
}