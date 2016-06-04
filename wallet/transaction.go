// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package wallet

import (
	"errors"

	"github.com/FactomProject/btcutil/base58"
	"github.com/FactomProject/factoid"
	"github.com/FactomProject/factom"
)

var (
	ErrTXExists = errors.New("wallet: Transaction name already exists")
	ErrTXNotExists = errors.New("wallet: Transaction name was not found")
)

func (w *Wallet) CreateTransaction(name string) error {
	if _, exists := w.transactions[name]; exists {
		return ErrTXExists
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

func (w *Wallet) AddInput (name string, address *factom.FactoidAddress, amount uint64) error {
	if _, exists := w.transactions[name]; !exists {
		return ErrTXNotExists
	}
	trans := w.transactions[name]

	adr := factoid.NewAddress(address.RCDHash())
	
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

// TODO func (w *Wallet) UpdateInput

func (w *Wallet) AddOutput (name, address string, amount uint64) error {
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

func (w *Wallet) AddECOutput (name, address string, amount uint64) error {
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

//func (w *Wallet) AddFee
//
//func (w *Wallet) SubFee
//
//func (w *Wallet) SignTransaction
//
////func (w *Wallet) SendTransaction
//
//func (w *Wallet) ComposeTransaction
//
//func (w *Wallet) ListTransactions
