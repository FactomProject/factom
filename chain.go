// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"crypto/sha256"
	"encoding/hex"
)

type Chain struct {
	ChainID    string
	FirstEntry *Entry
}

func NewChain(e *Entry) *Chain {
	c := new(Chain)
	c.FirstEntry = e

	// create the chainid from a series of hashes of the Entries ExtIDs
	hs := sha256.New()
	for _, id := range e.ExtIDs {
		h := sha256.Sum256(id)
		hs.Write(h[:])
	}
	c.ChainID = hex.EncodeToString(hs.Sum(nil))
	c.FirstEntry.ChainID = c.ChainID

	return c
}

type CHead struct {
	ChainHead string
}
