package wallet

import (
	"fmt"
	"sort"
	"strings"

	"github.com/FactomProject/factom"

	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/primitives"
	"github.com/FactomProject/factomd/database/databaseOverlay"
	"github.com/FactomProject/factomd/database/hybridDB"
	"github.com/FactomProject/factomd/database/mapdb"
)

// Database keys and key prefixes
var (
	fcDBPrefix    = []byte("Factoids")
	ecDBPrefix    = []byte("Entry Credits")
	seedDBKey     = []byte("DB Seed")
	nextSeedDBKey = []byte("Next Seed")
)

type WalletDatabaseOverlay struct {
	dbo databaseOverlay.Overlay
}

func NewWalletOverlay(db interfaces.IDatabase) *WalletDatabaseOverlay {
	answer := new(WalletDatabaseOverlay)
	answer.dbo.DB = db
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
	db := hybridDB.NewBoltMapHybridDB(nil, boltPath)

	fmt.Println("Database started from: " + boltPath)
	return NewWalletOverlay(db), nil
}

func (db *WalletDatabaseOverlay) InsertDBSeed(dbSeed []byte) error {
	if dbSeed == nil {
		return nil
	}
	if len(dbSeed) != SeedLength {
		return fmt.Errorf("Provided Seed is the wrong length: %d", len(dbSeed))
	}
	data := new(primitives.ByteSlice)
	err := data.UnmarshalBinary(dbSeed)
	if err != nil {
		return err
	}

	batch := []interfaces.Record{}
	batch = append(batch, interfaces.Record{seedDBKey, seedDBKey, data})

	return db.dbo.PutInBatch(batch)
}

func (db *WalletDatabaseOverlay) GetDBSeed() ([]byte, error) {
	data, err := db.dbo.Get(seedDBKey, seedDBKey, new(primitives.ByteSlice))
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	return data.MarshalBinary()
}

func (db *WalletDatabaseOverlay) InsertNextDBSeed(dbSeed []byte) error {
	if dbSeed == nil {
		return nil
	}
	if len(dbSeed) != SeedLength {
		return fmt.Errorf("Provided Seed is the wrong length: %d", len(dbSeed))
	}
	data := new(primitives.ByteSlice)
	err := data.UnmarshalBinary(dbSeed)
	if err != nil {
		return err
	}

	batch := []interfaces.Record{}
	batch = append(batch, interfaces.Record{nextSeedDBKey, nextSeedDBKey, data})

	return db.dbo.PutInBatch(batch)
}

func (db *WalletDatabaseOverlay) GetNextDBSeed() ([]byte, error) {
	data, err := db.dbo.Get(nextSeedDBKey, nextSeedDBKey, new(primitives.ByteSlice))
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	return data.MarshalBinary()
}

func (db *WalletDatabaseOverlay) InsertECAddress(e *factom.ECAddress) error {
	if e == nil {
		return nil
	}

	batch := []interfaces.Record{}
	batch = append(batch, interfaces.Record{ecDBPrefix, []byte(e.PubString()), e})

	return db.dbo.PutInBatch(batch)
}

func (db *WalletDatabaseOverlay) GetECAddress(pubString string) (*factom.ECAddress, error) {
	data, err := db.dbo.Get(ecDBPrefix, []byte(pubString), new(factom.ECAddress))
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	return data.(*factom.ECAddress), nil
}

func (db *WalletDatabaseOverlay) GetAllECAddresses() ([]*factom.ECAddress, error) {
	list, err := db.dbo.FetchAllBlocksFromBucket(ecDBPrefix, new(ECA))
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

	return db.dbo.PutInBatch(batch)
}

func (db *WalletDatabaseOverlay) GetFCTAddress(str string) (*factom.FactoidAddress, error) {
	data, err := db.dbo.Get(fcDBPrefix, []byte(str), new(factom.FactoidAddress))
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, ErrNoSuchAddress
	}
	return data.(*factom.FactoidAddress), nil
}

func (db *WalletDatabaseOverlay) GetAllFCTAddresses() ([]*factom.FactoidAddress, error) {
	list, err := db.dbo.FetchAllBlocksFromBucket(fcDBPrefix, new(FA))
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
