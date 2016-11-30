package wallet

import (
	"bytes"
	"crypto/rand"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/FactomProject/factom"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/primitives"
	"github.com/FactomProject/factomd/database/databaseOverlay"
	"github.com/FactomProject/factomd/database/hybridDB"
	"github.com/FactomProject/factomd/database/mapdb"
	"github.com/FactomProject/go-bip32"
	"github.com/FactomProject/go-bip39"
)

// Database keys and key prefixes
var (
	fcDBPrefix = []byte("Factoids")
	ecDBPrefix = []byte("Entry Credits")
	seedDBKey  = []byte("DB Seed")
)

type WalletDatabaseOverlay struct {
	DBO databaseOverlay.Overlay
}

func NewWalletOverlay(db interfaces.IDatabase) *WalletDatabaseOverlay {
	answer := new(WalletDatabaseOverlay)
	answer.DBO.DB = db
	return answer
}

func NewMapDB() *WalletDatabaseOverlay {
	return NewWalletOverlay(new(mapdb.MapDB))
}

func NewLevelDB(ldbpath string) (*WalletDatabaseOverlay, error) {
	db, err := hybridDB.NewLevelMapHybridDB(ldbpath, false)
	if err != nil {
		fmt.Printf("err opening db: %v\n", err)
	}

	if db == nil {
		fmt.Println("Creating new db ...")
		db, err = hybridDB.NewLevelMapHybridDB(ldbpath, true)

		if err != nil {
			return nil, err
		}
	}
	fmt.Println("Database started from: " + ldbpath)
	return NewWalletOverlay(db), nil
}

func NewBoltDB(boltPath string) (*WalletDatabaseOverlay, error) {
	// check if the file exists or if it is a directory
	fileInfo, err := os.Stat(boltPath)
	if err == nil {
		if fileInfo.IsDir() {
			return nil, fmt.Errorf("The path %s is a directory.  Please specify a file name.", boltPath)
		}
	}

	// create the wallet directory if it doesn't already exist
	if os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(boltPath), 0700); err != nil {
			fmt.Printf("database error %s\n", err)
		}
	}

	if err != nil && !os.IsNotExist(err) { //some other error, besides the file not existing
		fmt.Printf("database error %s\n", err)
		return nil, err
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Could not use wallet file \"%s\"\n%v\n", boltPath, r)
			os.Exit(1)
		}
	}()
	db := hybridDB.NewBoltMapHybridDB(nil, boltPath)

	fmt.Println("Database started from: " + boltPath)
	return NewWalletOverlay(db), nil
}

type DBSeedBase struct {
	MnemonicSeed            string
	NextFactoidAddressIndex uint32
	NextECAddressIndex      uint32
}

type DBSeed struct {
	DBSeedBase
}

var _ interfaces.BinaryMarshallable = (*DBSeed)(nil)

func (e *DBSeed) MarshalBinary() ([]byte, error) {
	var data primitives.Buffer

	enc := gob.NewEncoder(&data)

	err := enc.Encode(e.DBSeedBase)
	if err != nil {
		return nil, err
	}
	return data.DeepCopyBytes(), nil
}

func (e *DBSeed) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	dec := gob.NewDecoder(primitives.NewBuffer(data))
	dbsb := DBSeedBase{}
	err = dec.Decode(&dbsb)
	if err != nil {
		return nil, err
	}
	e.DBSeedBase = dbsb
	return nil, nil
}

func (e *DBSeed) UnmarshalBinary(data []byte) (err error) {
	_, err = e.UnmarshalBinaryData(data)
	return
}

func (e *DBSeed) JSONByte() ([]byte, error) {
	return primitives.EncodeJSON(e)
}

func (e *DBSeed) JSONString() (string, error) {
	return primitives.EncodeJSONString(e)
}

func (e *DBSeed) JSONBuffer(b *bytes.Buffer) error {
	return primitives.EncodeJSONToBuffer(e, b)
}

func (e *DBSeed) String() string {
	str, _ := e.JSONString()
	return str
}

func (e *DBSeed) NextFCTAddress() (*factom.FactoidAddress, error) {
	add, err := factom.MakeBIP44FactoidAddress(
		e.MnemonicSeed,
		bip32.FirstHardenedChild,
		0,
		e.NextFactoidAddressIndex,
	)
	if err != nil {
		return nil, err
	}
	e.NextFactoidAddressIndex++
	return add, nil
}

func (e *DBSeed) NextECAddress() (*factom.ECAddress, error) {
	add, err := factom.MakeBIP44ECAddress(
		e.MnemonicSeed,
		bip32.FirstHardenedChild,
		0,
		e.NextECAddressIndex,
	)
	if err != nil {
		return nil, err
	}
	e.NextECAddressIndex++
	return add, nil
}

func NewRandomSeed() (*DBSeed, error) {
	seed := make([]byte, 16)
	if n, err := rand.Read(seed); err != nil {
		panic(err)
		return nil, err
	} else if n != 16 {
		return nil, fmt.Errorf("Wrong number of bytes read: %d", n)
	}

	mnemonic, err := bip39.NewMnemonic(seed)
	if err != nil {
		panic(err)
		return nil, err
	}

	dbSeed := new(DBSeed)
	dbSeed.MnemonicSeed = mnemonic

	return dbSeed, nil
}

func (db *WalletDatabaseOverlay) InsertDBSeed(seed *DBSeed) error {
	if seed == nil {
		return nil
	}

	batch := []interfaces.Record{}
	batch = append(batch, interfaces.Record{seedDBKey, seedDBKey, seed})

	return db.DBO.PutInBatch(batch)
}

func (db *WalletDatabaseOverlay) GetDBSeed() (*DBSeed, error) {
	data, err := db.DBO.Get(seedDBKey, seedDBKey, new(DBSeed))
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	return data.(*DBSeed), nil
}

func (db *WalletDatabaseOverlay) GetOrCreateDBSeed() (*DBSeed, error) {
	data, err := db.DBO.Get(seedDBKey, seedDBKey, new(DBSeed))
	if err != nil {
		return nil, err
	}
	if data == nil {
		seed, err := NewRandomSeed()
		if err != nil {
			return nil, err
		}
		err = db.InsertDBSeed(seed)
		if err != nil {
			return nil, err
		}
		return seed, nil
	}
	return data.(*DBSeed), nil
}

func (db *WalletDatabaseOverlay) GetNextECAddress() (*factom.ECAddress, error) {
	seed, err := db.GetOrCreateDBSeed()
	if err != nil {
		return nil, err
	}
	add, err := seed.NextECAddress()
	if err != nil {
		return nil, err
	}
	err = db.InsertDBSeed(seed)
	if err != nil {
		return nil, err
	}
	err = db.InsertECAddress(add)
	if err != nil {
		return nil, err
	}
	return add, nil
}

func (db *WalletDatabaseOverlay) GetNextFCTAddress() (*factom.FactoidAddress, error) {
	seed, err := db.GetOrCreateDBSeed()
	if err != nil {
		return nil, err
	}
	add, err := seed.NextFCTAddress()
	if err != nil {
		return nil, err
	}
	err = db.InsertDBSeed(seed)
	if err != nil {
		return nil, err
	}
	err = db.InsertFCTAddress(add)
	if err != nil {
		return nil, err
	}
	return add, nil
}

func (db *WalletDatabaseOverlay) InsertECAddress(e *factom.ECAddress) error {
	if e == nil {
		return nil
	}

	batch := []interfaces.Record{}
	batch = append(batch, interfaces.Record{ecDBPrefix, []byte(e.PubString()), e})

	return db.DBO.PutInBatch(batch)
}

func (db *WalletDatabaseOverlay) GetECAddress(pubString string) (*factom.ECAddress, error) {
	data, err := db.DBO.Get(ecDBPrefix, []byte(pubString), new(factom.ECAddress))
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, ErrNoSuchAddress
	}
	return data.(*factom.ECAddress), nil
}

func (db *WalletDatabaseOverlay) GetAllECAddresses() ([]*factom.ECAddress, error) {
	list, err := db.DBO.FetchAllBlocksFromBucket(ecDBPrefix, new(ECA))
	if err != nil {
		return nil, err
	}
	return toECList(list), nil
}

func toECList(source []interfaces.BinaryMarshallableAndCopyable) []*factom.ECAddress {
	answer := make([]*factom.ECAddress, len(source))
	for i, v := range source {
		answer[i] = v.(*ECA).ECAddress
	}
	sort.Sort(byECName(answer))
	return answer
}

type byECName []*factom.ECAddress

func (f byECName) Len() int {
	return len(f)
}
func (f byECName) Less(i, j int) bool {
	a := strings.Compare(f[i].PubString(), f[j].PubString())
	return a < 0
}
func (f byECName) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

func (db *WalletDatabaseOverlay) InsertFCTAddress(e *factom.FactoidAddress) error {
	if e == nil {
		return nil
	}

	batch := []interfaces.Record{}
	batch = append(batch, interfaces.Record{fcDBPrefix, []byte(e.String()), e})

	return db.DBO.PutInBatch(batch)
}

func (db *WalletDatabaseOverlay) GetFCTAddress(str string) (*factom.FactoidAddress, error) {
	data, err := db.DBO.Get(fcDBPrefix, []byte(str), new(factom.FactoidAddress))
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, ErrNoSuchAddress
	}
	return data.(*factom.FactoidAddress), nil
}

func (db *WalletDatabaseOverlay) GetAllFCTAddresses() ([]*factom.FactoidAddress, error) {
	list, err := db.DBO.FetchAllBlocksFromBucket(fcDBPrefix, new(FA))
	if err != nil {
		return nil, err
	}
	return toFList(list), nil
}

func toFList(source []interfaces.BinaryMarshallableAndCopyable) []*factom.FactoidAddress {
	answer := make([]*factom.FactoidAddress, len(source))
	for i, v := range source {
		answer[i] = v.(*FA).FactoidAddress
	}
	sort.Sort(byFName(answer))
	return answer
}

func (db *WalletDatabaseOverlay) RemoveAddress(pubString string) error {
	if len(pubString) == 0 {
		return nil
	}
	if pubString[:1] == "F" {
		data, err := db.DBO.Get(fcDBPrefix, []byte(pubString), new(factom.FactoidAddress))
		if err != nil {
			return err
		}
		if data == nil {
			return ErrNoSuchAddress
		}
		err = db.DBO.Delete(fcDBPrefix, []byte(pubString))
		if err == nil {
			err := db.DBO.Delete(fcDBPrefix, []byte(pubString)) //delete twice to flush the db file
			return err
		} else {
			return err
		}
	} else if pubString[:1] == "E" {
		data, err := db.DBO.Get(ecDBPrefix, []byte(pubString), new(factom.ECAddress))
		if err != nil {
			return err
		}
		if data == nil {
			return ErrNoSuchAddress
		}
		err = db.DBO.Delete(ecDBPrefix, []byte(pubString))
		if err == nil {
			err := db.DBO.Delete(ecDBPrefix, []byte(pubString)) //delete twice to flush the db file
			return err
		} else {
			return err
		}
	} else {
		return fmt.Errorf("Unknown address type")
	}

	return nil
}

type byFName []*factom.FactoidAddress

func (f byFName) Len() int {
	return len(f)
}
func (f byFName) Less(i, j int) bool {
	a := strings.Compare(f[i].String(), f[j].String())
	return a < 0
}
func (f byFName) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

type ECA struct {
	*factom.ECAddress
}

var _ interfaces.BinaryMarshallableAndCopyable = (*ECA)(nil)

func (t *ECA) New() interfaces.BinaryMarshallableAndCopyable {
	e := new(ECA)
	e.ECAddress = factom.NewECAddress()
	return e
}

type FA struct {
	*factom.FactoidAddress
}

var _ interfaces.BinaryMarshallableAndCopyable = (*FA)(nil)

func (t *FA) New() interfaces.BinaryMarshallableAndCopyable {
	e := new(FA)
	e.FactoidAddress = factom.NewFactoidAddress()
	return e
}
