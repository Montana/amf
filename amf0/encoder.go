package amf0

import (
	"encoding/binary"
	"io"
	"math"
)

func (t Type) encode(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, t)
}

func encodeUTF8(s string, w io.Writer) error {
	var length uint16

	if len(s) >= math.MaxUint16 {
		length = math.MaxUint16
	} else {
		length = uint16(len(s))
	}

	if err := binary.Write(w, binary.BigEndian, length); err != nil {
		return err
	}

	if _, err := w.Write([]byte(s)[0:length]); err != nil {
		return err
	}

	return nil
}

func (p Number) Encode(w io.Writer) error {
	if err := p.Type().encode(w); err != nil {
		return err
	}
	return binary.Write(w, binary.BigEndian, float64(p))
}

var (
	booleanFalse []byte = []byte{0x00}
	booleanTrue  []byte = []byte{0x01}
)

func (p Boolean) Encode(w io.Writer) error {
	if err := p.Type().encode(w); err != nil {
		return err
	}

	var val []byte
	if p {
		val = booleanTrue
	} else {
		val = booleanFalse
	}
	if _, err := w.Write(val); err != nil {
		return err
	}

	return nil
}

func (p String) Encode(w io.Writer) error {
	l := len(p)
	// depending on the string length, encode as string or long string
	if l <= math.MaxUint16 {
		if err := p.Type().encode(w); err != nil {
			return err
		}

		length := uint16(l)
		if err := binary.Write(w, binary.BigEndian, length); err != nil {
			return err
		}

	} else {
		if err := TypeLongString.encode(w); err != nil {
			return err
		}

		length := uint32(l)
		if err := binary.Write(w, binary.BigEndian, length); err != nil {
			return err
		}
	}

	if _, err := w.Write([]byte(p)); err != nil {
		return err
	}

	return nil
}

func (p Object) Encode(w io.Writer) error {
	if err := p.Type().encode(w); err != nil {
		return err
	}

	for name, data := range p {
		// encode entry name
		if err := encodeUTF8(name, w); err != nil {
			return err
		}

		// encode entry value
		if err := data.Encode(w); err != nil {
			return err
		}
	}

	// empty name is the last element
	if err := encodeUTF8("", w); err != nil {
		return err
	}

	// finally object end marker
	return TypeObjectEnd.encode(w)
}

func (p Null) Encode(w io.Writer) error {
	return p.Type().encode(w)
}

func (p Undefined) Encode(w io.Writer) error {
	return p.Type().encode(w)
}

func (p ECMAArray) Encode(w io.Writer) error {
	if err := p.Type().encode(w); err != nil {
		return err
	}

	// encode array length
	length := uint32(len(p))
	if err := binary.Write(w, binary.BigEndian, length); err != nil {
		return err
	}

	for name, data := range p {
		// encode entry name
		if err := encodeUTF8(name, w); err != nil {
			return err
		}

		// encode entry value
		if err := data.Encode(w); err != nil {
			return err
		}
	}

	return nil
}

func (p StrictArray) Encode(w io.Writer) error {
	if err := p.Type().encode(w); err != nil {
		return err
	}
	
	// encode array length
	length := uint32(len(p))
	if err := binary.Write(w, binary.BigEndian, length); err != nil {
		return err
	}
	
	for _, data := range p {
		// encode entry value
		if err := data.Encode(w); err != nil {
			return err
		}
	}

	return nil
}
