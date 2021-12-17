package badgerRabbitmq

import (
	"encoding/json"
	"fmt"

	"github.com/Academics-and-Career-Council/Stargazer.git/internal/structure"
	"github.com/dgraph-io/badger/v3"
	//"github.com/streadway/amqp"
)

func checkHere(err error) {
	if err != nil {
		panic(err)
	}
}

func GetFromBadger(db *badger.DB, bID int) []Structure.Student{


	var studList[] Structure.Student
		
		var itr *badger.Iterator
		err := db.View(func(txn *badger.Txn) error {
			iopt := badger.DefaultIteratorOptions
			itr = txn.NewIterator(iopt)
			defer itr.Close()
			for itr.Rewind(); itr.Valid(); itr.Next() {
				err := db.View(func(txn *badger.Txn) error {

					item, err := txn.Get(itr.Item().Key())

					checkHere(err)
					item.Value(func(val []byte) error {p := append([]byte{}, val...)
						var new Structure.Student
						json.Unmarshal(p, &new)
						new.Batch = bID+1
						fmt.Println(new)
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
func DeleteFromBadger (db *badger.DB, bID int) {
	

	var itr *badger.Iterator
		err := db.View(func(txn *badger.Txn) error {
			iopt := badger.DefaultIteratorOptions
			itr = txn.NewIterator(iopt)
			defer itr.Close()
			for itr.Rewind(); itr.Valid(); itr.Next() {
				err := db.Update(func(txn *badger.Txn) error {

					item, err := txn.Get(itr.Item().Key())

					checkHere(err)
					item.Value(func(val []byte) error {p := append([]byte{}, val...)
						var new Structure.Student
						json.Unmarshal(p, &new)
						delBID := bID + 1
						key := func(i int) []byte {
							return []byte(fmt.Sprintf("%d", i))
						}
						if delBID >= new.Batch {
							txn.Delete(key(new.ID))
						} 
						// db.Update(func(txn *badger.Txn) error {
						// 	err := txn.SetEntry(badger.NewEntry(key(stud.ID), []byte(d.Body)))//to add Withttl
						// 	return err
						// })
						//new.Batch = bID+1
						//fmt.Println(new)
						//studList = append(studList, new)
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