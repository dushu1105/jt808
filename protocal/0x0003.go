package protocal

const TRegistCancelRequest = 0x0003

type TRegistCancelHandler struct {
	BaseHandler
}

func (t *TRegistCancelHandler) Cancel() byte{
	//todo
	return 0
}

func (t *TRegistCancelHandler) Do(msg *JT808Msg) (*Jt808ResultMsg, error) {
	var err error
	r := t.Cancel()

	v := CommonResp{Seq: msg.Header.Seq, Result:r}
	ret, err := v.Packet()
	if err != nil{
		return nil, err
	}

	return &Jt808ResultMsg{Msg:msg.CopyAndSet(PCommonResponse, ret), NeedFeedBack:true}, err
}
