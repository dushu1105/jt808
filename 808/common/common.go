package common

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
)

func Encode(data []byte) ([]byte, error) {
	var begin int

	buf := new(bytes.Buffer)
	err := buf.WriteByte(0x7e)
	if err != nil {
		return nil, err
	}
	for i, d := range data{
		if d == 0x7e{
			err = BWrite(buf, BigEndian, data[begin:i])
			if err != nil {
				return nil, err
			}
			err = BWrite(buf, BigEndian, []byte{0x7d, 0x02})
			if err != nil {
				return nil, err
			}
			begin = i + 1
		} else if d == 0x7d {
			err = BWrite(buf, BigEndian, data[begin:i])
			if err != nil {
				return nil, err
			}
			err = BWrite(buf, BigEndian, []byte{0x7d, 0x01})
			if err != nil {
				return nil, err
			}
			begin = i + 1
		}
	}
	err = BWrite(buf, BigEndian, data[begin:])
	if err != nil {
		return nil, err
	}
	err = buf.WriteByte(0x7e)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Decode(data []byte) ([]byte, error) {
	if len(data) < 1{
		return nil, errors.Errorf("空数据")
	}
	var begin, p int
	var err error
	begin = 1
	buf := new(bytes.Buffer)
	for i, d := range data[1:len(data) - 1]{
		j := i + 1
		if d == 0x7d{
			p = 1
			continue
		}

		if p == 1{
			if d == 0x01{
				err = BWrite(buf, BigEndian, data[begin:j])
				if err != nil {
					return nil, err
				}
				begin = j + 1
			} else if d == 0x02{
				err = BWrite(buf, BigEndian, data[begin:j-1])
				if err != nil {
					return nil, err
				}
				err = buf.WriteByte(0x7e)
				if err != nil {
					return nil, err
				}
				begin = j + 1
			}
		}
	}

	err = BWrite(buf, BigEndian, data[begin:len(data) - 1])
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}



func encrypt(data []byte) ([]byte, error){
	//808头的encrypt位如果是001，就需要做RSA加密，此处待做，000表示不加密
	return data, nil
}

func decrypt(data []byte) ([]byte, error){
	return data, nil
}

func BCD2Time(bcd []byte) string {
	var d string
	for _, s := range bcd {
		v := uint8(s)
		tmp := fmt.Sprintf("%d", v >> 4)
		d += tmp
		tmp = fmt.Sprintf("%d", v & 0b00001111)
		d += tmp
	}
	return d
}

func BCD2DEC(bcd []byte) string {
	var d string
	for _, s := range bcd {
		v := uint8(s)
		tmp := fmt.Sprintf("%d", v >> 4)
		d += tmp
		tmp = fmt.Sprintf("%d", v & 0b00001111)
		d += tmp
	}

	for i, s := range d{
		if s != int32('0'){
			d = d[i:]
			break
		}
	}
	return d
}

func Str2BCD(s string) []byte {
	if len(s) % 2 != 0{
		s = "0" + s
	}
	ret := make([]byte, len(s))
	for i := 0; i < len(s); i += 2 {
		a := int(s[i]) - int('0')
		b := int(s[i+1]) - int('0')
		ret[i/2] = uint8(a<<4 | b)
	}
	return ret
}

func DEC2BCD(dec string) []byte {
	var r = make([]byte, 0)
	for _, s := range dec{
		v := uint8(s)
		r = append(r, (v+(v/10)*6))
	}
	return  r
}

func Parse(data []byte, h interface{})  error {
	buf := bytes.NewBuffer(data)
	err := BRead(buf, BigEndian, h)
	return err
}
func Packet(h interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := BWrite(buf, BigEndian, h)
	if err != nil{
		return nil, err
	}
	return buf.Bytes(), nil
}

func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}