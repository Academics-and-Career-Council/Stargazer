package database

import (
	"encoding/json"
	"fmt"
	//"time"

	"github.com/dgraph-io/badger/v3"
	//"github.com/go-co-op/gocron"
	"github.com/Academics-and-Career-Council/Stargazer.git/internal/models"
	//"github.com/Academics-and-Career-Council/Stargazer.git/internal/api"
)

func checkHere(err error) {
	if err != nil {
		panic(err)
	}
}

func GetFromBadger(db *badger.DB, bID int) []Models.Student {

	var studList []Models.Student

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
					var new Models.Student
					json.Unmarshal(p, &new)
					new.Batch = bID + 1
					fmt.Println(new)
					key := func(i int) []byte {
						return []byte(fmt.Sprintf("%d", i))
					}
					temp, err := json.Marshal(new)
					checkHere(err)
					err = txn.Set(key(new.ID), []byte(temp))
					checkHere(err)
					studList = append(studList, new)
					return nil
				})
				return nil
			})
			checkHere(err)
		}

		return nil
	})
	checkHere(err)

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
					var new Models.Student
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

func BulkWriteToBadger(db *badger.DB, batchID int, flag bool) {
		if !flag {
			flag = true
		} else {
			studList := GetFromBadger(db, batchID)
			fmt.Println(studList)
			MongoClient.BulkWriteInStudents(studList, db, batchID)
			batchID = MongoClient.GetLastBatchID()
			fmt.Println(batchID)
		}
}
// 	// s.StartAsync()
// 	// go API.GetSyslog()
// 	// Services.GetFromRabbitMQ(db)
// 	// //syslog.GetSyslog()
// 	// s.StartBlocking()
// }

func WriteToBadger(db *badger.DB, key []byte, body []byte) {
	db.Update(func(txn *badger.Txn) error {
		err := txn.SetEntry(badger.NewEntry(key, body)) //to add Withttl
		//fmt.Println("sent to badgerDB", d.Body)
		return err
	})
}

func OpenBadgerDB() *badger.DB {
	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		//fmt.Println("is it here?")
		panic(err)
	}

	defer db.Close()
	return db
}
