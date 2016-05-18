// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"bytes"
	
	ed "github.com/FactomProject/ed25519"
	"github.com/FactomProject/btcutil/base58"
)

type ECAddress struct {
	Pub *[ed.PublicKeySize]byte
	Sec *[ed.PrivateKeySize]byte
}

func (e *ECAddress) convertToUser() string {
	buf := new(bytes.Buffer)
	
	// EC address prefix
	buf.Write([]byte{0x59, 0x2a})
	
	// Public key
	buf.Write(e.Pub[:])
	
	// Checksum
	check := shad(buf.Bytes())[:4]
	buf.Write(check)
	
	return base58.Encode(buf.Bytes())
}