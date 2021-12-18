package database

import (
	"encoding/json"
	"fmt"

	"github.com/Academics-and-Career-Council/Stargazer.git/internal/models"

	"github.com/dgraph-io/badger/v3"
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
 func WriteToBadger(db *badger.DB, key []byte, body []byte) {
	db.Update(func(txn *badger.Txn) error {
		err := txn.SetEntry(badger.NewEntry(key, body))//to add Withttl
		//fmt.Println("sent to badgerDB", d.Body)
		return err
	})
 }