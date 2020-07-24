package protocal

const THeartBeat = 0x0002

type HeartBeatHandler struct {
	BaseHandler
}

func (h *HeartBeatHandler) Do(msg *JT808Msg) (*Jt808ResultMsg, error) {
	var err error

	v := CommonResp{Seq: msg.Header.Seq, Result:0}
	ret, err := v.Packet()
	if err != nil{
		return nil, err
	}

	return &Jt808ResultMsg{Msg:msg.CopyAndSet(PCommonResponse, ret), NeedFeedBack:true}, err
}
