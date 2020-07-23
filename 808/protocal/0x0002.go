package protocal

const THeartBeat = 0x0002

type HeartBeatHandler struct {
	BaseHandler
}

func (h *HeartBeatHandler) Do(msg *JT808Msg) (*JT808Msg, error) {
	var err error

	v := CommonRespHandler{Seq:msg.Header.Seq, Result:0}
	ret, err := v.Packet()
	if err != nil{
		return nil, err
	}

	return msg.CopyAndSet(PCommonResponse, ret), err
}
