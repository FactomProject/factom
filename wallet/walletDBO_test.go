package wallet_test

import (
	"strings"
	"testing"

	. "github.com/FactomProject/factom/wallet"
)

func TestWalletDBO(t *testing.T) {
	db := NewMapDB()
	seed, err := db.GetOrCreateDBSeed()
	if err != nil {
		t.Errorf("%v", err)
	}

	seed.MnemonicSeed = "yellow yellow yellow yellow yellow yellow yellow yellow yellow yellow yellow yellow"

	err = db.InsertDBSeed(seed)
	if err != nil {
		t.Errorf("%v", err)
	}

	ec, err := db.GetNextECAddress()
	if err != nil {
		t.Errorf("%v", err)
	}
	if ec.String() != "EC2KnJQN86MYq4pQyeSGTHSiVdkhRCPXS3udzD4im6BXRBjZFMmR" {
		t.Errorf("%v", ec.String())
	}
	ec, err = db.GetNextECAddress()
	if err != nil {
		t.Errorf("%v", err)
	}
	if ec.String() != "EC2UNG5LztGN3BNiVMEgkBP8ra8ud3HjjWWXKjrQozJ98rTvXKYy" {
		t.Errorf("%v", ec.String())
	}

	f, err := db.GetNextFCTAddress()
	if err != nil {
		t.Errorf("%v", err)
	}
	if f.String() != "FA22de5NSG2FA2HmMaD4h8qSAZAJyztmmnwgLPghCQKoSekwYYct" {
		t.Errorf("%v", f.String())
	}
	f, err = db.GetNextFCTAddress()
	if err != nil {
		t.Errorf("%v", err)
	}
	if f.String() != "FA3heCmxKCk1tCCfiAMDmX8Ctg6XTQjRRaJrF5Jagc9rbo7wqQLV" {
		t.Errorf("%v", f.String())
	}
}

func TestDBSeed(t *testing.T) {
	seed, err := NewRandomSeed()
	if err != nil {
		t.Errorf("%v", err)
	}
	if seed.MnemonicSeed == "" {
		t.Errorf("Empty mnemonic seed returned")
	}
	l := len(strings.Fields(seed.MnemonicSeed))
	if l < 12 {
		t.Errorf("Not enough words in mnemonic. Expecitng 12, found %d", l)
	}
	if l > 12 {
		t.Errorf("Too many words in mnemonic. Expecitng 12, found %d", l)
	}

	seed.MnemonicSeed = "yellow yellow yellow yellow yellow yellow yellow yellow yellow yellow yellow yellow"

	ec, err := seed.NextECAddress()
	if err != nil {
		t.Errorf("%v", err)
	}
	if ec.String() != "EC2KnJQN86MYq4pQyeSGTHSiVdkhRCPXS3udzD4im6BXRBjZFMmR" {
		t.Errorf("%v", ec.String())
	}
	ec, err = seed.NextECAddress()
	if err != nil {
		t.Errorf("%v", err)
	}
	if ec.String() != "EC2UNG5LztGN3BNiVMEgkBP8ra8ud3HjjWWXKjrQozJ98rTvXKYy" {
		t.Errorf("%v", ec.String())
	}

	f, err := seed.NextFCTAddress()
	if err != nil {
		t.Errorf("%v", err)
	}
	if f.String() != "FA22de5NSG2FA2HmMaD4h8qSAZAJyztmmnwgLPghCQKoSekwYYct" {
		t.Errorf("%v", f.String())
	}
	f, err = seed.NextFCTAddress()
	if err != nil {
		t.Errorf("%v", err)
	}
	if f.String() != "FA3heCmxKCk1tCCfiAMDmX8Ctg6XTQjRRaJrF5Jagc9rbo7wqQLV" {
		t.Errorf("%v", f.String())
	}

	if seed.NextFactoidAddressIndex != 2 {
		t.Errorf("Wrong NextFactoidAddressIndex")
	}
	if seed.NextECAddressIndex != 2 {
		t.Errorf("Wrong NextECAddressIndex")
	}
}
