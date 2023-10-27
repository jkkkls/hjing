// Copyright 2017 guangbo. All rights reserved.

// watch kv数据库模块
package utils

import (
	"log"

	"github.com/boltdb/bolt"
)

type KfvInfo struct {
	Key   string
	Field string
	Value string
}

var boltDB *bolt.DB

func InitDb(file string) {
	var err error
	boltDB, err = bolt.Open(file, 0600, nil)
	if err != nil {
		log.Println(err)
	}
}

func GetDB() *bolt.DB {
	return boltDB
}

func GetAllKfv() []*KfvInfo {
	var kfvs []*KfvInfo
	if boltDB == nil {
		return kfvs
	}

	boltDB.View(func(tx *bolt.Tx) error {
		tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			b.ForEach(func(k, v []byte) error {
				kfvs = append(kfvs, &KfvInfo{string(name), string(k), string(v)})
				return nil
			})
			return nil
		})

		return nil
	})
	return kfvs
}

func AddKvs(key, field, value string) {
	if boltDB == nil {
		return
	}

	boltDB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(key))
		if err != nil {
			return err
		}
		b.Put([]byte(field), []byte(value))
		return nil
	})
}

func UpdateKvs(key, field, value string) {
	if boltDB == nil {
		return
	}

	AddKvs(key, field, value)
}

func DelKvs(key, field string) {
	if boltDB == nil {
		return
	}
	boltDB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(key))
		if err != nil {
			return err
		}

		b.Delete([]byte(field))

		return nil
	})
}

func DelKey(key string) {
	if boltDB == nil {
		return
	}
	boltDB.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket([]byte(key))

		return nil
	})
}
