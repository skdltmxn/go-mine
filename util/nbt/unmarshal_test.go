package nbt

import (
	"encoding/hex"
	"testing"
)

func TestNbtUnmarshalLongArray(t *testing.T) {
	data, _ := hex.DecodeString("0C000000000003000000000000010100000000000002210000000000003322")
	var val []int32
	err := Unmarshal(data, &val)
	if err != nil {
		t.Errorf("Unmarshal failed: %s", err)
		return
	}

	t.Logf("Unmarshalled data: %+v", val)
}

func TestNbtUnmarshalIntArray(t *testing.T) {
	data, _ := hex.DecodeString("0B00000000000400000001000000020000000300000004")
	var val []int64
	err := Unmarshal(data, &val)
	if err != nil {
		t.Errorf("Unmarshal failed: %s", err)
		return
	}

	t.Logf("Unmarshalled data: %+v", val)
}

func TestNbtUnmarshalCompound(t *testing.T) {
	data, _ := hex.DecodeString("0A000568656C6C6F08000161000474657374030001620000000106000163BFB41205C28F5C2904000164000000000001234500")
	type dummy struct {
		Abc string `nbt:"a"`
		B   int32  `nbt:"b"`
		D   int64  `nbt:"d"`
	}
	var val dummy
	err := Unmarshal(data, &val)
	if err != nil {
		t.Errorf("Unmarshal failed: %s", err)
		return
	}

	t.Logf("Unmarshalled data: %+v", val)
}

func TestNbtUnmarshalList(t *testing.T) {
	data, _ := hex.DecodeString("090003506F7306000000034049CCF44787744040518000000000004028E09D1CEB0051")
	var val [2]float64
	err := Unmarshal(data, &val)
	if err != nil {
		t.Errorf("Unmarshal failed: %s", err)
		return
	}

	t.Logf("Unmarshalled data: %+v", val)
}

func TestNbtUnmarshalString(t *testing.T) {
	data, _ := hex.DecodeString("08000D67656E657261746F724E616D65000764656661756C74")
	var val string
	err := Unmarshal(data, &val)
	if err != nil {
		t.Errorf("Unmarshal failed: %s", err)
		return
	}

	t.Logf("Unmarshalled data: %s", val)
}

func TestNbtUnmarshalByteArray(t *testing.T) {
	data, _ := hex.DecodeString("070000000000050102030405")
	val := make([]byte, 5)
	err := Unmarshal(data, &val)
	if err != nil {
		t.Errorf("Unmarshal failed: %s", err)
		return
	}

	t.Logf("Unmarshalled data: %+v", val)
}

func TestNbtUnmarshalDouble(t *testing.T) {
	data, _ := hex.DecodeString("060000BFB41205C28F5C29")
	var val float64
	err := Unmarshal(data, &val)
	if err != nil {
		t.Errorf("Unmarshal failed: %s", err)
		return
	}

	t.Logf("Unmarshalled data: %f", val)
}

func TestNbtUnmarshalFloat(t *testing.T) {
	data, _ := hex.DecodeString("050000C2ABDE84")
	var val float32
	err := Unmarshal(data, &val)
	if err != nil {
		t.Errorf("Unmarshal failed: %s", err)
		return
	}

	t.Logf("Unmarshalled data: %f", val)
}

func TestNbtUnmarshalInt64(t *testing.T) {
	data, _ := hex.DecodeString("04000A52616E646F6D53656564000000000001824F")
	var val int64
	err := Unmarshal(data, &val)
	if err != nil {
		t.Errorf("Unmarshal failed: %s", err)
		return
	}

	t.Logf("Unmarshalled data: %d", val)
}

func TestNbtUnmarshalInt32(t *testing.T) {
	data, _ := hex.DecodeString("03001A57616E646572696E67547261646572537061776E4368616E636500000019")
	var val int
	err := Unmarshal(data, &val)
	if err != nil {
		t.Errorf("Unmarshal failed: %s", err)
		return
	}

	t.Logf("Unmarshalled data: %d", val)
}

func TestNbtUnmarshalInt16(t *testing.T) {
	data, _ := hex.DecodeString("02001A57616E646572696E67547261646572537061776E4368616E63651234")
	var val int
	err := Unmarshal(data, &val)
	if err != nil {
		t.Errorf("Unmarshal failed: %s", err)
		return
	}

	t.Logf("Unmarshalled data: %d", val)
}

func TestNbtUnmarshalByte(t *testing.T) {
	data, _ := hex.DecodeString("01001A57616E646572696E67547261646572537061776E4368616E6365FF")
	var val int
	err := Unmarshal(data, &val)
	if err != nil {
		t.Errorf("Unmarshal failed: %s", err)
		return
	}

	t.Logf("Unmarshalled data: %d", val)
}
