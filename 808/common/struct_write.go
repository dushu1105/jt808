package common

import (
	"bytes"
	"fmt"
	"math"
	"reflect"
	"errors"
	"strings"
)

func writeNormal(w *bytes.Buffer, order ByteOrder, data interface{}, n int) error {
	bs := make([]byte, n)
	switch v := data.(type) {
	case *bool:
		if *v {
			bs[0] = 1
		} else {
			bs[0] = 0
		}
	case bool:
		if v {
			bs[0] = 1
		} else {
			bs[0] = 0
		}
	case []bool:
		for i, x := range v {
			if x {
				bs[i] = 1
			} else {
				bs[i] = 0
			}
		}
	case *int8:
		bs[0] = byte(*v)
	case int8:
		bs[0] = byte(v)
	case []int8:
		for i, x := range v {
			bs[i] = byte(x)
		}
	case *uint8:
		bs[0] = *v
	case uint8:
		bs[0] = v
	case []uint8:
		bs = v // TODO(josharian): avoid allocating bs in this case?
	case *int16:
		order.PutUint16(bs, uint16(*v))
	case int16:
		order.PutUint16(bs, uint16(v))
	case []int16:
		for i, x := range v {
			order.PutUint16(bs[2*i:], uint16(x))
		}
	case *uint16:
		order.PutUint16(bs, *v)
	case uint16:
		order.PutUint16(bs, v)
	case []uint16:
		for i, x := range v {
			order.PutUint16(bs[2*i:], x)
		}
	case *int32:
		order.PutUint32(bs, uint32(*v))
	case int32:
		order.PutUint32(bs, uint32(v))
	case []int32:
		for i, x := range v {
			order.PutUint32(bs[4*i:], uint32(x))
		}
	case *uint32:
		order.PutUint32(bs, *v)
	case uint32:
		order.PutUint32(bs, v)
	case []uint32:
		for i, x := range v {
			order.PutUint32(bs[4*i:], x)
		}
	case *int64:
		order.PutUint64(bs, uint64(*v))
	case int64:
		order.PutUint64(bs, uint64(v))
	case []int64:
		for i, x := range v {
			order.PutUint64(bs[8*i:], uint64(x))
		}
	case *uint64:
		order.PutUint64(bs, *v)
	case uint64:
		order.PutUint64(bs, v)
	case []uint64:
		for i, x := range v {
			order.PutUint64(bs[8*i:], x)
		}
	case *float32:
		order.PutUint32(bs, math.Float32bits(*v))
	case float32:
		order.PutUint32(bs, math.Float32bits(v))
	case []float32:
		for i, x := range v {
			order.PutUint32(bs[4*i:], math.Float32bits(x))
		}
	case *float64:
		order.PutUint64(bs, math.Float64bits(*v))
	case float64:
		order.PutUint64(bs, math.Float64bits(v))
	case []float64:
		for i, x := range v {
			order.PutUint64(bs[8*i:], math.Float64bits(x))
		}
	}
	_, err := w.Write(bs)
	return err
}

func (e *encoder) getLenByTag(v reflect.Value, t reflect.StructTag) (int, error){
	var err error

	vlen, err := getLenByTag(v, t)
	if err != nil{
		return -1, err
	}

	if vlen == ALL_DATA_LEN{
		vlen = e.end - e.offset
	}

	if e.offset+vlen > e.end  || e.offset >= e.end {
		return -1, errors.New(ERR_LENGTH_NOT_ENOUGH)
	}
	return vlen, nil
}

func (e *encoder) transformByTag(t reflect.StructTag, buf string, bufLen int) ([]byte, error){
	vLenStr, ok := t.Lookup("type")
	if ok{
		if vLenStr == "bcd" ||  vLenStr == "time"{
			data := Str2BCD(buf)
			return data, nil
		}
	}

	data, err := Utf8ToGbk([]byte(buf))
	if err != nil{
		return nil, err
	}
	return data, nil
}

func (e *encoder) value1(v reflect.Value) error {
	var err error
	switch v.Kind() {
	case reflect.Ptr:
		err = e.value1(v.Elem())
		if err != nil{
			return err
		}
	case reflect.Array, reflect.Slice:
		l := v.Len()
		for i := 0; i < l; i++ {
			err = e.value1(v.Index(i))
			if err != nil{
				return err
			}
		}

	case reflect.Struct:
		t := v.Type()
		l := v.NumField()
		for i := 0; i < l; i++ {
			// see comment for corresponding code in decoder.value()
			v1 := v.Field(i)
			if strings.HasPrefix(t.Field(i).Name, "_"){
				continue
			}

			if v1.Kind() == reflect.String{
				vlen, err := e.getLenByTag(v, t.Field(i).Tag)
				if err != nil{
					return err
				}

				tmp, err := e.transformByTag(t.Field(i).Tag, v1.String(), vlen)
				if err != nil{
					fmt.Println("transfer data to string failed", t.Field(i).Name)
					return err
				}

				for j:=0;j<vlen;j++{
					if j < len(tmp){
						e.buf[e.offset+j] = tmp[j]
					} else {
						e.buf[e.offset+j] = 0
					}
				}
				e.offset += vlen

			} else {
				err = e.value1(v1)
			}
			if err != nil{
				return err
			}
		}
	case reflect.Bool:
		e.bool(v.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch v.Type().Kind() {
		case reflect.Int8:
			e.int8(int8(v.Int()))
		case reflect.Int16:
			e.int16(int16(v.Int()))
		case reflect.Int32:
			e.int32(int32(v.Int()))
		case reflect.Int64:
			e.int64(v.Int())
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		switch v.Type().Kind() {
		case reflect.Uint8:
			e.uint8(uint8(v.Uint()))
		case reflect.Uint16:
			e.uint16(uint16(v.Uint()))
		case reflect.Uint32:
			e.uint32(uint32(v.Uint()))
		case reflect.Uint64:
			e.uint64(v.Uint())
		}

	case reflect.Float32, reflect.Float64:
		switch v.Type().Kind() {
		case reflect.Float32:
			e.uint32(math.Float32bits(float32(v.Float())))
		case reflect.Float64:
			e.uint64(math.Float64bits(v.Float()))
		}

	case reflect.Complex64, reflect.Complex128:
		switch v.Type().Kind() {
		case reflect.Complex64:
			x := v.Complex()
			e.uint32(math.Float32bits(float32(real(x))))
			e.uint32(math.Float32bits(float32(imag(x))))
		case reflect.Complex128:
			x := v.Complex()
			e.uint64(math.Float64bits(real(x)))
			e.uint64(math.Float64bits(imag(x)))
		}
	}
	return nil
}

func WriteStruct(w *bytes.Buffer, order ByteOrder, data interface{}) error {
	if n := intDataSize(data); n != 0 {
		err := writeNormal(w, order, data, n)
		if err != nil{
			return err
		}
	}

	v := reflect.Indirect(reflect.ValueOf(data))
	size := sizeof1(v)
	if size < 0 {
		return errors.New("WriteStruct: invalid type " + reflect.TypeOf(data).String())
	}
	buf := make([]byte, size)
	e := &encoder{order: order, buf: buf, end:size}
	err := e.value1(v)
	if err != nil{
		return err
	}

	_, err = w.Write(buf)
	return err
}

// sizeof returns the size >= 0 of variables for the given type or -1 if the type is not acceptable.
func sizeof1(v reflect.Value) int {
	switch v.Kind() {
	case reflect.Ptr:
		s := sizeof1(v.Elem())
		if s < 0 {
			return -1
		}
		return s
	case reflect.Array,reflect.Slice:
		l := v.Len()
		sum := 0
		for i := 0; i < l; i++ {
			s := sizeof1(v.Index(i))
			if s < 0{
				return -1
			}
			sum += s
		}
		return sum
	case reflect.Struct:
		t := v.Type()
		l := v.NumField()

		sum := 0
		for i := 0; i < l; i++ {
			v1 := v.Field(i)
			if strings.HasPrefix(t.Field(i).Name, "_"){
				continue
			}

			if v1.Kind() == reflect.String{
				s, err := getLenByTag(v, t.Field(i).Tag)
				if err != nil{
					fmt.Println("length < 0", t.Field(i).Type, t.Field(i).Name)
					return -1
				}
				if s == ALL_DATA_LEN{
					s = v1.Len()
				}
				sum += s
			} else {
				s := sizeof1(v.Field(i))
				if s < 0 {
					fmt.Println("length < 0", t.Field(i).Type, t.Field(i).Name)
					return -1
				}
				sum += s
			}

		}
		return sum
	case reflect.Bool,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		return int(v.Type().Size())
	case reflect.String:
		return v.Len()
	}

	return -1
}