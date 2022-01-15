package main

import (
	"fmt"
	"github.com/dgraph-io/badger"
	"log"
)

func main() {
	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte("answer"), []byte("42"))
		return err
	})
	if err != nil {
		log.Fatal(err)
	}

	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("answer"))
		if err != nil {
			log.Fatal(err)
		}

		var valNot, valCopy []byte
		err = item.Value(func(val []byte) error {
			fmt.Printf("The answer is: %s\n", val)
			valCopy = append([]byte{}, val...)
			valNot = val // don't do this
			return nil
		})

		return err
	})
}