package protocal

import "testing"

func TestCommonRespHandler(t *testing.T) {
	var c CommonResp
	var v = []byte{0x00, 0x01, 0x00, 0x02, 0x03}

	err := c.Parse(v)
	if err != nil{
		t.Log(err.Error())
		return
	}
	t.Log(c)

	d, err := c.Packet()
	if err != nil{
		t.Log(err.Error())
		return
	}
	t.Logf("% x", d)
}
