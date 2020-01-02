package nbt

import (
	"bytes"
	"errors"
	"math"
	"reflect"
)

func Marshal(v interface{}, name string) ([]byte, error) {
	e := &encoder{}
	return e.marshal(v, name)
}

type MarshalTypeError struct {
	Kind reflect.Kind
}

func (e *MarshalTypeError) Error() string {
	return "nbt: cannot marshal " + e.Kind.String()
}

type encoder struct {
	bytes.Buffer
}

func (e *encoder) marshal(v interface{}, name string) ([]byte, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil, errors.New("nbt: cannot marshal nil")
		}
		rv = rv.Elem()
	}

	err := e.writeValue(rv, name)
	if err != nil {
		return nil, err
	}

	return e.Bytes(), nil
}

func (e *encoder) writeValue(rv reflect.Value, name string) (err error) {
	switch k := rv.Kind(); k {
	case reflect.Int8:
		e.WriteByte(TagByte)
		e.writeString(name)
		err = e.WriteByte(byte(rv.Int()))
	case reflect.Int16:
		e.WriteByte(TagShort)
		e.writeString(name)
		err = e.writeInt16(int16(rv.Int()))
	case reflect.Int32, reflect.Int:
		e.WriteByte(TagInt)
		e.writeString(name)
		err = e.writeInt32(int32(rv.Int()))
	case reflect.Int64:
		e.WriteByte(TagLong)
		e.writeString(name)
		err = e.writeInt64(rv.Int())
	case reflect.Float32:
		e.WriteByte(TagFloat)
		e.writeString(name)
		err = e.writeFloat(float32(rv.Float()))
	case reflect.Float64:
		e.WriteByte(TagDouble)
		e.writeString(name)
		err = e.writeDouble(rv.Float())
	case reflect.Array:
		switch ek := rv.Type().Elem().Kind(); ek {
		case reflect.Int8:
			e.WriteByte(TagByteArray)
		case reflect.Int32, reflect.Int:
			e.WriteByte(TagIntArray)
		case reflect.Int64:
			e.WriteByte(TagLongArray)
		default:
			err = errors.New("nbt: cannot marshal array of " + ek.String())
		}
		e.writeString(name)
		err = e.writeArray(rv)
	case reflect.String:
		e.WriteByte(TagString)
		e.writeString(name)
		err = e.writeString(rv.String())
	case reflect.Slice:
		e.WriteByte(TagList)
		e.writeString(name)
		err = e.writeSlice(rv)
	case reflect.Struct:
		e.WriteByte(TagCompound)
		e.writeString(name)
		err = e.writeCompound(rv)
	default:
		err = &MarshalTypeError{k}
	}

	return
}

func (e *encoder) writeInt16(n int16) error {
	e.WriteByte(byte(n >> 8))
	e.WriteByte(byte(n))
	return nil
}

func (e *encoder) writeInt32(n int32) error {
	e.WriteByte(byte(n >> 24))
	e.WriteByte(byte(n >> 16))
	e.WriteByte(byte(n >> 8))
	e.WriteByte(byte(n))
	return nil
}

func (e *encoder) writeInt64(n int64) error {
	e.WriteByte(byte(n >> 56))
	e.WriteByte(byte(n >> 48))
	e.WriteByte(byte(n >> 40))
	e.WriteByte(byte(n >> 32))
	e.WriteByte(byte(n >> 24))
	e.WriteByte(byte(n >> 16))
	e.WriteByte(byte(n >> 8))
	e.WriteByte(byte(n))
	return nil
}

func (e *encoder) writeFloat(n float32) error {
	return e.writeInt32(int32(math.Float32bits(n)))
}

func (e *encoder) writeDouble(n float64) error {
	return e.writeInt64(int64(math.Float64bits(n)))
}

func (e *encoder) writeString(name string) error {
	n := len(name)
	if n > 0xffff {
		return errors.New("string too long")
	}

	e.writeInt16(int16(n))
	_, err := e.WriteString(name)

	return err
}

func (e *encoder) writeArray(rv reflect.Value) error {
	e.writeInt32(int32(rv.Len()))

	var write func(v reflect.Value) error
	switch ek := rv.Type().Elem().Kind(); ek {
	case reflect.Int8:
		write = func(v reflect.Value) error {
			e.WriteByte(byte(v.Int()))
			return nil
		}
	case reflect.Int32, reflect.Int:
		write = func(v reflect.Value) error {
			e.writeInt32(int32(v.Int()))
			return nil
		}
	case reflect.Int64:
		write = func(v reflect.Value) error {
			e.writeInt64(v.Int())
			return nil
		}
	default:
		return errors.New("nbt: cannot marshal array of " + ek.String())
	}

	for i := 0; i < rv.Len(); i++ {
		err := write(rv.Index(i))
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *encoder) writeSlice(rv reflect.Value) error {
	var write func(v reflect.Value) error
	switch ek := rv.Type().Elem().Kind(); ek {
	case reflect.Int8:
		e.WriteByte(TagByte)
		write = func(v reflect.Value) error {
			return e.WriteByte(byte(v.Int()))
		}
	case reflect.Int16:
		e.WriteByte(TagShort)
		write = func(v reflect.Value) error {
			return e.writeInt16(int16(v.Int()))
		}
	case reflect.Int32:
		e.WriteByte(TagInt)
		write = func(v reflect.Value) error {
			return e.writeInt32(int32(v.Int()))
		}
	case reflect.Int64:
		e.WriteByte(TagLong)
		write = func(v reflect.Value) error {
			return e.writeInt64(v.Int())
		}
	case reflect.Float32:
		e.WriteByte(TagFloat)
		write = func(v reflect.Value) error {
			return e.writeFloat(float32(v.Float()))
		}
	case reflect.Float64:
		e.WriteByte(TagDouble)
		write = func(v reflect.Value) error {
			return e.writeDouble(v.Float())
		}
	case reflect.String:
		e.WriteByte(TagString)
		write = func(v reflect.Value) error {
			return e.writeString(v.String())
		}
	case reflect.Slice:
		e.WriteByte(TagList)
		write = func(v reflect.Value) error {
			return e.writeSlice(v)
		}
	case reflect.Struct, reflect.Map:
		e.WriteByte(TagCompound)
		write = func(v reflect.Value) error {
			return e.writeCompound(v)
		}
	default:
		return errors.New("nbt: cannot marshal slice with " + ek.String())
	}

	e.writeInt32(int32(rv.Len()))

	for i := 0; i < rv.Len(); i++ {
		err := write(rv.Index(i))
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *encoder) writeCompound(rv reflect.Value) error {
	nameIndexMap, err := getTargetFieldNames(rv.Type())
	if err != nil {
		return err
	}

	for k, v := range nameIndexMap {
		err = e.writeValue(rv.Field(v), k)
		if err != nil {
			return err
		}
	}

	return e.WriteByte(TagEnd)
}
