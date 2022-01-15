package persistentstore

import (
	"fmt"
	"github.com/dgraph-io/badger"
	"github.com/libp2p/go-libp2p-core/peer"
	"log"
)
// SetUsername
// 用户设置用户名，并且存储在本地，用户再次接入网络后会直接从本地存储读取用户信息

const dbPath = "/tmp/pworld_chat"

func SaveUsername(id peer.ID, uname string) error {
	if !PassedDetect(uname) {
		log.Fatal("用户名重复，请更换")
	}
	db, err := badger.Open(badger.DefaultOptions(dbPath))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(id.Pretty()), []byte(uname))

		return err
	})
	if err != nil {
		log.Fatal(err)
	}

	return err
}

func GetUserNameFromLocalStore(pid string) string {
	db, err := badger.Open(badger.DefaultOptions(dbPath))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var uname string
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(pid))
		if err != nil {
			log.Println(err)
			return err
		}

		var valCopy []byte
		err = item.Value(func(val []byte) error {
			fmt.Printf("The username is: %s\n", val)
			valCopy = append([]byte{}, val...)
			uname = string(valCopy)
			return nil
		})

		return err
	})
	if err != nil {
		log.Fatal(err)
	}


	return uname
}

// PassedDetect 设置用户名的时候为了确保名字不重复，可以向某个提供名字检测服务的节点发出检测请求
// 这样就确保了每个加入网络的用户的名字是不重复的，通过用户名去识别用户才有意义
// 为什么不用pid识别呢，pid是对应设备的，确切说是对应链接的，比如我两台电脑同时加入聊天室，
// 但是我期望的是两个设备都对应的我的用户名
// 所以uname 和 pid 的关系是1对多的关系
func PassedDetect(uname string) bool {

	return true
}


func SaveFavorRoom(uname string, roomname string) {

}

func GetFavorRooms() []string{
	rooms := make([]string, 0)

	return rooms
}


func SavePID(id peer.ID, uname string) error {
	if !PassedDetect(uname) {
		log.Fatal("用户名重复，请更换")
	}
	db, err := badger.Open(badger.DefaultOptions(dbPath))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(uname), []byte(id))

		return err
	})
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func GetPIDFromLocalStorage(uname string) string {
	db, err := badger.Open(badger.DefaultOptions(dbPath))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var pid string

	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(uname))
		if err != nil {
			log.Println(err)
			return err
		}

		var valCopy []byte
		err = item.Value(func(val []byte) error {
			fmt.Printf("The pid is: %s\n", val)
			valCopy = append([]byte{}, val...)
			pid = string(valCopy)
			return nil
		})

		return err
	})
	if err != nil {
		log.Fatal(err)
	}


	return pid
}