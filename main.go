package main

import (
	"808/jtnet"
	"808/protocal"
)

func addHandler(s *jtnet.Server){
	s.AddHandler(protocal.TCommonResponse, &protocal.CommonRespHandler{})
	s.AddHandler(protocal.TRegistRequest2013, &protocal.TRegistReqHandler2013{})
	s.AddHandler(protocal.TRegistRequest2019, &protocal.TRegistReqHandler2019{})
	s.AddHandler(protocal.THeartBeat, &protocal.HeartBeatHandler{})
	s.AddHandler(protocal.TVerifyRequest2013, &protocal.TVerifyHandler2013{})
	s.AddHandler(protocal.TVerifyRequest2019, &protocal.TVerifyHandler2019{})
	s.AddHandler(protocal.TRegistCancelRequest, &protocal.TRegistCancelHandler{})
	s.AddHandler(protocal.TQueryTimeRequest, &protocal.TQueryTimeHandler{})
	s.AddHandler(protocal.TPositionRequest, &protocal.TPositionHandler{})
	s.AddHandler(protocal.TPositionRequest1, &protocal.TPositionHandler{})
	s.AddHandler(protocal.TQueryAttrResponse, &protocal.TQueryAttrHandler{})
	s.AddHandler(protocal.TDriverInfoResponse, &protocal.TDriverInfoHandler{})
}

func main() {
	//1 创建一个server句柄
	s := jtnet.NewServer()

	//2 配置路由
	addHandler(s)

	//3 开启服务
	s.Serve()
}
