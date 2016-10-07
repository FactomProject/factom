package wallet_test

import (
	"strings"
	"testing"

	. "github.com/FactomProject/factom/wallet"
)

/*
func TestWalletDBO(t *testing.T) {
	db := NewMapDB()

}*/

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
	if ec.String() != "EC2HAuRUAeK9aBAafkWwmPPbfHFBF67pa71j7sZ8PdFTvZtuzqfF" {
		t.Errorf("%v", ec.String())
	}
	ec, err = seed.NextECAddress()
	if err != nil {
		t.Errorf("%v", err)
	}
	if ec.String() != "EC3aCUr3PgrzSXZ2dczD4raXmxcqCxvk1EvwVVKwbB9fx71yA9pp" {
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
