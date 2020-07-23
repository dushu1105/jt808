package protocal

const PQueryAttrRequest = 0x8107

type PQueryAttrHandler struct{}

func (p *PQueryAttrHandler) Packet() ([]byte, error) {
	return nil, nil
}