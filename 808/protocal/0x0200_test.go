package protocal

import (
	"testing"
)

func TestTPositionHandler_Parse(t *testing.T) {
	var data = "7e 02 00 40 1c 01 00 00 00 00 02 22 22 22 22 22 00 03 00 00 00 05 00 00 00 0a 00 01 ad b0 00 01 86 a0 00 32 00 32 00 1e 20 07 16 11 31 03 44 7e"
	var j JT808Msg
	err := j.Parse(s2b(data))
	if err != nil{
		t.Log(err)
	}
	t.Log(j.Header)

	var tt TPositionHandler

	err = tt.Parse(j.Body)
	if err != nil{
		t.Log(err)
		return
	}
	t.Log(tt)
	tt.Do(&j)
}
