// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package wallet

import (
	"crypto/rand"
	"crypto/sha512"
	"fmt"
	"sync"

	"github.com/FactomProject/factom"
	"github.com/FactomProject/goleveldb/leveldb"
	"github.com/FactomProject/goleveldb/leveldb/opt"
	"github.com/FactomProject/goleveldb/leveldb/util"
)

var (
	fcDBPrefix    = []byte("Factoids")
	ecDBPrefix    = []byte("Entry Credits")
	seedDBKey     = []byte("DB Seed")
	nextSeedDBKey = []byte("Next Seed")
)

type Wallet struct {
	lock sync.RWMutex
	ldb  *leveldb.DB
}

func NewWallet(path string) (*Wallet, error) {
	o := &opt.Options{ErrorIfExist: true}
	wdb := new(Wallet)
	if l, err := leveldb.OpenFile(path, o); err != nil {
		return nil, err
	} else {
		wdb.ldb = l
	}

	// generate a random seed for new address generation in this wallet
	seed := make([]byte, 64)
	if n, err := rand.Read(seed); err != nil {
		return nil, err
	} else if n != 64 {
		return nil, fmt.Errorf("Wrong number of bytes read: %d", n)
	}
	wdb.ldb.Put(seedDBKey, seed, nil)
	wdb.ldb.Put(nextSeedDBKey, seed, nil)

	return wdb, nil
}

func OpenWallet(path string) (*Wallet, error) {
	o := &opt.Options{ErrorIfMissing: true}
	wdb := new(Wallet)
	if l, err := leveldb.OpenFile(path, o); err != nil {
		wdb.ldb = l
		return nil, err
	}
	// TODO - validate database
	// ? - check if db is corrrupt and recover
	return wdb, nil
}

func (w *Wallet) Close() error {
	w.lock.Lock()
	defer w.lock.Unlock()

	return w.ldb.Close()
}

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

func (w *Wallet) PutECAddress(e *factom.ECAddress) error {
	key := append(ecDBPrefix, e.PubString()...)
	return w.ldb.Put(key, []byte(e.SecString()), nil)
}

func (w *Wallet) PutFCTAddress(f *factom.FactoidAddress) error {
	key := append(fcDBPrefix, f.PubString()...)
	return w.ldb.Put(key, []byte(f.SecString()), nil)
}

func (w *Wallet) getSeed() ([]byte, error) {
	w.lock.Lock()
	defer w.lock.Unlock()
	return w.ldb.Get(seedDBKey, nil)
}

func (w *Wallet) getNextSeed() ([]byte, error) {
	w.lock.Lock()
	defer w.lock.Unlock()
	return w.ldb.Get(nextSeedDBKey, nil)
}

func (w *Wallet) putSeed(seed []byte) error {
	return w.ldb.Put(seedDBKey, seed, nil)
}

func (w *Wallet) putNextSeed(seed []byte) error {
	return w.ldb.Put(nextSeedDBKey, seed, nil)
}
