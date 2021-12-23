package database

import (
	"encoding/json"
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
					new.Batch = bID + 1//changes Batch ID to make it unique from previous batch
					temp, err := json.Marshal(new)
					checkHere(err)
					err = txn.Set(itr.Item().Key(), []byte(temp))
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
					
					if delBID == new.Batch {
						txn.Delete(itr.Item().Key())//deletes all data in badger from previous batch
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

func BulkWrite(db *badger.DB) {//write to mongo from badgerDB
	batchID := MongoClient.GetLastBatchID()//gets last batch ID from mongoDB

	s := gocron.NewScheduler(time.UTC)
	time.Sleep(20*time.Second)
	s.Every(1).Hour().Do(func ()  {	
		studList := GetFromBadger(db, batchID)//returns all students in badger 
		MongoClient.BulkWriteInSyslog(studList, db, batchID)//writes everything to mongo finally
		batchID = MongoClient.GetLastBatchID()//updates the batchID
	})
	s.StartAsync()
	s.StartBlocking()
}



func WriteToBadger(db *badger.DB, key []byte, body []byte) {
	db.Update(func(txn *badger.Txn) error {
		err := txn.SetEntry(badger.NewEntry(key, body))//object sent to badgerDB
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