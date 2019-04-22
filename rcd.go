// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	ed "github.com/FactomProject/ed25519"
)

// RCD is a Redeem Condition Datastructure representing a Factoid account. The
// RCD Hash is the address of the account. Different RCD types may conform to
// this interface and be used as part of the Factoid Transactions.
type RCD interface {
	Type() byte
	Hash() []byte
}

// RCD1 is a Type 1 Redeem Condition Datastructure which contains a public key
// used to sign transactions made with a Factoid Address.
type RCD1 struct {
	Pub *[ed.PublicKeySize]byte
}

// NewRCD1 creates a new 0 value Type 1 Factoid RCD.
func NewRCD1() *RCD1 {
	r := new(RCD1)
	r.Pub = new([ed.PublicKeySize]byte)
	return r
}

func (r *RCD1) Type() uint8 {
	return byte(1)
}

// Hash of the Type 1 RCD is the double sha256 hash of the type byte (1) and the
// RCD public key.
func (r *RCD1) Hash() []byte {
	p := append([]byte{r.Type()}, r.Pub[:]...)
	return shad(p)
}

// PubBytes may be used to validate a signature from a Type 1 Factoid RCD.
func (r *RCD1) PubBytes() []byte {
	return r.Pub[:]
}
