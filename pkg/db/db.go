package db

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/mrbttf/bridge-server/pkg/game"
)

const (
	sessionBucket = "session"
	playerBucket  = "player"
)

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
		_, err := tx.CreateBucketIfNotExists([]byte(sessionBucket))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(playerBucket))
		return err
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (db *DB) CreateSession(sessionName, hostPlayer string) (string, error) {
	session := game.New(sessionName, hostPlayer)
	err := db.SaveSession(session)
	return session.ID, err
}

func (db *DB) JoinSession(sessionID, playerName string) error {
	session, err := db.GetSession(sessionID)
	session.Players[playerName] = nil
	err = db.SaveSession(session)
	return err
}

func (db *DB) SaveSession(session *game.Session) error {
	return db.bolt.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(sessionBucket))

		encoded, err := json.Marshal(session)
		fmt.Println()
		fmt.Println(string(encoded))
		if err != nil {
			return err
		}
		return b.Put([]byte(session.ID), encoded)
	})
}

func (db *DB) GetSession(sessionID string) (*game.Session, error) {
	var session *game.Session
	err := db.bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(sessionBucket))
		v := b.Get([]byte(sessionID))
		if v == nil {
			return fmt.Errorf("Session doesn't exist: %s", sessionID)
		}
		return json.Unmarshal(v, &session)
	})
	return session, err
}

func (db DB) Close() error {
	return db.bolt.Close()
}
