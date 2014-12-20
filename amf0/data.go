// Package amf0 represents AMF0 values.
package amf0

import "io"

// Type represents the type marker of AMF0 data.
type Type byte

// Type constants represent every AMF0 supported type.
const (
	TypeNumber        = Type(0x00)
	TypeBoolean       = Type(0x01)
	TypeString        = Type(0x02)
	TypeObject        = Type(0x03)
	TypeMovieClip     = Type(0x04)
	TypeNull          = Type(0x05)
	TypeUndefined     = Type(0x06)
	TypeReference     = Type(0x07)
	TypeECMAArray     = Type(0x08)
	TypeObjectEnd     = Type(0x09)
	TypeStrictArray   = Type(0x0A)
	TypeDate          = Type(0x0B)
	TypeLongString    = Type(0x0C)
	TypeUnsupported   = Type(0x0D)
	TypeRecordSet     = Type(0x0E)
	TypeXMLDocument   = Type(0x0F)
	TypeTypedObject   = Type(0x10)
	TypeAMVPlusObject = Type(0x11)
)

// Data is the base interface representing AMF data.
type Data interface {
	Type() Type
	Encode(w io.Writer) error
}

// Number represents an AMF0 number.
type Number float64

func (Number) Type() Type {
	return TypeNumber
}

// Boolean represents an AMF0 boolean value.
type Boolean bool

func (Boolean) Type() Type {
	return TypeBoolean
}

// String represents an AMF0 UTF-8 string.
type String string

func (String) Type() Type {
	return TypeString
}

// Object represents an anonymous AMF0 object.
type Object map[string]Data

func (Object) Type() Type {
	return TypeObject
}

// Null represents the AMF0 null value.
type Null struct{}

func (Null) Type() Type {
	return TypeNull
}

// Undefined represents an AMF0 undefined value.
type Undefined struct{}

func (Undefined) Type() Type {
	return TypeUndefined
}

// ECMAArray represents an associative array, almost similar to an object.
type ECMAArray map[string]Data

func (ECMAArray) Type() Type {
	return TypeECMAArray
}

// StrictArray represents an array with ordinal indices.
type StrictArray []Data

func (StrictArray) Type() Type {
	return TypeStrictArray
}
