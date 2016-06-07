// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package wallet

import (
	"crypto/rand"
	"crypto/sha512"
	"fmt"
	"os"
	"sync"

	"github.com/FactomProject/factoid"
	"github.com/FactomProject/factom"
	"github.com/FactomProject/goleveldb/leveldb"
	"github.com/FactomProject/goleveldb/leveldb/opt"
	"github.com/FactomProject/goleveldb/leveldb/util"
)

// Database keys and key prefixes
var (
	fcDBPrefix    = []byte("Factoids")
	ecDBPrefix    = []byte("Entry Credits")
	seedDBKey     = []byte("DB Seed")
	nextSeedDBKey = []byte("Next Seed")
)

// Wallet is a connection to a Factom Wallet Database
type Wallet struct {
	lock         sync.RWMutex
	ldb          *leveldb.DB
	transactions map[string]factoid.ITransaction
}

// NewWallet creates a new Factom Wallet Database. It will return an error if a
// database already exists at the specified path.
func NewWallet(path string) (*Wallet, error) {
	o := &opt.Options{ErrorIfExist: true}
	w := new(Wallet)
	w.transactions = make(map[string]factoid.ITransaction)
	if l, err := leveldb.OpenFile(path, o); err != nil {
		return nil, err
	} else {
		w.ldb = l
	}

	// generate a random seed for new address generation in this wallet
	seed := make([]byte, 64)
	if n, err := rand.Read(seed); err != nil {
		return nil, err
	} else if n != 64 {
		return nil, fmt.Errorf("Wrong number of bytes read: %d", n)
	}
	w.ldb.Put(seedDBKey, seed, nil)
	w.ldb.Put(nextSeedDBKey, seed, nil)

	return w, nil
}

// OpenWallet opens an existing Factom Wallet Database. It will return an error
// if no database exists at the path.
func OpenWallet(path string) (*Wallet, error) {
	o := &opt.Options{ErrorIfMissing: true}
	w := new(Wallet)
	w.transactions = make(map[string]factoid.ITransaction)
	
	// open the db file
	if l, err := leveldb.OpenFile(path, o); err != nil {
		// try an recover the file if there is an error
		r, err := leveldb.RecoverFile(path, nil)
		if err != nil {
			return nil, err
		}
		w.ldb = r		
	} else {
		w.ldb = l
	}

	// make sure it is a wallet db
	if s, err := w.ldb.Has(seedDBKey, nil); err != nil {
		return nil, err
	} else if !s {
		return nil, fmt.Errorf("wallet is missing its seed")
	} else if n, err := w.ldb.Has(nextSeedDBKey, nil); err != nil {
		return nil, err
	} else if !n {
		return nil, fmt.Errorf("wallet is missing its next seed")
	}
	
	return w, nil
}

func NewOrOpenWallet(path string) (*Wallet, error) {
	w, err := NewWallet(path)
	if err != nil {
		if !os.IsExist(err) {
			return nil, err
		}
		return OpenWallet(path)
	}
	return w, err
}

// Close closes a Factom Wallet Database
func (w *Wallet) Close() error {
	w.lock.Lock()
	defer w.lock.Unlock()

	return w.ldb.Close()
}

// GenerateECAddress creates and stores a new Entry Credit Address in the
// Wallet. The address can be reproduced in the future using the Wallet Seed.
func (w *Wallet) GenerateECAddress() (*factom.ECAddress, error) {
	// get the next seed from the db
	seed, err := w.getNextSeed()
	if err != nil {
		return nil, err
	}

	// create the new seed
	newseed := sha512.Sum512(seed)
	a, err := factom.MakeECAddress(newseed[:32])
	if err != nil {
		return nil, err
	}

	// save the new seed and the address in the db
	if err := w.putNextSeed(newseed[:]); err != nil {
		return nil, err
	}

	if err := w.PutECAddress(a); err != nil {
		return nil, err
	}

	return a, nil
}

// GenerateFCTAddress creates and stores a new Factoid Address in the Wallet.
// The address can be reproduced in the future using the Wallet Seed.
func (w *Wallet) GenerateFCTAddress() (*factom.FactoidAddress, error) {
	// get the next seed from the db
	seed, err := w.getNextSeed()
	if err != nil {
		return nil, err
	}

	// create the new seed
	newseed := sha512.Sum512(seed)
	a, err := factom.MakeFactoidAddress(newseed[:32])
	if err != nil {
		return nil, err
	}

	// save the new seed and the address in the db
	if err := w.putNextSeed(newseed[:]); err != nil {
		return nil, err
	}

	if err := w.PutFCTAddress(a); err != nil {
		return nil, err
	}

	return a, nil
}

// GetAllAddresses retrieves all Entry Credit and Factoid Addresses from the
// Wallet Database.
func (w *Wallet) GetAllAddresses() ([]*factom.FactoidAddress, []*factom.ECAddress, error) {
	fcs := make([]*factom.FactoidAddress, 0)
	for iter := w.ldb.NewIterator(util.BytesPrefix(fcDBPrefix), nil); iter.Next(); {
		f, err := factom.GetFactoidAddress(string(iter.Value()))
		if err != nil {
			return nil, nil, err
		}
		fcs = append(fcs, f)
	}
	
	ecs := make([]*factom.ECAddress, 0)
	for iter := w.ldb.NewIterator(util.BytesPrefix(ecDBPrefix), nil); iter.Next(); {
		e, err := factom.GetECAddress(string(iter.Value()))
		if err != nil {
			return nil, nil, err
		}
		ecs = append(ecs, e)
	}

	return fcs, ecs, nil
}

// GetECAddress retrieves a specific Entry Credit Address from the Wallet using
// the Public Address String.
func (w *Wallet) GetECAddress(a string) (*factom.ECAddress, error) {
	w.lock.RLock()
	defer w.lock.RUnlock()

	key := append(ecDBPrefix, a...)
	p, err := w.ldb.Get(key, nil)
	if err != nil {
		return nil, err
	}
	e, err := factom.GetECAddress(string(p))
	if err != nil {
		return nil, err
	}
	return e, nil
}

// GetFCTAddress retrieves a specific Factoid Address from the Wallet using the
// Public Address String.
func (w *Wallet) GetFCTAddress(a string) (*factom.FactoidAddress, error) {
	w.lock.RLock()
	defer w.lock.RUnlock()

	key := append(fcDBPrefix, a...)
	p, err := w.ldb.Get(key, nil)
	if err != nil {
		return nil, err
	}
	f, err := factom.GetFactoidAddress(string(p))
	if err != nil {
		return nil, err
	}
	return f, nil
}

// PutECAddress stores an Entry Credit Address in the Wallet Database.
func (w *Wallet) PutECAddress(e *factom.ECAddress) error {
	key := append(ecDBPrefix, e.PubString()...)
	return w.ldb.Put(key, []byte(e.SecString()), nil)
}

// PutFCTAddress stores a Factoid Address in the Wallet Database.
func (w *Wallet) PutFCTAddress(f *factom.FactoidAddress) error {
	key := append(fcDBPrefix, f.PubString()...)
	return w.ldb.Put(key, []byte(f.SecString()), nil)
}

// GetSeed returns the string representaion of the Wallet Seed. The Wallet Seed
// can be used to regenerate the Factoid and Entry Credit Addresses previously
// generated by the wallet. Note that Addresses that are imported into the
// Wallet cannot be regenerated using the Wallet Seed.
func (w *Wallet) GetSeed() (string, error) {
	seed, err := w.getSeed()
	if err != nil {
		return "", err
	}
	
	return seedString(seed), nil
}

// getSeed retrieves the raw Wallet Seed from the Database.
func (w *Wallet) getSeed() ([]byte, error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	p, err:=  w.ldb.Get(seedDBKey, nil)
	if err != nil {
		return nil, err
	} else if len(p) != SeedLength {
		return nil, fmt.Errorf("Wallet Seed is the wrong length: %d", len(p))
	}

	return p, nil
}

// getNextSeed retrieves the raw Next Wallet Seed from the Database.
func (w *Wallet) getNextSeed() ([]byte, error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	p, err := w.ldb.Get(nextSeedDBKey, nil)
	if err != nil {
		return nil, err
	} else if len(p) != SeedLength {
		return nil, fmt.Errorf("Wallet Seed is the wrong length: %d", len(p))
	}

	return p, nil
}

// putNextSeed stores the Next Wallet Seed in the Wallet Database.
func (w *Wallet) putNextSeed(seed []byte) error {
	if len(seed) != SeedLength {
		return fmt.Errorf("Provided Seed is the wrong length: %d", len(seed))
	}
	return w.ldb.Put(nextSeedDBKey, seed, nil)
}
