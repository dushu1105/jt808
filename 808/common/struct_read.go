package common

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

const ALL_DATA_LEN  = -2
const ERR_LENGTH_NOT_ENOUGH = "data length not enough for struct"
const ERR_TAG = "member type string must specify length member in tag or pos as last member"

func readNormal(bs []byte, order ByteOrder, data interface{}) int {
	switch data := data.(type) {
	case *bool:
		*data = bs[0] != 0
	case *int8:
		*data = int8(bs[0])
	case *uint8:
		*data = bs[0]
	case *int16:
		*data = int16(order.Uint16(bs))
	case *uint16:
		*data = order.Uint16(bs)
	case *int32:
		*data = int32(order.Uint32(bs))
	case *uint32:
		*data = order.Uint32(bs)
	case *int64:
		*data = int64(order.Uint64(bs))
	case *uint64:
		*data = order.Uint64(bs)
	case *float32:
		*data = math.Float32frombits(order.Uint32(bs))
	case *float64:
		*data = math.Float64frombits(order.Uint64(bs))
	case []bool:
		for i, x := range bs { // Easier to loop over the input for 8-bit values.
			data[i] = x != 0
		}
	case []int8:
		for i, x := range bs {
			data[i] = int8(x)
		}
	case []uint8:
		copy(data, bs)
	case []int16:
		for i := range data {
			data[i] = int16(order.Uint16(bs[2*i:]))
		}
	case []uint16:
		for i := range data {
			data[i] = order.Uint16(bs[2*i:])
		}
	case []int32:
		for i := range data {
			data[i] = int32(order.Uint32(bs[4*i:]))
		}
	case []uint32:
		for i := range data {
			data[i] = order.Uint32(bs[4*i:])
		}
	case []int64:
		for i := range data {
			data[i] = int64(order.Uint64(bs[8*i:]))
		}
	case []uint64:
		for i := range data {
			data[i] = order.Uint64(bs[8*i:])
		}
	case []float32:
		for i := range data {
			data[i] = math.Float32frombits(order.Uint32(bs[4*i:]))
		}
	case []float64:
		for i := range data {
			data[i] = math.Float64frombits(order.Uint64(bs[8*i:]))
		}
	default:
		return 0
	}

	return len(bs)
}
func ReadStruct(buf []byte, order ByteOrder, data interface{}) error {
	l := len(buf)
	if n := intDataSize(data); n != 0 {
		if n > l{
			return errors.New(ERR_LENGTH_NOT_ENOUGH)
		}
		n = readNormal(buf[:n], order, data)
		if n != 0{
			return nil
		}
	}

	v := reflect.ValueOf(data)
	d := &decoder{order: order, buf: buf, end:len(buf)}
	err := d.value1(v)

	return err
}

func getLenByTag(v reflect.Value, t reflect.StructTag) (int, error){
	var err error
	var vlen = -1

	vLenStr, ok := t.Lookup("len")
	if ok{
		vlen, err = strconv.Atoi(vLenStr)
		if err != nil{
			return -1, err
		}
	} else {
		vLenStr, ok = t.Lookup("ref")
		if ok{
			lenField := v.FieldByName(vLenStr)
			if !lenField.IsValid(){
				return -1, errors.New(ERR_TAG)
			}
			vlen = int(lenField.Uint())
		} else {
			return ALL_DATA_LEN, nil
		}
	}

	if vlen < 0{
		return -1, errors.New(ERR_TAG)
	}
	return vlen, nil
}
func (d *decoder) getLenByTag(v reflect.Value, t reflect.StructTag) (int, error){
	var err error

	vlen, err := getLenByTag(v, t)
	if err != nil{
		return -1, err
	}

	if vlen == ALL_DATA_LEN{
		vlen = d.end - d.offset
	}

	if d.offset+vlen > d.end  || d.offset >= d.end {
		return -1, errors.New(ERR_LENGTH_NOT_ENOUGH)
	}
	return vlen, nil
}

func (d *decoder) transformByTag(t reflect.StructTag, buf []byte) (string, error){
	vLenStr, ok := t.Lookup("type")
	if ok{
		if vLenStr == "bcd"{
			data := BCD2DEC(buf)
			return data, nil
		}

		if vLenStr == "time"{
			data := BCD2Time(buf)
			return data, nil
		}
	}

	data, err := GbkToUtf8(buf)
	if err != nil{
		return "", err
	}
	return string(data), nil
}

func (d *decoder) value1(v reflect.Value)  error {
	var err error
	switch v.Kind() {
	case reflect.Ptr:
		err = d.value1(v.Elem())
		if err != nil{
			return err
		}
	case reflect.Array:
		l := v.Len()
		for i := 0; i < l; i++ {
			err = d.value1(v.Index(i))
			if err != nil{
				return err
			}
		}

	case reflect.Struct:
		t := v.Type()
		l := v.NumField()

		for i := 0; i < l; i++ {
			// Note: Calling v.CanSet() below is an optimization.
			// It would be sufficient to check the field name,
			// but creating the StructField info for each field is
			// costly (run "go test -bench=ReadStruct" and compare
			// results when making changes to this code).
			v1 := v.Field(i)
			if !v1.CanSet() || strings.HasPrefix(t.Field(i).Name, "_"){
				continue
			}

			switch v1.Kind() {
			case reflect.Ptr:
				if v1.IsNil(){
					v.Field(i).Set(reflect.New(v.Field(i).Type().Elem()))
				}
			case reflect.Slice:
				if v1.IsNil(){
					vlen, err := d.getLenByTag(v, t.Field(i).Tag)
					if err != nil{
						return err
					}

					arr := reflect.MakeSlice(v.Field(i).Type(), vlen, vlen)
					v.Field(i).Set(arr)
					for j := 0; j < vlen; j++ {
						if v1.Index(j).Kind() == reflect.Ptr{
							if v1.Index(j).IsNil(){
								v1.Index(j).Set(reflect.New(v1.Index(j).Type().Elem()))
							}
						}
					}
				}
			}

			if v1.Kind() == reflect.String{
				vlen, err := d.getLenByTag(v, t.Field(i).Tag)
				if err != nil{
					return err
				}

				tmp, err := d.transformByTag(t.Field(i).Tag, d.buf[d.offset:d.offset+vlen])
				if err != nil{
					fmt.Println("transfer data to string failed", t.Field(i).Name)
					return err
				}
				d.offset += vlen
				v1.SetString(tmp)
			} else {
				err = d.value1(v1)
				if err != nil{
					return err
				}
			}
		}

	case reflect.Slice:
		l := v.Len()
		for i := 0; i < l; i++ {
			err = d.value1(v.Index(i))
			if err != nil{
				return err
			}
		}

	case reflect.Bool:
		if d.offset + 1 > d.end || d.offset >= d.end {
			return errors.New(ERR_LENGTH_NOT_ENOUGH)
		}
		v.SetBool(d.bool())

	case reflect.Int8:
		if d.offset + 1 > d.end || d.offset >= d.end {
			return errors.New(ERR_LENGTH_NOT_ENOUGH)
		}
		v.SetInt(int64(d.int8()))
	case reflect.Int16:
		if d.offset + 2 > d.end || d.offset >= d.end {
			return errors.New(ERR_LENGTH_NOT_ENOUGH)
		}
		v.SetInt(int64(d.int16()))
	case reflect.Int32:
		if d.offset + 4 > d.end || d.offset >= d.end {
			return errors.New(ERR_LENGTH_NOT_ENOUGH)
		}
		v.SetInt(int64(d.int32()))
	case reflect.Int64:
		if d.offset + 8 > d.end || d.offset >= d.end {
			return errors.New(ERR_LENGTH_NOT_ENOUGH)
		}
		v.SetInt(d.int64())

	case reflect.Uint8:
		if d.offset + 1 > d.end || d.offset >= d.end {
			return errors.New(ERR_LENGTH_NOT_ENOUGH)
		}
		v.SetUint(uint64(d.uint8()))
	case reflect.Uint16:
		if d.offset + 2 > d.end || d.offset >= d.end {
			return errors.New(ERR_LENGTH_NOT_ENOUGH)
		}
		v.SetUint(uint64(d.uint16()))
	case reflect.Uint32:
		if d.offset + 4 > d.end || d.offset >= d.end {
			return errors.New(ERR_LENGTH_NOT_ENOUGH)
		}
		v.SetUint(uint64(d.uint32()))
	case reflect.Uint64:
		if d.offset + 8 > d.end || d.offset >= d.end {
			return errors.New(ERR_LENGTH_NOT_ENOUGH)
		}
		v.SetUint(d.uint64())

	case reflect.Float32:
		if d.offset + 4 > d.end || d.offset >= d.end {
			return errors.New(ERR_LENGTH_NOT_ENOUGH)
		}
		v.SetFloat(float64(math.Float32frombits(d.uint32())))
	case reflect.Float64:
		if d.offset + 8 > d.end || d.offset >= d.end {
			return errors.New(ERR_LENGTH_NOT_ENOUGH)
		}
		v.SetFloat(math.Float64frombits(d.uint64()))

	case reflect.Complex64:
		if d.offset + 8 > d.end || d.offset >= d.end {
			return errors.New(ERR_LENGTH_NOT_ENOUGH)
		}
		v.SetComplex(complex(
			float64(math.Float32frombits(d.uint32())),
			float64(math.Float32frombits(d.uint32())),
		))
	case reflect.Complex128:
		if d.offset + 16 > d.end || d.offset >= d.end {
			return errors.New(ERR_LENGTH_NOT_ENOUGH)
		}
		v.SetComplex(complex(
			math.Float64frombits(d.uint64()),
			math.Float64frombits(d.uint64()),
		))
	}
	return nil
}