// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	ed "github.com/FactomProject/ed25519"
)

type RCD interface {
	Type() byte
	Hash() []byte
}

type RCD1 struct {
	Pub *[ed.PublicKeySize]byte
}

func NewRCD1() *RCD1 {
	r := new(RCD1)
	r.Pub = new([ed.PublicKeySize]byte)
	return r
}

func (r *RCD1) Type() uint8 {
	return byte(1)
}

func (r *RCD1) Hash() []byte {
	p := append([]byte{r.Type()}, r.Pub[:]...)
	return shad(p)
}

func (r *RCD1) PubBytes() []byte {
	return r.Pub[:]
}
