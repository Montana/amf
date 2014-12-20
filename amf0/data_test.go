package amf0

import (
	"bytes"
	"strings"
	"testing"
)

func TestNumber(t *testing.T) {
	v := Number(5.0)

	if v.Type() != TypeNumber {
		t.Error("Invalid type")
	}
}

func TestString(t *testing.T) {
	v := String("foobar")

	if v.Type() != TypeString {
		t.Error("Invalid type")
	}
}

func TestBoolean(t *testing.T) {
	v := Boolean(false)

	if v.Type() != TypeBoolean {
		t.Error("Invalid type")
	}
}

func TestObject(t *testing.T) {
	v := make(Object)

	if v.Type() != TypeObject {
		t.Error("Invalid type")
	}
}

func TestNull(t *testing.T) {
	v := Null{}

	if v.Type() != TypeNull {
		t.Error("Invalid type")
	}
}

func TestUndefined(t *testing.T) {
	v := Undefined{}

	if v.Type() != TypeUndefined {
		t.Error("Invalid type")
	}
}

func TestECMAArray(t *testing.T) {
	v := make(ECMAArray)

	if v.Type() != TypeECMAArray {
		t.Error("Invalid type")
	}
}

func TestStrictArray(t *testing.T) {
	v := make(StrictArray, 5)
	if v.Type() != TypeStrictArray {
		t.Error("Invalid type")
	}

	v = StrictArray{Boolean(true), Number(21.0), String("erf1")}
	if len(v) != 3 {
		t.Error("Invalid size")
	}
	if v[0] != Boolean(true) {
		t.Error("Invalid value")
	}
	if v[1] != Number(21.0) {
		t.Error("Invalid value")
	}
	if v[2] != String("erf1") {
		t.Error("Invalid value")
	}
}

func TestLongString(t *testing.T) {
	s := String(strings.Repeat("abc", 30000))
	if len(s) != 90000 {
		t.Fatal("Invalid string length")
	}

	b := new(bytes.Buffer)
	if err := s.Encode(b); err != nil {
		t.Fatal(err)
	}

	expected := []byte{
		0x0c,                   // long string type marker
		0x00, 0x01, 0x5f, 0x90, // length = 90000, big endian, 32-bit
	}
	if !bytes.Equal(b.Bytes()[0:5], expected) {
		t.Error("Expected ", expected, " got ", b.Bytes()[0:5])
	}
}
