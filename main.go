package main

import (
	"github.com/dushu1105/jt808/jtnet"
	"github.com/dushu1105/jt808/protocal"
	"github.com/dushu1105/jt808/webapi"
)

func addHandler(s *jtnet.Server){
	s.AddHandler(protocal.TCommonResponse, &protocal.CommonResp{})
	s.AddHandler(protocal.TRegistRequest2013, &protocal.TRegistReqHandler2013{})
	s.AddHandler(protocal.TRegistRequest2019, &protocal.TRegistReqHandler2019{})
	s.AddHandler(protocal.THeartBeat, &protocal.HeartBeatHandler{})
	s.AddHandler(protocal.TVerifyRequest2013, &protocal.TVerifyHandler2013{})
	s.AddHandler(protocal.TVerifyRequest2019, &protocal.TVerifyHandler2019{})
	s.AddHandler(protocal.TRegistCancelRequest, &protocal.TRegistCancelHandler{})
	s.AddHandler(protocal.TQueryTimeRequest, &protocal.TQueryTimeHandler{})
	s.AddHandler(protocal.TPositionRequest, &protocal.TPositionHandler{})
	s.AddHandler(protocal.TPositionRequest1, &protocal.TPositionHandler{})
	s.AddHandler(protocal.TQueryAttrResponse2013, &protocal.TQueryAttrHandler2013{})
	s.AddHandler(protocal.TQueryAttrResponse2019, &protocal.TQueryAttrHandler2019{})
	s.AddHandler(protocal.TDriverInfoResponse, &protocal.TDriverInfoHandler{})
}

func main() {
	//1 创建一个server句柄
	s := jtnet.NewServer()

	//2 配置路由
	addHandler(s)

	go webapi.RunWebServer(s)
	//3 开启服务
	s.Serve()
}
