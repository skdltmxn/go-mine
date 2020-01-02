package nbt

import (
	"encoding/hex"
	"testing"
)

func TestMarshalFloat(t *testing.T) {
	data := float32(2.71828)
	raw, err := Marshal(data, "float")
	if err != nil {
		t.Fatalf("Marshal failed: %s", err)
	}

	t.Logf("Marshalled data: %s", hex.EncodeToString(raw))
}

func TestMarshalDouble(t *testing.T) {
	data := float64(3.141592)
	raw, err := Marshal(data, "double")
	if err != nil {
		t.Fatalf("Marshal failed: %s", err)
	}

	t.Logf("Marshalled data: %s", hex.EncodeToString(raw))
}

func TestMarshalString(t *testing.T) {
	data := "minecraft:brand=go-mine"
	raw, err := Marshal(data, "string")
	if err != nil {
		t.Fatalf("Marshal failed: %s", err)
	}

	t.Logf("Marshalled data: %s", hex.EncodeToString(raw))
}

func TestMarshalCompount(t *testing.T) {
	type test_struct struct {
		A string `nbt:"a"`
		B int    `nbt:"b"`
	}

	ts := &test_struct{"test", 1234}
	raw, err := Marshal(ts, "")
	if err != nil {
		t.Fatalf("Marshal failed: %s", err)
	}

	t.Logf("Marshalled data: %s", hex.EncodeToString(raw))
}

func TestMarshalListOfCompounds(t *testing.T) {
	type Inventory struct {
		Slot  int8
		Id    string `nbt:"id"`
		Count int8
	}
	data := []Inventory{
		Inventory{0, "minecraft:gold_ore", 1},
		Inventory{1, "minecraft:stone", 1},
	}

	raw, err := Marshal(data, "Inventory")
	if err != nil {
		t.Fatalf("Marshal failed: %s", err)
	}

	t.Logf("Marshalled data: %s", hex.EncodeToString(raw))
}

func TestMarshalByteArray(t *testing.T) {
	data := [5]int8{1, 2, 3, 4, 5}
	raw, err := Marshal(data, "ByteArray")
	if err != nil {
		t.Fatalf("Marshal failed: %s", err)
	}

	t.Logf("Marshalled data: %s", hex.EncodeToString(raw))
}

func TestMarshalIntArray(t *testing.T) {
	data := [8]int{1, 2, 3, 4, 5, 6, 7, 8}
	raw, err := Marshal(data, "IntArray")
	if err != nil {
		t.Fatalf("Marshal failed: %s", err)
	}

	t.Logf("Marshalled data: %s", hex.EncodeToString(raw))
}

func TestMarshalLongArray(t *testing.T) {
	data := [4]int64{1, 2, 3, 4}
	raw, err := Marshal(data, "LongArray")
	if err != nil {
		t.Fatalf("Marshal failed: %s", err)
	}

	t.Logf("Marshalled data: %s", hex.EncodeToString(raw))
}
