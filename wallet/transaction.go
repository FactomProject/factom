// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package wallet

import (
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"

	"github.com/FactomProject/btcutil/base58"
	ed "github.com/FactomProject/ed25519"
	"github.com/FactomProject/factoid"
	"github.com/FactomProject/factom"
)

var (
	ErrTXExists      = errors.New("wallet: Transaction name already exists")
	ErrTXNotExists   = errors.New("wallet: Transaction name was not found")
	ErrTXInvalidName = errors.New("wallet: Transaction name is not valid")
)

func (w *Wallet) NewTransaction(name string) error {
	if _, exists := w.transactions[name]; exists {
		return ErrTXExists
	}
	
	// check that the transaction name is valid
	if name == "" {
		return ErrTXInvalidName
	}
	if len(name) > 32 {
		return ErrTXInvalidName
	}
	if match, err := regexp.MatchString("[^a-zA-Z0-9_-]", name); err != nil {
		return err
	} else if match {
		return ErrTXInvalidName
	}
	
	t := new(factoid.Transaction)
	t.SetMilliTimestamp(milliTime())
	w.transactions[name] = t
	return nil
}

func (w *Wallet) DeleteTransaction(name string) error {
	if _, exists := w.transactions[name]; !exists {
		return ErrTXNotExists
	}
	delete(w.transactions, name)
	return nil
}

func (w *Wallet) AddInput(name, address string, amount uint64) error {
	if _, exists := w.transactions[name]; !exists {
		return ErrTXNotExists
	}
	trans := w.transactions[name]

	a, err := w.GetFCTAddress(address)
	if err != nil {
		return err
	}
	adr := factoid.NewAddress(a.RCDHash())

	// First look if this is really an update
	for _, input := range trans.GetInputs() {
		if input.GetAddress().IsSameAs(adr) {
			input.SetAmount(amount)
			return nil
		}
	}

	// Add our new input
	trans.AddInput(adr, amount)

	return nil
}

func (w *Wallet) AddOutput(name, address string, amount uint64) error {
	if _, exists := w.transactions[name]; !exists {
		return ErrTXNotExists
	}
	trans := w.transactions[name]

	if !factom.IsValidAddress(address) {
		return errors.New("Invalid Address")
	}

	adr := factoid.NewAddress(base58.Decode(address)[2:34])

	trans.AddOutput(adr, amount)

	return nil
}

func (w *Wallet) AddECOutput(name, address string, amount uint64) error {
	if _, exists := w.transactions[name]; !exists {
		return ErrTXNotExists
	}
	trans := w.transactions[name]

	if !factom.IsValidAddress(address) {
		return errors.New("Invalid Address")
	}

	adr := factoid.NewAddress(base58.Decode(address)[2:34])

	trans.AddECOutput(adr, amount)

	return nil
}

func (w *Wallet) AddFee(name, address string, rate uint64) error {
	if _, exists := w.transactions[name]; !exists {
		return ErrTXNotExists
	}
	trans := w.transactions[name]

	{
		ins, err := trans.TotalInputs()
		if err != nil {
			return err
		}
		outs, err := trans.TotalOutputs()
		if err != nil {
			return err
		}
		ecs, err := trans.TotalECs()
		if err != nil {
			return err
		}

		if ins != outs+ecs {
			return fmt.Errorf("Inputs and outputs don't add up")
		}
	}

	transfee, err := trans.CalculateFee(rate)
	if err != nil {
		return err
	}

	a, err := w.GetFCTAddress(address)
	if err != nil {
		return err
	}
	adr := factoid.NewAddress(a.RCDHash())

	for _, input := range trans.GetInputs() {
		if input.GetAddress().IsSameAs(adr) {
			amt, err := factoid.ValidateAmounts(input.GetAmount(), transfee)
			if err != nil {
				return err
			}
			input.SetAmount(amt)
			return nil
		}
	}
	return fmt.Errorf("%s is not an input to the transaction.", address)
}

func (w *Wallet) SubFee(name, address string, rate uint64) error {
	if _, exists := w.transactions[name]; !exists {
		return ErrTXNotExists
	}
	trans := w.transactions[name]

	if !factom.IsValidAddress(address) {
		return errors.New("Invalid Address")
	}

	{
		ins, err := trans.TotalInputs()
		if err != nil {
			return err
		}
		outs, err := trans.TotalOutputs()
		if err != nil {
			return err
		}
		ecs, err := trans.TotalECs()
		if err != nil {
			return err
		}

		if ins != outs+ecs {
			return fmt.Errorf("Inputs and outputs don't add up")
		}
	}

	transfee, err := trans.CalculateFee(rate)
	if err != nil {
		return err
	}

	adr := factoid.NewAddress(base58.Decode(address)[2:34])

	for _, output := range trans.GetOutputs() {
		if output.GetAddress().IsSameAs(adr) {
			output.SetAmount(output.GetAmount() - transfee)
			return nil
		}
	}
	return fmt.Errorf("%s is not an output to the transaction.", address)
}

func (w *Wallet) SignTransaction(name string) error {
	if _, exists := w.transactions[name]; !exists {
		return ErrTXNotExists
	}
	trans := w.transactions[name]

	data, err := trans.MarshalBinarySig()
	if err != nil {
		return err
	}

	var errMsg []byte
	for i, rcd := range trans.GetRCDs() {
		rcd1, ok := rcd.(*factoid.RCD_1)
		if ok {
			pub := rcd1.GetPublicKey()
			key := base58.CheckEncodeWithVersionBytes(pub[:], 0x5f, 0xb1)
			adr, err := w.GetFCTAddress(key)
			if err != nil {
				errMsg = append(errMsg,
					[]byte("Do not have the private key for: "+
						factoid.ConvertFctAddressToUserStr(factoid.NewAddress(pub))+"\n")...)
			} else {
				sec := new([factoid.SIGNATURE_LENGTH]byte)
				copy(sec[:], adr.SecBytes())
				bsig := ed.Sign(sec, data)
				sig := new(factoid.Signature)
				sig.SetSignature(bsig[:])
				sigblk := new(factoid.SignatureBlock)
				sigblk.AddSignature(sig)
				trans.SetSignatureBlock(i, sigblk)
			}
		}
	}
	if errMsg != nil {
		return errors.New(string(errMsg))
	}

	return nil
}

func (w *Wallet) GetTransactions() map[string]factoid.ITransaction {
	return w.transactions
}

func (w *Wallet) ComposeTransaction(name string) (*factom.JSON2Request, error) {
	if _, exists := w.transactions[name]; !exists {
		return nil, ErrTXNotExists
	}
	trans := w.transactions[name]

	type txreq struct {
		Transaction string `json:"transaction"`
	}

	param := new(txreq)
	if p, err := trans.MarshalBinary(); err != nil {
		return nil, err
	} else {
		param.Transaction = hex.EncodeToString(p)
	}

	req := factom.NewJSON2Request("factoid-transaction", apiCounter(), param)

	return req, nil
}

// TODO ---
//
//func (w *Wallet) SendTransaction
