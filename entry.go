package factom

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
)

type Entry struct {
	ChainID string
	ExtIDs  []string
	Content string
}

func (e *Entry) Hash() []byte {
	a, err := e.MarshalBinary()
	if err != nil {
		return make([]byte, 32)
	}
	return sha23(a)
}

func (e *Entry) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	c, err := hex.DecodeString(e.Content)
	if err != nil {
		return buf.Bytes(), err
	}
	x, err := e.MarshalExtIDsBinary()
	if err != nil {
		return buf.Bytes(), err
	}

	// Header

	// 1 byte Version
	buf.Write([]byte{0})

	// 32 byte chainid
	if p, err := hex.DecodeString(e.ChainID); err != nil {
		return buf.Bytes(), err
	} else {
		buf.Write(p)
	}

	// 2 byte size of extids
	if err := binary.Write(buf, binary.BigEndian, int16(len(x))); err != nil {
		return buf.Bytes(), err
	}

	// Payload

	// extids
	buf.Write(x)

	// data
	buf.Write(c)

	return buf.Bytes(), nil
}

func (e *Entry) MarshalExtIDsBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	for _, v := range e.ExtIDs {
		p, err := hex.DecodeString(v)
		if err != nil {
			return buf.Bytes(), err
		}
		// 2 byte length of extid
		binary.Write(buf, binary.BigEndian, int16(len(p)))
		// extid
		buf.Write(p)
	}

	return buf.Bytes(), nil
}
