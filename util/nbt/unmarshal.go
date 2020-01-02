package nbt

import (
	"bytes"
	"errors"
	"io"
	"log"
	"math"
	"reflect"
	"strconv"
)

func Unmarshal(data []byte, v interface{}) error {
	d := &decoder{bytes.NewReader(data)}
	return d.unmarshal(v)
}

type UnmarshalError struct {
	Type reflect.Type
}

func (e *UnmarshalError) Error() string {
	if e.Type == nil {
		return "nbt: cannot unmarshal nil"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "nbt: cannot unmarshal non-ptr (" + e.Type.String() + ")"
	}

	return "nbt: cannot unmarshal " + e.Type.String()
}

type UnmarshalTypeError struct {
	Src string
	Dst reflect.Kind
}

func (e *UnmarshalTypeError) Error() string {
	return "nbt: cannot unmarshal from " + e.Src + " to " + e.Dst.String()
}

type decoder struct {
	r *bytes.Reader
}

func (d *decoder) unmarshal(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &UnmarshalError{reflect.TypeOf(v)}
	}

	t, name, err := d.readTag()
	if err != nil {
		log.Printf("Failed to unmarshal: %s", err)
		return err
	}

	return d.readValue(t, name, rv.Elem())
}

func (d *decoder) readTag() (tag byte, name string, err error) {
	// tag is always 1 byte
	tag, err = d.r.ReadByte()

	if tag == TagEnd || tag > TagLongArray || err != nil {
		return
	}

	// read name
	name, err = d.readString()
	if err != nil {
		return
	}

	return
}

func (d *decoder) readValue(tag byte, name string, v reflect.Value) error {
	switch tag {
	case TagByte:
		switch k := v.Kind(); k {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			value, err := d.r.ReadByte()
			if err != nil {
				return err
			}

			v.SetInt(int64(value))
		default:
			return &UnmarshalTypeError{"Byte", k}
		}
	case TagShort:
		switch k := v.Kind(); k {
		case reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			value, err := d.readInt16()
			if err != nil {
				return err
			}

			v.SetInt(int64(value))
		default:
			return &UnmarshalTypeError{"Short", k}
		}
	case TagInt:
		switch k := v.Kind(); k {
		case reflect.Int, reflect.Int32, reflect.Int64:
			value, err := d.readInt32()
			if err != nil {
				return err
			}

			v.SetInt(int64(value))
		default:
			return &UnmarshalTypeError{"Int", k}
		}
	case TagLong:
		switch k := v.Kind(); k {
		case reflect.Int64:
			value, err := d.readInt64()
			if err != nil {
				return err
			}

			v.SetInt(int64(value))
		default:
			return &UnmarshalTypeError{"Long", k}
		}
	case TagFloat:
		switch k := v.Kind(); k {
		case reflect.Float32, reflect.Float64:
			value, err := d.readInt32()
			if err != nil {
				return err
			}

			f := math.Float32frombits(uint32(value))
			v.SetFloat(float64(f))
		default:
			return &UnmarshalTypeError{"Float", k}
		}
	case TagDouble:
		switch k := v.Kind(); k {
		case reflect.Float64:
			value, err := d.readInt64()
			if err != nil {
				return err
			}

			f := math.Float64frombits(uint64(value))
			v.SetFloat(f)
		default:
			return &UnmarshalTypeError{"Double", k}
		}
	case TagByteArray:
		switch k := v.Kind(); k {
		case reflect.Slice:
			length, err := d.readInt32()
			if err != nil {
				return err
			}

			bs := make([]byte, length)
			_, err = d.r.Read(bs)
			if err != nil {
				return err
			}

			v.SetBytes(bs)
		default:
			return &UnmarshalTypeError{"ByteArray", k}
		}
	case TagString:
		switch k := v.Kind(); k {
		case reflect.String:
			s, err := d.readString()
			if err != nil {
				return err
			}

			v.SetString(s)
		default:
			return &UnmarshalTypeError{"String", k}
		}
	case TagList:
		t, err := d.r.ReadByte()
		if err != nil {
			return err
		}

		l, err := d.readInt32()
		if err != nil {
			return err
		}

		var value reflect.Value
		switch k := v.Kind(); k {
		case reflect.Slice:
			value = reflect.MakeSlice(v.Type(), int(l), int(l))
		case reflect.Array:
			if v.Len() < int(l) {
				return errors.New("nbt: given array size is smaller than payload (" + strconv.Itoa(v.Len()) + " < " + strconv.Itoa(int(l)) + ")")
			}
			value = v

		default:
			return &UnmarshalTypeError{"List", k}
		}

		for i := 0; i < int(l); i++ {
			if err := d.readValue(t, "", value.Index(i)); err != nil {
				return err
			}
		}

		v.Set(value)
	case TagCompound:
		switch k := v.Kind(); k {
		case reflect.Struct:
			nameIndexMap, err := getTargetFieldNames(v.Type())
			if err != nil {
				return err
			}

			for {
				t, name, err := d.readTag()
				if err != nil {
					return err
				}

				if t == TagEnd {
					break
				}

				if index, ok := nameIndexMap[name]; ok {
					if err = d.readValue(t, name, v.Field(index)); err != nil {
						return err
					}
				} else {
					if err = d.skip(t); err != nil {
						return err
					}
				}
			}
		default:
			return &UnmarshalTypeError{"Compound", k}
		}
	case TagIntArray:
		l, err := d.readInt32()
		if err != nil {
			return err
		}

		switch k := v.Kind(); k {
		case reflect.Slice:
			if elemKind := v.Type().Elem().Kind(); elemKind != reflect.Int && elemKind != reflect.Int32 {
				return &UnmarshalTypeError{"IntArray", elemKind}
			}

			value := reflect.MakeSlice(v.Type(), int(l), int(l))

			for i := 0; i < int(l); i++ {
				n, err := d.readInt32()
				if err != nil {
					return err
				}
				value.Index(i).SetInt(int64(n))
			}

			v.Set(value)
		default:
			return &UnmarshalTypeError{"IntArray", k}
		}
	case TagLongArray:
		l, err := d.readInt32()
		if err != nil {
			return err
		}

		switch k := v.Kind(); k {
		case reflect.Slice:
			if elemKind := v.Type().Elem().Kind(); elemKind != reflect.Int64 {
				return &UnmarshalTypeError{"LongArray", elemKind}
			}

			value := reflect.MakeSlice(v.Type(), int(l), int(l))

			for i := 0; i < int(l); i++ {
				n, err := d.readInt64()
				if err != nil {
					return err
				}
				value.Index(i).SetInt(int64(n))
			}

			v.Set(value)
		default:
			return &UnmarshalTypeError{"LongArray", k}
		}
	}

	return nil
}

func (d *decoder) skip(t byte) (err error) {
	switch t {
	case TagByte:
		_, err = d.r.ReadByte()
	case TagShort:
		_, err = d.readInt16()
	case TagInt, TagFloat:
		_, err = d.readInt32()
	case TagLong, TagDouble:
		_, err = d.readInt64()
	case TagByteArray:
		var l int32
		l, err = d.readInt32()
		if err != nil {
			return
		}
		_, err = d.r.Seek(int64(l), io.SeekCurrent)
		if err != nil {
			return
		}
	case TagString:
		_, err = d.readString()
	case TagList:
		var listType byte
		listType, err = d.r.ReadByte()
		if err != nil {
			return
		}
		var l int32
		l, err = d.readInt32()
		if err != nil {
			return
		}
		for i := 0; i < int(l); i++ {
			d.skip(listType)
		}
	case TagCompound:
		var tag byte
		for {
			tag, _, err = d.readTag()
			if tag == TagEnd {
				break
			}
			d.skip(tag)
		}
	case TagIntArray, TagLongArray:
		var delta int64 = 4
		if t == TagLongArray {
			delta = 8
		}
		var l int32
		l, err = d.readInt32()
		if err != nil {
			return
		}
		_, err = d.r.Seek(int64(l)*delta, io.SeekCurrent)
		if err != nil {
			return
		}
	}

	return
}

func (d *decoder) readInt16() (v int16, err error) {
	bs := make([]byte, 2)
	_, err = d.r.Read(bs)
	if err != nil {
		return
	}
	v = int16(bs[0])<<8 | int16(bs[1])

	return
}

func (d *decoder) readInt32() (v int32, err error) {
	bs := make([]byte, 4)
	_, err = d.r.Read(bs)
	if err != nil {
		return
	}
	v = int32(bs[0])<<24 |
		int32(bs[1])<<16 |
		int32(bs[2])<<8 |
		int32(bs[3])

	return
}

func (d *decoder) readInt64() (v int64, err error) {
	bs := make([]byte, 8)
	_, err = d.r.Read(bs)
	if err != nil {
		return
	}

	v = int64(bs[0])<<56 |
		int64(bs[1])<<48 |
		int64(bs[2])<<40 |
		int64(bs[3])<<32 |
		int64(bs[4])<<24 |
		int64(bs[5])<<16 |
		int64(bs[6])<<8 |
		int64(bs[7])

	return
}

func (d *decoder) readString() (v string, err error) {
	l, err := d.readInt16()
	if err != nil {
		return
	}

	s := make([]byte, l)
	_, err = d.r.Read(s)
	if err != nil {
		return
	}

	v = string(s)
	return
}
