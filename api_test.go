package factom

import (
	"encoding/json"
	//	"fmt"
	"io"
	"log"
	"strings"
	//	"io/ioutil"
	//	"net/http"
	//	"strconv"
	"testing"

	"github.com/FactomProject/FactomCode/database"
	"github.com/FactomProject/FactomCode/database/ldb"
	"github.com/FactomProject/btcd/wire"
	//	"github.com/FactomProject/FactomCode/wsapi"
)

/*******************************************************************************
 *
 *
 *      This test hangs if you are running FactomDB.  It is locked out
 *      by LevelDB!
 *
 *
 *
 *******************************************************************************/
var (
	db      database.Db
	ldbpath = "/tmp/ldb9"
	MsgQ    = make(chan wire.Message, 100)
)

func initDB() {
	var err error
	db, err = ldb.OpenLevelDB(ldbpath, false)

	if err != nil {
		log.Printf("err opening db: %v\n", err)
	}

	if db == nil {
		log.Println("Creating new db ...")
		db, err = ldb.OpenLevelDB(ldbpath, true)

		if err != nil {
			panic(err)
		}
	}
	log.Println("Database started from: " + ldbpath)
}

func TestDecode(t *testing.T) {
	const stream = `{"Header":{"BlockID":0,"PrevBlockHash":"0000000000000000000000000000000000000000000000000000000000000000","MerkleRoot":"29a960e1e98fe3881cbe96b18498fbab0bdda5fc3c5e13fc465a8d2ee33e2b1e","Version":1,"Timestamp":1424298447,"BatchFlag":0,"EntryCount":2},"DBEntries":[{"MerkleRoot":"2d8fc252e8ce40ee7ff0396621f69854e3f058fe640533510b81b89e9b68408d","ChainID":"f4f614fd9b59fe26827137937d401e2b82125c4eb48f966e6e8d30f187184cb0"},{"ChainID":"0100000000000000000000000000000000000000000000000000000000000000","MerkleRoot":"2cc92b09333c7d6172939be031332a12ca5a47b4716f4e8fcbfd69435a394e6f"}]}
{"Header":{"PrevBlockHash":"74c052d99050a334d35e6cbf196e2242921140308e12f500bb73298622f7395d","MerkleRoot":"2d7cb0911ede2948eec478d18e59baa50b5ccb1fa225a8c71b056f77bdb0df6b","Version":1,"Timestamp":1424298507,"BatchFlag":0,"EntryCount":2,"BlockID":1},"DBEntries":[{"MerkleRoot":"0538abefba2bb9c96b75698c1d18a2c32e0fed88bc1e1deac8901963697dbd69","ChainID":"f4f614fd9b59fe26827137937d401e2b82125c4eb48f966e6e8d30f187184cb0"},{"MerkleRoot":"8ff0698ebcb2d034d52b583f8b3646d4a9e992867fbe77ebbfea2d8c0d0d8547","ChainID":"0100000000000000000000000000000000000000000000000000000000000000"}]}`

	dblocks := make([]DBlock, 0)

	dec := json.NewDecoder(strings.NewReader(stream))
	for i := 0; i < 10; i++ {
		var dblock DBlock
		if err := dec.Decode(&dblock); err == io.EOF {
			break
		} else if err != nil {
			t.Errorf(err.Error())
			return
		}
		dblocks = append(dblocks, dblock)
	}
}

//func TestGetDBlocks(t *testing.T) {
//	server = "localhost:8088"
//	wsapi.Start(db, MsgQ)
//	defer wsapi.Stop()
//
//	dblocks, err := GetDBlocks(0, 5)
//	if err != nil {
//		t.Errorf(err.Error())
//	}
//	fmt.Println(dblocks)
//}

//func TestGetBlockHeight(t *testing.T) {
//	var server = "http://localhost:8088"
//	wsapi.Start(db, MsgQ)
//
//	resp, err := http.Get(server + "/v1/blockheight")
//	if err != nil {
//		t.Errorf(err.Error())
//	}
//	defer resp.Body.Close()
//	p, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		t.Errorf(err.Error())
//	}
//
//	fmt.Println(strconv.Atoi(string(p)))
//}
