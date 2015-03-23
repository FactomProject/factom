package factom

import (
	"fmt"
	"log"
//	"io/ioutil"
//	"net/http"
//	"strconv"
	"testing"

	"github.com/FactomProject/btcd/wire"
	"github.com/FactomProject/FactomCode/database"
	"github.com/FactomProject/FactomCode/database/ldb"
	"github.com/FactomProject/FactomCode/wsapi"
)

var (
	db      database.Db
	ldbpath = "/tmp/ldb9"
	MsgQ    = make(chan wire.Message, 100)
	server = "http://localhost:8088"
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

func TestGetDBlocks(t *testing.T) {
	wsapi.Start(db, MsgQ)
	
	dblocks, err := GetDBlocks(0, 5)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(dblocks)
}

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
