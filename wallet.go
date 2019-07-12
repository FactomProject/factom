// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"encoding/json"
	"fmt"
)

// BackupWallet returns a formatted string with the wallet seed and the secret
// keys for all of the wallet addresses.
func BackupWallet() (string, error) {
	req := NewJSON2Request("wallet-backup", APICounter(), nil)
	resp, err := walletRequest(req)
	if err != nil {
		return "", err
	}
	if resp.Error != nil {
		return "", resp.Error
	}

	w := new(struct {
		Seed         string             `json:"wallet-seed"`
		Addresses    []*addressResponse `json:"addresses"`
		IdentityKeys []*addressResponse `json:"identity-keys"`
	})
	if err := json.Unmarshal(resp.JSONResult(), w); err != nil {
		return "", err
	}

	s := fmt.Sprintln(w.Seed)
	s += fmt.Sprintln()
	for _, adr := range w.Addresses {
		s += fmt.Sprintln(adr.Public)
		s += fmt.Sprintln(adr.Secret)
		s += fmt.Sprintln()
	}
	for _, k := range w.IdentityKeys {
		s += fmt.Sprintln(k.Public)
		s += fmt.Sprintln(k.Secret)
		s += fmt.Sprintln()
	}
	return s, nil
}

// GenerateFactoidAddress creates a new Factoid Address and stores it in the
// Factom Wallet.
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

// GenerateECAddress creates a new Entry Credit Address and stores it in the
// Factom Wallet.
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

func GenerateIdentityKey() (*IdentityKey, error) {
	req := NewJSON2Request("generate-identity-key", APICounter(), nil)
	resp, err := walletRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	k := new(addressResponse)
	if err := json.Unmarshal(resp.JSONResult(), k); err != nil {
		return nil, err
	}
	e, err := GetIdentityKey(k.Secret)
	if err != nil {
		return nil, err
	}

	return e, nil
}

// ImportAddresses takes a number of Factoid and Entry Creidit secure keys and
// stores the Facotid and Entry Credit addresses in the Factom Wallet.
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

// ImportKoinify creates a Factoid Address from a secret 12 word koinify
// mnumonic.
//
// This functionality is used only to recover addresses that were funded by the
// Factom Genisis block to pay participants in the initial Factom network crowd
// funding.
func ImportKoinify(mnemonic string) (*FactoidAddress, error) {
	params := &struct {
		Words string `json:"words"`
	}{
		Words: mnemonic,
	}

	req := NewJSON2Request("import-koinify", APICounter(), params)
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

// RemoveAddress removes an address from the Factom Wallet database.
// (Be careful!)
func RemoveAddress(address string) error {
	params := new(addressRequest)
	params.Address = address

	req := NewJSON2Request("remove-address", APICounter(), params)
	resp, err := walletRequest(req)
	if err != nil {
		return err
	}
	if resp.Error != nil {
		return resp.Error
	}

	return nil
}

// FetchAddresses requests all of the addresses in the Factom Wallet database.
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

// FetchECAddress requests an Entry Credit address from the Factom Wallet.
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

// FetchFactoidAddress requests a Factom address from the Factom Wallet.
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

func ImportIdentityKeys(pubs ...string) ([]*IdentityKey, error) {
	params := new(struct {
		IdentityKeys []secretRequest `json:"keys"`
	})
	for _, pub := range pubs {
		s := secretRequest{Secret: pub}
		params.IdentityKeys = append(params.IdentityKeys, s)
	}

	req := NewJSON2Request("import-identity-keys", APICounter(), params)
	resp, err := walletRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	r := new(multiIdentityKeyResponse)
	if err := json.Unmarshal(resp.JSONResult(), r); err != nil {
		return nil, err
	}
	keys := make([]*IdentityKey, 0)
	for _, v := range r.IdentityKeys {
		k, err := GetIdentityKey(v.Secret)
		if err != nil {
			return nil, err
		}
		keys = append(keys, k)
	}

	return keys, nil
}

func FetchIdentityKey(pub string) (*IdentityKey, error) {
	params := new(struct {
		Public string `json:"public"`
	})
	params.Public = pub

	req := NewJSON2Request("identity-key", APICounter(), params)
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

	return GetIdentityKey(r.Secret)
}

func FetchIdentityKeys() ([]*IdentityKey, error) {
	req := NewJSON2Request("all-identity-keys", APICounter(), nil)
	resp, err := walletRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	keys := make([]*IdentityKey, 0)

	multiKeyResp := new(multiIdentityKeyResponse)
	if err := json.Unmarshal(resp.JSONResult(), multiKeyResp); err != nil {
		return nil, err
	}

	for _, v := range multiKeyResp.IdentityKeys {
		if IdentityKeyStringType(v.Public) == IDPub {
			k, err := GetIdentityKey(v.Secret)
			if err != nil {
				return nil, err
			}
			keys = append(keys, k)
		} else {
			return nil, fmt.Errorf("%s is not a valid public identity key", v.Public)
		}
	}

	return keys, nil
}

func RemoveIdentityKey(pub string) error {
	params := new(struct {
		Public string `json:"public"`
	})
	params.Public = pub

	req := NewJSON2Request("remove-identity-key", APICounter(), params)
	resp, err := walletRequest(req)
	if err != nil {
		return err
	}
	if resp.Error != nil {
		return resp.Error
	}

	return nil
}

// GetWalletHeight requests the current block heights known to the Factom
// Wallet.
func GetWalletHeight() (uint32, error) {
	req := NewJSON2Request("get-height", APICounter(), nil)
	resp, err := walletRequest(req)
	if err != nil {
		return 0, err
	}
	if resp.Error != nil {
		return 0, resp.Error
	}

	r := new(heightResponse)
	if err := json.Unmarshal(resp.JSONResult(), r); err != nil {
		return 0, err
	}

	return uint32(r.Height), nil
}

func UnlockWallet(passphrase string, seconds int64) (int64, error) {
	req := NewJSON2Request("unlock-wallet", APICounter(), &passphraseRequest{Password: passphrase, Timeout: seconds})
	resp, err := walletRequest(req)
	if err != nil {
		return 0, err
	}
	if resp.Error != nil {
		return 0, resp.Error
	}

	r := new(unlockResponse)
	if err := json.Unmarshal(resp.JSONResult(), r); err != nil {
		return 0, err
	}

	return r.UnlockedUntil, nil
}

type addressResponse struct {
	Public string `json:"public"`
	Secret string `json:"secret"`
}

type multiAddressResponse struct {
	Addresses []*addressResponse `json:"addresses"`
}

type multiIdentityKeyResponse struct {
	IdentityKeys []*addressResponse `json:"keys"`
}

type composeEntryRequest struct {
	Entry Entry  `json:"entry"`
	ECPub string `json:"ecpub"`
	Force bool   `json:"force"`
}

type composeChainRequest struct {
	Chain Chain  `json:"chain"`
	ECPub string `json:"ecpub"`
	Force bool   `json:"force"`
}

type composeEntryResponse struct {
	Commit *JSON2Request `json:"commit"`
	Reveal *JSON2Request `json:"reveal"`
}

// WalletComposeChainCommitReveal composes commit and reveal json objects that
// may be used to make API calls to the factomd API to create a new Factom
// Chain.
//
// WalletComposeChainCommitReveal may be used by an offline wallet to create the
// calls needed to create new chains while keeping addresses secure in an
// offline wallet.
func WalletComposeChainCommitReveal(chain *Chain, ecPub string, force bool) (*JSON2Request, *JSON2Request, error) {
	params := new(composeChainRequest)
	params.Chain = *chain
	params.ECPub = ecPub
	params.Force = force

	req := NewJSON2Request("compose-chain", APICounter(), params)
	resp, err := walletRequest(req)
	if err != nil {
		return nil, nil, err
	}
	if resp.Error != nil {
		return nil, nil, resp.Error
	}

	r := new(composeEntryResponse)
	if err := json.Unmarshal(resp.JSONResult(), r); err != nil {
		return nil, nil, err
	}

	return r.Commit, r.Reveal, nil
}

// WalletComposeEntryCommitReveal composes commit and reveal json objects that
// may be used to make API calls to the factomd API to create a new Factom
// Entry.
//
// WalletComposeEntryCommitReveal may be used by an offline wallet to create the
// calls needed to create new entries while keeping addresses secure in an
// offline wallet.
func WalletComposeEntryCommitReveal(entry *Entry, ecPub string, force bool) (*JSON2Request, *JSON2Request, error) {
	params := new(composeEntryRequest)
	params.Entry = *entry
	params.ECPub = ecPub
	params.Force = force

	req := NewJSON2Request("compose-entry", APICounter(), params)
	resp, err := walletRequest(req)
	if err != nil {
		return nil, nil, err
	}
	if resp.Error != nil {
		return nil, nil, resp.Error
	}

	r := new(composeEntryResponse)
	if err := json.Unmarshal(resp.JSONResult(), r); err != nil {
		return nil, nil, err
	}

	return r.Commit, r.Reveal, nil
}

type heightResponse struct {
	Height int64 `json:"height"`
}
