package protocal

const TRegistCancelRequest = 0x0003

type TRegistCancelHandler struct {
	BaseHandler
}

func (t *TRegistCancelHandler) Cancel() byte{
	//todo
	return 0
}

func (t *TRegistCancelHandler) Do(msg *JT808Msg) (*JT808Msg, error) {
	var err error
	r := t.Cancel()

	v := CommonRespHandler{Seq:msg.Header.Seq, Result:r}
	ret, err := v.Packet()
	if err != nil{
		return nil, err
	}

	return msg.CopyAndSet(PCommonResponse, ret), err
}
