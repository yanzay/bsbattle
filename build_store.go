package main

import (
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"
	"github.com/yanzay/log"
)

type BuildStore struct {
	db *bolt.DB
}

var bucketName = []byte("buildings")

func NewBuildStore(filename string) *BuildStore {
	db, err := bolt.Open(filename, 0600, &bolt.Options{Timeout: 5 * time.Second})
	if err != nil {
		log.Fatalf("can't open database %s: %q", filename, err)
	}
	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists(bucketName)
		return nil
	})
	return &BuildStore{
		db: db,
	}
}

func (bs *BuildStore) GetBuildings(name string) *Buildings {
	var buildBytes []byte
	bs.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		buildBytes = b.Get([]byte(name))
		return nil
	})
	builds := &Buildings{}
	err := json.Unmarshal(buildBytes, builds)
	if err != nil {
		log.Errorf("can't unmarshal buildings for %s: %q", name, err)
	}
	return builds
}

func (bs *BuildStore) SaveBuildings(name string, builds *Buildings) {
	buildBytes, err := json.Marshal(builds)
	if err != nil {
		log.Errorf("can't marshal buildings for %s: %q", name, err)
		return
	}
	bs.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		return b.Put([]byte(name), buildBytes)
	})
}
