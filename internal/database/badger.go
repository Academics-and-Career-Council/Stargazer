package database

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/spf13/viper"
	"github.com/go-co-op/gocron"
	"github.com/Academics-and-Career-Council/Stargazer.git/internal/models"
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
					if new.ServiceName != "" {
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
