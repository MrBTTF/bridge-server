package db

import "github.com/boltdb/bolt"

const SessionBucket = "session"

type DB struct {
	bolt *bolt.DB
}

func New(dbPath string) (*DB, error) {
	db := new(DB)
	var err error
	db.bolt, err = bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.bolt.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(SessionBucket))
		return err
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (db DB) Close() error {
	return db.bolt.Close()
}
