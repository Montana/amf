package amf0

import (
	"bytes"
	"math"
	"testing"
)

type dataEncodeTest struct {
	value    Data
	expected []byte
}

var dataEncodeTests = []dataEncodeTest{
	// Number: header, followed by big-endian IEEE-754 representation
	{Number(0.0), []byte{
		0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}},
	{Number(1.0), []byte{
		0x00,
		0x3f, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}},
	{Number(2.0), []byte{
		0x00,
		0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}},
	{Number(-2.0), []byte{
		0x00,
		0xc0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}},
	{Number(1.0000000000000002), []byte{
		0x00,
		0x3f, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
	}},
	{Number(3.4e12), []byte{
		0x00,
		0x42, 0x88, 0xbc, 0xfe, 0x56, 0x80, 0x00, 0x00,
	}},
	{Number(-100000.0), []byte{
		0x00,
		0xC0, 0xF8, 0x6a, 0x00, 0x00, 0x00, 0x00, 0x00,
	}},
	{Number(math.Pi), []byte{
		0x00,
		0x40, 0x09, 0x21, 0xfb, 0x54, 0x44, 0x2d, 0x18,
	}},

	// Boolean
	{Boolean(true), []byte{0x01, 0x01}},
	{Boolean(false), []byte{0x01, 0x00}},

	// String: header followed by big endian unsigned 16-bit length, and UTF-8 encoded string
	{String(""), []byte{
		0x02,
		0x00, 0x00,
	}},
	{String("foo"), []byte{
		0x02,
		0x00, 0x03,
		0x66, 0x6f, 0x6f,
	}},
	{String("FOO BAR"), []byte{
		0x02,
		0x00, 0x07,
		0x46, 0x4f, 0x4f, 0x20, 0x42, 0x41, 0x52,
	}},
	{String("580€"), []byte{
		0x02,
		0x00, 0x06,
		0x35, 0x38, 0x30, 0xe2, 0x82, 0xac,
	}},

	// Object
	{make(Object), []byte{
		0x03,       // object type marker
		0x00, 0x00, // empty utf-8 string (only 16-bit length field)
		0x09, // object end marker
	}},
	{Object{"test": Boolean(true)}, []byte{
		0x03,       // marker
		0x00, 0x04, // name length
		0x74, 0x65, 0x73, 0x74, // name 'test'
		0x01,       // boolean marker
		0x01,       // true
		0x00, 0x00, // empty string
		0x09, // object end marker
	}},
	{Object{"test": make(Object)}, []byte{
		0x03,       // marker
		0x00, 0x04, // name length
		0x74, 0x65, 0x73, 0x74, // name 'test'
		0x03,       // object marker
		0x00, 0x00, // empty string
		0x09,       // object end marker
		0x00, 0x00, // empty string
		0x09, // object end marker
	}},

	// Null
	{Null{}, []byte{0x05}},

	// Undefined
	{Undefined{}, []byte{0x06}},

	// ECMA Array
	{make(ECMAArray), []byte{
		0x08,                   // object type marker
		0x00, 0x00, 0x00, 0x00, // big endian 32-bit length
	}},
	{ECMAArray{"test": Boolean(true)}, []byte{
		0x08,                   // marker
		0x00, 0x00, 0x00, 0x01, // length
		0x00, 0x04, // name length
		0x74, 0x65, 0x73, 0x74, // name 'test'
		0x01, // boolean marker
		0x01, // true
	}},
	{ECMAArray{"test": make(ECMAArray)}, []byte{
		0x08,                   // marker
		0x00, 0x00, 0x00, 0x01, // length
		0x00, 0x04, // name length
		0x74, 0x65, 0x73, 0x74, // name 'test'
		0x08,                   // object marker
		0x00, 0x00, 0x00, 0x00, // length
	}},

	// Strict Array
	{StrictArray{}, []byte{
		0x0a,                   // marker
		0x00, 0x00, 0x00, 0x00, // empty array length
	}},
	{StrictArray{Number(1.0)}, []byte{
		0x0a,                   // marker
		0x00, 0x00, 0x00, 0x01, // empty array length
		0x00,                                           // number marker
		0x3f, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // value
	}},
	{StrictArray{Boolean(true), Boolean(false)}, []byte{
		0x0a,                   // marker
		0x00, 0x00, 0x00, 0x02, // empty array length
		0x01, 0x01, // true
		0x01, 0x00, // false
	}},
}

func TestDataEncode(t *testing.T) {
	for _, entry := range dataEncodeTests {
		b := new(bytes.Buffer)
		if err := entry.value.Encode(b); err != nil {
			t.Error(err)
			continue
		}

		if !bytes.Equal(entry.expected, b.Bytes()) {
			t.Error("Expected ", entry.expected, " got ", b.Bytes())
		}
	}
}

func TestObjectEncode(t *testing.T) {
	o := Object{"test": Boolean(true), "test2": Boolean(false)}

	// handle random map interation
	expected1 := []byte{
		0x03,       // marker
		0x00, 0x04, // name length
		0x74, 0x65, 0x73, 0x74, // name 'test'
		0x01,       // boolean marker
		0x01,       // true
		0x00, 0x05, // name length
		0x74, 0x65, 0x73, 0x74, 0x32, // name 'test2'
		0x01,       // boolean marker
		0x00,       // false
		0x00, 0x00, // empty string
		0x09, // object end marker
	}
	expected2 := []byte{
		0x03,       // marker
		0x00, 0x05, // name length
		0x74, 0x65, 0x73, 0x74, 0x32, // name 'test2'
		0x01,       // boolean marker
		0x00,       // false
		0x00, 0x04, // name length
		0x74, 0x65, 0x73, 0x74, // name 'test'
		0x01,       // boolean marker
		0x01,       // true
		0x00, 0x00, // empty string
		0x09, // object end marker
	}

	b := new(bytes.Buffer)
	if err := o.Encode(b); err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(b.Bytes(), expected1) && !bytes.Equal(b.Bytes(), expected2) {
		t.Error("Invalid output ", b.Bytes())
	}
}

func TestECMAArrayEncode(t *testing.T) {
	a := ECMAArray{"test": Boolean(true), "test2": Boolean(false)}

	// handle random map interation
	expected1 := []byte{
		0x08,                   // marker
		0x00, 0x00, 0x00, 0x02, // length
		0x00, 0x04, // name length
		0x74, 0x65, 0x73, 0x74, // name 'test'
		0x01,       // boolean marker
		0x01,       // true
		0x00, 0x05, // name length
		0x74, 0x65, 0x73, 0x74, 0x32, // name 'test2'
		0x01, // boolean marker
		0x00, // false
	}
	expected2 := []byte{
		0x08,       // marker
		0x00, 0x05, // name length
		0x74, 0x65, 0x73, 0x74, 0x32, // name 'test2'
		0x01,                   // boolean marker
		0x00,                   // false
		0x00, 0x00, 0x00, 0x02, // length
		0x00, 0x04, // name length
		0x74, 0x65, 0x73, 0x74, // name 'test'
		0x01, // boolean marker
		0x01, // true
	}

	b := new(bytes.Buffer)
	if err := a.Encode(b); err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(b.Bytes(), expected1) && !bytes.Equal(b.Bytes(), expected2) {
		t.Error("Invalid output ", b.Bytes())
	}
}
