// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom_test

import (
	"bytes"
	"crypto/rand"
	"testing"

	ed "github.com/FactomProject/ed25519"
	. "github.com/FactomProject/factom"
	"github.com/FactomProject/go-bip32"
)

var ()

func TestMarshalAddresses(t *testing.T) {
	for i := 0; i < 100; i++ {
		sec := make([]byte, 32)
		_, err := rand.Read(sec)
		if err != nil {
			t.Error(err)
		}

		ec, err := MakeECAddress(sec)
		if err != nil {
			t.Error(err)
		}

		data, err := ec.MarshalBinary()
		if err != nil {
			t.Error(err)
		}

		ec2 := new(ECAddress)
		newdata, err := ec2.UnmarshalBinaryData(data)
		if err != nil {
			t.Error(err)
		}

		if len(newdata) != 0 {
			t.Errorf("UnmarshalBinary left %d bytes remaining", len(newdata))
		}

		if bytes.Compare(ec.SecBytes(), ec2.SecBytes()) != 0 {
			t.Errorf("Unmarshaled object has different secret.")
		}

		if bytes.Compare(ec.PubBytes(), ec2.PubBytes()) != 0 {
			t.Errorf("Unmarshaled object has different public.")
		}
	}

	for i := 0; i < 100; i++ {
		sec := make([]byte, 32)
		_, err := rand.Read(sec)
		if err != nil {
			t.Error(err)
		}

		fa, err := MakeFactoidAddress(sec)
		if err != nil {
			t.Error(err)
		}

		data, err := fa.MarshalBinary()
		if err != nil {
			t.Error(err)
		}

		fa2 := new(FactoidAddress)
		newdata, err := fa2.UnmarshalBinaryData(data)
		if err != nil {
			t.Error(err)
		}

		if len(newdata) != 0 {
			t.Errorf("UnmarshalBinary left %d bytes remaining", len(newdata))
		}

		if bytes.Compare(fa.SecBytes(), fa2.SecBytes()) != 0 {
			t.Errorf("Unmarshaled object has different secret.")
		}

		if bytes.Compare(fa.PubBytes(), fa2.PubBytes()) != 0 {
			t.Errorf("Unmarshaled object has different public.")
		}
	}
}

func TestAddressStringType(t *testing.T) {
	var (
		a0 = "FX1zT4aFpEvcnPqPCigB3fvGu4Q4mTXY22iiuV69DqE1pNhdF2MX"
		a1 = "FA1zT4aFpEvcnPqPCigB3fvGu4Q4mTXY22iiuV69DqE1pNhdF2MC"
		a2 = "Fs1KWJrpLdfucvmYwN2nWrwepLn8ercpMbzXshd1g8zyhKXLVLWj"
		a3 = "EC2DKSYyRcNWf7RS963VFYgMExoHRYLHVeCfQ9PGPmNzwrcmgm2r"
		a4 = "Es2Rf7iM6PdsqfYCo3D1tnAR65SkLENyWJG1deUzpRMQmbh9F3eG"
	)

	if v := AddressStringType(a0); v != InvalidAddress {
		t.Errorf("invalid address has wrong type %s %#v", a0, v)
	}
	if v := AddressStringType(a1); v != FactoidPub {
		t.Errorf("wrong address type %s %#v", a1, v)
	}
	if v := AddressStringType(a2); v != FactoidSec {
		t.Errorf("wrong address type %s %#v", a1, v)
	}
	if v := AddressStringType(a3); v != ECPub {
		t.Errorf("wrong address type %s %#v", a1, v)
	}
	if v := AddressStringType(a4); v != ECSec {
		t.Errorf("wrong address type %s %#v", a1, v)
	}
}

func TestNewECAddress(t *testing.T) {
	zPub := "EC1m9mouvUQeEidmqpUYpYtXg8fvTYi6GNHaKg8KMLbdMBrFfmUa"
	e := NewECAddress()
	if e.PubString() != zPub {
		t.Errorf("new address %s did not match %s", e.PubString(), zPub)
	}
}

func TestECAddress(t *testing.T) {
	zPub := "EC1m9mouvUQeEidmqpUYpYtXg8fvTYi6GNHaKg8KMLbdMBrFfmUa"
	zSec := "Es2Rf7iM6PdsqfYCo3D1tnAR65SkLENyWJG1deUzpRMQmbh9F3eG"
	e := NewECAddress()
	e.Pub = &[32]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	e.Sec = &[64]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01}

	if e.PubString() != zPub {
		t.Errorf("%s did not match %s", e.PubString(), zPub)
	}

	if e.SecString() != zSec {
		t.Errorf("%s did not match %s", e.SecString(), zSec)
	}
}

func TestIsValidECAddress(t *testing.T) {
	zPub := "EC1m9mouvUQeEidmqpUYpYtXg8fvTYi6GNHaKg8KMLbdMBrFfmUa"
	zSec := "Es2Rf7iM6PdsqfYCo3D1tnAR65SkLENyWJG1deUzpRMQmbh9F3eG"
	badEmpty := ""
	badLen := "EC1m9mouvUQeEidmqpUYpYtXgfvTYi6GNHaKg8KMLbdMBrFfmUa"
	badPrePub := "Ec1m9mouvUQeEidmqpUYpYtXg8fvTYi6GNHaKg8KMLbdMBrFfmUa"
	badPreSec := "ER2Rf7iM6PdsqfYCo3D1tnAR65SkLENyWJG1deUzpRMQmbh9F3eG"
	badCheckPub := "EC1m9mouvUQeEidmqpUYpYtXg8fvTYi6GNHaKg8KMLbdMBrFfgUa"
	badCheckSec := "Es2Rf7iM6PdsqfYCo3D1tnAR65SkLENyWJG1deUzpRMQmbh9F3ea"

	if !IsValidAddress(zPub) {
		t.Errorf("%s was not considered valid", zPub)
	}
	if !IsValidAddress(zSec) {
		t.Errorf("%s was not considered valid", zSec)
	}

	if IsValidAddress(badEmpty) {
		t.Errorf("%s was considered valid", badEmpty)
	}
	if IsValidAddress(badLen) {
		t.Errorf("%s was considered valid", badLen)
	}
	if IsValidAddress(badPrePub) {
		t.Errorf("%s was considered valid", badPrePub)
	}
	if IsValidAddress(badPreSec) {
		t.Errorf("%s was considered valid", badPreSec)
	}
	if IsValidAddress(badCheckPub) {
		t.Errorf("%s was considered valid", badCheckPub)
	}
	if IsValidAddress(badCheckSec) {
		t.Errorf("%s was considered valid", badCheckSec)
	}
}

func TestGetECAddress(t *testing.T) {
	zSec := "Es2Rf7iM6PdsqfYCo3D1tnAR65SkLENyWJG1deUzpRMQmbh9F3eG"
	e, err := GetECAddress(zSec)
	if err != nil {
		t.Error(err)
	}

	// verify that the keys work
	msg := []byte("Hello Factom!")
	sig := ed.Sign(e.SecFixed(), msg)
	if !ed.Verify(e.PubFixed(), msg, sig) {
		t.Errorf("Key signature did not match")
	}
}

func TestIsValidFactoidAddress(t *testing.T) {
	zPub := "EC1m9mouvUQeEidmqpUYpYtXg8fvTYi6GNHaKg8KMLbdMBrFfmUa"
	zSec := "Es2Rf7iM6PdsqfYCo3D1tnAR65SkLENyWJG1deUzpRMQmbh9F3eG"
	badEmpty := ""
	badLen := "FA1y5ZGuHSLmf2TqNf6hVMkPiNGyQpQDFJvDLRkKQaoPo4bmbgu"
	badPrePub := "Fe1y5ZGuHSLmf2TqNf6hVMkPiNGyQpQDTFJvDLRkKQaoPo4bmbgu"
	badPreSec := "Fb1KWJrpLdfucvmYwN2nWrwepLn8ercpMbzXshd1g8zyhKXLVLWj"
	badCheckPub := "FA1y5ZGuHSLmf2TqNf6hVMkPiNGyQpQDTFJvDLRkKQaoPo4bmggu"
	badCheckSec := "Fs1KWJrpLdfucvmYwN2nWrwepLn8ercpMbzXshd1g8zyhKXLVLwj"

	if !IsValidAddress(zPub) {
		t.Errorf("%s was not considered valid", zPub)
	}
	if !IsValidAddress(zSec) {
		t.Errorf("%s was not considered valid", zSec)
	}

	if IsValidAddress(badEmpty) {
		t.Errorf("%s was considered valid", badEmpty)
	}
	if IsValidAddress(badLen) {
		t.Errorf("%s was considered valid", badLen)
	}
	if IsValidAddress(badPrePub) {
		t.Errorf("%s was considered valid", badPrePub)
	}
	if IsValidAddress(badPreSec) {
		t.Errorf("%s was considered valid", badPreSec)
	}
	if IsValidAddress(badCheckPub) {
		t.Errorf("%s was considered valid", badCheckPub)
	}
	if IsValidAddress(badCheckSec) {
		t.Errorf("%s was considered valid", badCheckSec)
	}
}

func TestGetFactoidAddress(t *testing.T) {
	zSec := "Fs1KWJrpLdfucvmYwN2nWrwepLn8ercpMbzXshd1g8zyhKXLVLWj"

	if _, err := GetFactoidAddress(zSec); err != nil {
		t.Error(err)
	}

	// ? test factoid key validity here
}

func TestMakeFactoidAddressFromMnemonic(t *testing.T) {
	m := "yellow yellow yellow yellow yellow yellow yellow yellow yellow yellow yellow yellow"
	cannonAdr := "FA3cih2o2tjEUsnnFR4jX1tQXPpSXFwsp3rhVp6odL5PNCHWvZV1"

	fct, err := MakeFactoidAddressFromKoinify(m)
	if err != nil {
		t.Error(err)
	}

	if fct.String() != cannonAdr {
		t.Errorf(
			"incorrect factoid address from 12 words: got %s expecting %s",
			fct.String(), cannonAdr)
	}
}

func TestMakeBIP44FactoidAddress(t *testing.T) {
	m := "yellow yellow yellow yellow yellow yellow yellow yellow yellow yellow yellow yellow"
	cannonAdr := "FA22de5NSG2FA2HmMaD4h8qSAZAJyztmmnwgLPghCQKoSekwYYct"

	fct, err := MakeBIP44FactoidAddress(m, bip32.FirstHardenedChild, 0, 0)
	if err != nil {
		t.Error(err)
	}

	if fct.String() != cannonAdr {
		t.Errorf(
			"incorrect factoid address from 12 words: got %s expecting %s",
			fct.String(), cannonAdr)
	}
}

func TestParseAndValidateMnemonic(t *testing.T) {
	ms := []string{
		"yellow yellow yellow yellow yellow yellow yellow yellow yellow yellow yellow yellow",   //valid
		"yellow  yellow yellow yellow yellow yellow yellow yellow yellow yellow yellow yellow",  //extra space
		"YELLOW yellow yellow yellow yellow yellow yellow yellow yellow yellow yellow yellow",   //capitalization
		" yellow yellow yellow yellow yellow yellow yellow yellow yellow yellow yellow yellow ", //spaces on sides
	}
	for i, m := range ms {
		_, err := ParseAndValidateMnemonic(m)
		if err != nil {
			t.Errorf("Error for mnemonic %v - `%v` - err", i, m)
		}
	}
}
