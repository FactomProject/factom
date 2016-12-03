// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/FactomProject/btcutil/base58"
	ed "github.com/FactomProject/ed25519"
)

// BackupWallet returns a formatted string with the wallet seed and the secret
// keys for all of the wallet addresses.
func BackupWallet() (string, error) {
	type walletBackupResponse struct {
		Seed      string             `json:"wallet-seed"`
		Addresses []*addressResponse `json:"addresses"`
	}

	req := NewJSON2Request("wallet-backup", APICounter(), nil)
	resp, err := walletRequest(req)
	if err != nil {
		return "", err
	}
	if resp.Error != nil {
		return "", resp.Error
	}

	w := new(walletBackupResponse)
	if err := json.Unmarshal(resp.JSONResult(), w); err != nil {
		return "", err
	}

	s := fmt.Sprintln("Seed:")
	s += fmt.Sprintln(w.Seed)
	s += fmt.Sprintln("Addresses:")
	for _, adr := range w.Addresses {
		s += fmt.Sprintln(adr.Secret)
	}
	return s, nil
}

func GenerateFactoidAddress() (*FactoidAddress, error) {
	req := NewJSON2Request("generate-factoid-address", APICounter(), nil)
	resp, err := walletRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	a := new(addressResponse)
	if err := json.Unmarshal(resp.JSONResult(), a); err != nil {
		return nil, err
	}
	f, err := GetFactoidAddress(a.Secret)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func GenerateECAddress() (*ECAddress, error) {
	req := NewJSON2Request("generate-ec-address", APICounter(), nil)
	resp, err := walletRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	a := new(addressResponse)
	if err := json.Unmarshal(resp.JSONResult(), a); err != nil {
		return nil, err
	}
	e, err := GetECAddress(a.Secret)
	if err != nil {
		return nil, err
	}

	return e, nil
}

func ImportAddresses(addrs ...string) (
	[]*FactoidAddress,
	[]*ECAddress,
	error) {

	params := new(importRequest)
	for _, addr := range addrs {
		s := secretRequest{Secret: addr}
		params.Addresses = append(params.Addresses, s)
	}
	req := NewJSON2Request("import-addresses", APICounter(), params)
	resp, err := walletRequest(req)
	if err != nil {
		return nil, nil, err
	}
	if resp.Error != nil {
		return nil, nil, resp.Error
	}

	r := new(multiAddressResponse)
	if err := json.Unmarshal(resp.JSONResult(), r); err != nil {
		return nil, nil, err
	}
	fs := make([]*FactoidAddress, 0)
	es := make([]*ECAddress, 0)
	for _, address := range r.Addresses {
		switch AddressStringType(address.Secret) {
		case FactoidSec:
			f, err := GetFactoidAddress(address.Secret)
			if err != nil {
				return nil, nil, err
			}
			fs = append(fs, f)
		case ECSec:
			e, err := GetECAddress(address.Secret)
			if err != nil {
				return nil, nil, err
			}
			es = append(es, e)
		default:
			return nil, nil, fmt.Errorf("Unrecognized address type")
		}
	}

	return fs, es, nil
}

func ImportMnemonic(mnemonic string) (*FactoidAddress, error) {
	params := new(importMnemonicRequest)
	params.Words = mnemonic

	req := NewJSON2Request("import-mnemonic", APICounter(), params)
	resp, err := walletRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	r := new(addressResponse)
	if err := json.Unmarshal(resp.JSONResult(), r); err != nil {
		return nil, err
	}
	f, err := GetFactoidAddress(r.Secret)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func FetchAddresses() ([]*FactoidAddress, []*ECAddress, error) {
	req := NewJSON2Request("all-addresses", APICounter(), nil)
	resp, err := walletRequest(req)
	if err != nil {
		return nil, nil, err
	}
	if resp.Error != nil {
		return nil, nil, resp.Error
	}

	fs := make([]*FactoidAddress, 0)
	es := make([]*ECAddress, 0)

	as := new(multiAddressResponse)
	if err := json.Unmarshal(resp.JSONResult(), as); err != nil {
		return nil, nil, err
	}

	for _, adr := range as.Addresses {
		switch AddressStringType(adr.Public) {
		case FactoidPub:
			f, err := GetFactoidAddress(adr.Secret)
			if err != nil {
				return nil, nil, err
			}
			fs = append(fs, f)
		case ECPub:
			e, err := GetECAddress(adr.Secret)
			if err != nil {
				return nil, nil, err
			}
			es = append(es, e)
		default:
			return nil, nil, fmt.Errorf("%s is not a valid address", adr.Public)
		}
	}

	return fs, es, nil
}

func FetchECAddress(ecpub string) (*ECAddress, error) {
	if AddressStringType(ecpub) != ECPub {
		return nil, fmt.Errorf(
			"%s is not an Entry Credit Public Address", ecpub)
	}
	params := new(addressRequest)
	params.Address = ecpub

	req := NewJSON2Request("address", APICounter(), params)
	resp, err := walletRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	r := new(addressResponse)
	if err := json.Unmarshal(resp.JSONResult(), r); err != nil {
		return nil, err
	}

	return GetECAddress(r.Secret)
}

func FetchFactoidAddress(fctpub string) (*FactoidAddress, error) {
	if AddressStringType(fctpub) != FactoidPub {
		return nil, fmt.Errorf("%s is not a Factoid Address", fctpub)
	}
	params := new(addressRequest)
	params.Address = fctpub

	req := NewJSON2Request("address", APICounter(), params)
	resp, err := walletRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	r := new(addressResponse)
	if err := json.Unmarshal(resp.JSONResult(), r); err != nil {
		return nil, err
	}

	return GetFactoidAddress(r.Secret)
}

func SignMessage(p string, m string) (string, string, error) {
	if AddressStringType(p) != FactoidPub && AddressStringType(p) != ECPub {
		return "", "", fmt.Errorf("%s is neither a Factoid Address nor an EntryCredit Address", p)
	}

	var (
		sig            *[ed.SignatureSize]byte
		pubKeyPrefixed string
	)

	pubKey := new(bytes.Buffer)

	switch AddressStringType(p) {
	case FactoidPub:
		addr, err := FetchFactoidAddress(p)
		if err != nil {
			return "", "", err
		}

		sig = ed.Sign(addr.Sec, []byte(m))

		pubKey.Write(fcPubPrefix)
		pubKey.Write(addr.PubBytes())
		pubKeyPrefixed = base64.StdEncoding.EncodeToString(pubKey.Bytes())
	case ECPub:
		addr, err := FetchECAddress(p)
		if err != nil {
			return "", "", err
		}

		sig = ed.Sign(addr.Sec, []byte(m))

		pubKey.Write(ecPubPrefix)
		pubKey.Write(addr.PubBytes())
		pubKeyPrefixed = base64.StdEncoding.EncodeToString(pubKey.Bytes())
	}

	return pubKeyPrefixed, base64.StdEncoding.EncodeToString(sig[:]), nil
}

func VerifyMessage(p string, s string, m string) (bool, string, error) {
	var addr string

	pubBytesPrefixed, err := base64.StdEncoding.DecodeString(p)
	if err != nil {
		fmt.Println("Error reading public key")
		return false, "", err
	}

	prefix := pubBytesPrefixed[:PrefixLength]
	pubBytes := pubBytesPrefixed[PrefixLength:]

	if bytes.Equal(prefix, ecPubPrefix) {
		buf := new(bytes.Buffer)
		buf.Write(pubBytesPrefixed)

		check := shad(buf.Bytes())[:ChecksumLength]
		buf.Write(check)

		addr = base58.Encode(buf.Bytes())

	} else if bytes.Equal(prefix, fcPubPrefix) {
		a := new(FactoidAddress)
		r := NewRCD1()
		copy(r.Pub[:], pubBytesPrefixed[PrefixLength:])
		a.RCD = r

		addr = a.String()

	} else {
		fmt.Println("Public key is not valid")
		return false, "", nil
	}

	sigDecoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return false, "", err
	}

	sig := new([ed.SignatureSize]byte)
	pub := new([ed.PublicKeySize]byte)
	copy(sig[:], sigDecoded)
	copy(pub[:], pubBytes)
	res := ed.Verify(pub, []byte(m), sig)

	return res, addr, nil
}

type addressResponse struct {
	Public string `json:"public"`
	Secret string `json:"secret"`
}

type multiAddressResponse struct {
	Addresses []*addressResponse `json:"addresses"`
}
