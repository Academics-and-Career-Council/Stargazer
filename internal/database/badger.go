package database

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Academics-and-Career-Council/Stargazer.git/internal/models"
	"github.com/dgraph-io/badger/v3"
	"github.com/go-co-op/gocron"
	"github.com/spf13/viper"
)

func checkHere(err error) {
	if err != nil {
		panic(err)
	}
}

func GetFromBadger(db *badger.DB, bID int) []Models.Syslog {

	var studList []Models.Syslog
	flag := false

	var itr *badger.Iterator
	err := db.View(func(txn *badger.Txn) error {
		iopt := badger.DefaultIteratorOptions
		itr = txn.NewIterator(iopt)
		defer itr.Close()
		for itr.Rewind(); itr.Valid(); itr.Next() {
			err := db.Update(func(txn *badger.Txn) error {

				item, err := txn.Get(itr.Item().Key())

				checkHere(err)

				item.Value(func(val []byte) error {
					p := append([]byte{}, val...)
					var new Models.Syslog
					json.Unmarshal(p, &new)
					new.Batch = bID + 1
					key := func(i int) []byte {
						return []byte(fmt.Sprintf("%d", i))
					}
					temp, err := json.Marshal(new)
					checkHere(err)
					err = txn.Set(key(new.ID), []byte(temp))
					checkHere(err)
					sevLvl := ConvertSevirity(new.Severity)
					reqSeverity := ConvertSevirity(viper.GetString("severityCheck.severity"))
					allow := false
					if sevLvl <= reqSeverity {
						allow = true
					}
					
					if new.ServiceName != "" && allow {
						studList = append(studList, new)
						flag = true
					}
					return nil
				})
				return nil
			})
			checkHere(err)
		}

		return nil
	})
	checkHere(err)
	if !flag {
		return nil
	}
	log.Println("successfully read one batch from BadgerDB")
	return studList
}
func DeleteFromBadger(db *badger.DB, bID int) {

	var itr *badger.Iterator
	err := db.View(func(txn *badger.Txn) error {
		iopt := badger.DefaultIteratorOptions
		itr = txn.NewIterator(iopt)
		defer itr.Close()
		for itr.Rewind(); itr.Valid(); itr.Next() {
			err := db.Update(func(txn *badger.Txn) error {

				item, err := txn.Get(itr.Item().Key())

				checkHere(err)
				item.Value(func(val []byte) error {
					p := append([]byte{}, val...)
					var new Models.Syslog
					json.Unmarshal(p, &new)
					delBID := bID + 1
					key := func(i int) []byte {
						return []byte(fmt.Sprintf("%d", i))
					}
					if delBID == new.Batch {
						txn.Delete(key(new.ID))
					}

					return nil
				})
				return nil
			})
			checkHere(err)
		}

		return nil
	})
	checkHere(err)
	log.Println("Deleted previous batch")

}

func BulkWrite(db *badger.DB) {
	flag := false
	batchID := MongoClient.GetLastBatchID()

	s := gocron.NewScheduler(time.UTC)

	s.Every(30).Seconds().Do(func ()  {	
		if !flag {
			flag = true
		} else {
			studList := GetFromBadger(db, batchID)
			MongoClient.BulkWriteInSyslog(studList, db, batchID)
			batchID = MongoClient.GetLastBatchID()
		}
	})
	s.StartAsync()
	s.StartBlocking()
}



func WriteToBadger(db *badger.DB, key []byte, body []byte) {
	db.Update(func(txn *badger.Txn) error {
		err := txn.SetEntry(badger.NewEntry(key, body))
		if err != nil {
			panic(err)
		} else {
			log.Println("written to BadgerDB")
		}
		return err
	})
}

func OpenBadgerDB() *badger.DB {
	db, err := badger.Open(badger.DefaultOptions(viper.GetString("badger.fileLoc")))
	if err != nil {
		panic(err)
	}

	//defer db.Close()
	return db
}

func ConvertSevirity(sev string) int {
	var lvl int
	if sev == "emerg" {
		lvl = 0
	} else if sev == "alert" {
		lvl = 1
	} else if sev == "crit" {
		lvl = 2
	} else if sev == "err" {
		lvl = 3
	} else if sev == "warning" {
		lvl = 4
	} else if sev == "notice" {
		lvl = 5
	} else if sev == "info" {
		lvl = 6
	} else if sev == "debug" {
		lvl = 7
	}else {
		panic("invalid severity type")
	}
	return lvl
}