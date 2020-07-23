package protocal

import (
	"time"
)

const TQueryTimeRequest = 0x0004
const PQueryTimeResponse = 0x8004

type TQueryTimeHandler struct {
	BaseHandler
}

type PQueryTimeResp struct {
	Time  [6]byte
}

func dec2bcd(dec int) uint8 {
	return uint8(((dec/10)<<4) + (dec%10));
}

func Time() [6]byte {
	t := time.Now().UTC()//.Format("2006-01-02 15:04:05")

	var r [6]byte

	r[0] = dec2bcd(t.Year()%100)
	r[1] = dec2bcd(int(t.Month()))
	r[2] = dec2bcd(t.Day())
	r[3] = dec2bcd(t.Hour())
	r[4] = dec2bcd(t.Minute())
	r[5] = dec2bcd(t.Second())
	return r
}

func (t *TQueryTimeHandler) Do(msg *JT808Msg) (*JT808Msg, error) {
	var err error
	r := Time()
	return msg.CopyAndSet(PQueryTimeResponse, r[:]), err
}
