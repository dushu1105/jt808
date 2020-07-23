package protocal

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func s2b(ss string) []byte{
	b := make([]byte, 0)
	for _, s := range strings.Split(ss, " "){
		i, err := strconv.ParseInt(s, 16, 0)
		if err != nil{
			fmt.Println(err)
			return nil
		}
		b = append(b, uint8(i))
	}
	return b
}

func TestTRegistReqHandler_Do(t *testing.T) {
	var data = "7e 01 00 00 31 01 11 11 11 11 11 00 01 00 00 00 00 4b 6f 69 6b 65 30 30 31 4b 6f 69 6b 65 30 30 31 30 30 30 30 30 30 30 30 30 30 30 30 4b 6f 69 6b 65 30 30 31 02 d4 c1 41 31 32 33 34 35 34 7e"
	//var data = "7e 01 00 40 54 01 00 00 00 00 02 22 22 22 22 22 00 01 00 00 00 00 30 30 30 4b 6f 69 6b 65 30 30 32 30 30 30 30 30 30 30 30 30 30 30 30 30 30 30 30 30 30 30 30 30 30 4b 6f 69 6b 65 30 30 32 30 30 30 30 30 30 30 30 30 30 30 30 30 30 30 30 30 30 30 30 30 30 4b 6f 69 6b 65 30 30 32 02 d4 c1 41 31 32 33 34 36 10 7e"
	var j JT808Msg
	err := j.Parse(s2b(data))
	if err != nil{
		t.Log(err)
	}
	t.Log(j.Header)

	var tt TRegistReqHandler2013
	err = tt.Parse(j.Body)
	if err != nil{
		t.Log(err)
	}
	t.Log(string(tt.Manu[:]), string(tt.TType[:]), string(tt.Licence[:]),)
}

func TestTVerifyHandler2019_Parse(t *testing.T) {
	var data = "7e 01 02 40 1f 01 00 00 00 00 02 22 22 22 22 22 00 02 04 36 36 36 36 31 32 33 34 35 36 76 31 2e 30 2e 30 30 30 30 30 30 30 30 30 30 30 30 30 30 30 3b 7e"
	var j JT808Msg
	err := j.Parse(s2b(data))
	if err != nil{
		t.Log(err)
	}
	t.Log(j.Header)

	var tt TVerifyHandler2019

	err = tt.Parse(j.Body)
	if err != nil{
		t.Log(err)
		return
	}
	t.Log(tt)
}

func TestTime(t *testing.T) {
	fmt.Printf("% x\n", Time())

	fmt.Printf("% x", 0x200716)
	fmt.Printf("% x\n", []byte{0b00100000, 0x07, 0x16})

}
