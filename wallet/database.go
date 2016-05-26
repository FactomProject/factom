// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package wallet

import (
	"crypto/rand"
	"fmt"
	"sync"
	
	"github.com/FactomProject/factom"
	"github.com/FactomProject/goleveldb/leveldb"
	"github.com/FactomProject/goleveldb/leveldb/opt"
)

var (
	fcDBPrefix    = []byte("Factoids")
	ecDBPrefix    = []byte("Entry Credits")
	seedDBKey     = []byte("DB Seed")
	nextSeedDBKey = []byte("Next Seed")
)

type WalletDB struct {
	lock sync.RWMutex
	ldb *leveldb.DB
}

func NewWalletDB(path string) (*WalletDB, error) {
	o := &opt.Options{ ErrorIfExist: true }
	wdb := new(WalletDB)
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

func OpenFile(path string) (*WalletDB, error) {
	o := &opt.Options{ ErrorIfMissing: true }
	wdb := new(WalletDB)
	if l, err := leveldb.OpenFile(path, o); err != nil {
		wdb.ldb = l
		return nil, err
	}
	// TODO - validate database
	// ? - check if db is corrrupt and recover
	return wdb, nil
}

func (w *WalletDB) Close() error {
	w.lock.Lock()
	defer w.lock.Unlock()

	return w.ldb.Close()
}

func (w *WalletDB) GetECAddress(a string) (*factom.ECAddress, error) {
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

func (w *WalletDB) GetFCTAddress(a string) (*factom.FactoidAddress, error) {
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



func (w *WalletDB) PutECAddress(e *factom.ECAddress) error {
	key := append(ecDBPrefix, e.PubString()...)
	return w.ldb.Put(key, []byte(e.SecString()), nil)
}

func (w *WalletDB) PutFCTAddress(f *factom.FactoidAddress) error {
	key := append(fcDBPrefix, f.PubString()...)
	return w.ldb.Put(key, []byte(f.SecString()), nil)
}
