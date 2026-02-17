package db

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"time"

	bolt "go.etcd.io/bbolt"
)

const (
	defaultDBPath  = "my.db"
	ideasBucketKey = "ideas"
	metaBucketKey  = "meta"

	legacyMigratedFlag = "legacy_migrated"
)

var (
	db             *bolt.DB
	ErrIdeaNotFound = errors.New("idea not found")
)

type Idea struct {
	Title      string `json:"title"`
	Raw        string `json:"raw"`
	Elaborated string `json:"elaborated"`
	CreatedAt  int64  `json:"created_at"`
}

type StoredIdea struct {
	ID   uint64
	Idea Idea
}

func InitDB() error {
	var err error
	db, err = bolt.Open(defaultDBPath, 0600, nil)
	if err != nil {
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		ideasBucket, createErr := tx.CreateBucketIfNotExists([]byte(ideasBucketKey))
		if createErr != nil {
			return createErr
		}

		metaBucket, metaErr := tx.CreateBucketIfNotExists([]byte(metaBucketKey))
		if metaErr != nil {
			return metaErr
		}

		if string(metaBucket.Get([]byte(legacyMigratedFlag))) == "1" {
			return nil
		}

		if err := migrateLegacyBuckets(tx, ideasBucket); err != nil {
			return err
		}

		return metaBucket.Put([]byte(legacyMigratedFlag), []byte("1"))
	})
}

func CloseDB() {
	if db != nil {
		_ = db.Close()
		db = nil
	}
}

func PutIdea(idea Idea) (uint64, error) {
	if db == nil {
		return 0, errors.New("db is not initialized")
	}
	if idea.CreatedAt == 0 {
		idea.CreatedAt = time.Now().Unix()
	}

	var ideaID uint64
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ideasBucketKey))
		if b == nil {
			return errors.New("ideas bucket not found")
		}

		nextID, seqErr := b.NextSequence()
		if seqErr != nil {
			return seqErr
		}
		payload, marshalErr := json.Marshal(idea)
		if marshalErr != nil {
			return marshalErr
		}

		ideaID = nextID
		return b.Put(u64ToBytes(ideaID), payload)
	})

	return ideaID, err
}

func GetIdea(id uint64) (StoredIdea, error) {
	if id == 0 {
		return StoredIdea{}, ErrIdeaNotFound
	}
	if db == nil {
		return StoredIdea{}, errors.New("db is not initialized")
	}

	var (
		stored StoredIdea
		found  bool
	)
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ideasBucketKey))
		if b == nil {
			return errors.New("ideas bucket not found")
		}

		raw := b.Get(u64ToBytes(id))
		if raw == nil {
			return nil
		}

		idea, decodeErr := decodeIdea(raw)
		if decodeErr != nil {
			return decodeErr
		}

		stored = StoredIdea{ID: id, Idea: idea}
		found = true
		return nil
	})
	if err != nil {
		return StoredIdea{}, err
	}
	if !found {
		return StoredIdea{}, ErrIdeaNotFound
	}

	return stored, nil
}

func UpdateIdea(id uint64, next Idea) error {
	if id == 0 {
		return ErrIdeaNotFound
	}
	if db == nil {
		return errors.New("db is not initialized")
	}

	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ideasBucketKey))
		if b == nil {
			return errors.New("ideas bucket not found")
		}

		key := u64ToBytes(id)
		existingRaw := b.Get(key)
		if existingRaw == nil {
			return ErrIdeaNotFound
		}

		existingIdea, decodeErr := decodeIdea(existingRaw)
		if decodeErr != nil {
			return decodeErr
		}

		if next.CreatedAt == 0 {
			next.CreatedAt = existingIdea.CreatedAt
		}
		if next.Elaborated == "" && existingIdea.Elaborated != "" {
			next.Elaborated = existingIdea.Elaborated
		}

		payload, marshalErr := json.Marshal(next)
		if marshalErr != nil {
			return marshalErr
		}

		return b.Put(key, payload)
	})
}

func DeleteIdea(id uint64) error {
	if id == 0 {
		return ErrIdeaNotFound
	}
	if db == nil {
		return errors.New("db is not initialized")
	}

	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ideasBucketKey))
		if b == nil {
			return errors.New("ideas bucket not found")
		}

		key := u64ToBytes(id)
		if b.Get(key) == nil {
			return ErrIdeaNotFound
		}
		return b.Delete(key)
	})
}

func ListIdeas() ([]StoredIdea, error) {
	if db == nil {
		return nil, errors.New("db is not initialized")
	}

	out := make([]StoredIdea, 0)
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ideasBucketKey))
		if b == nil {
			return errors.New("ideas bucket not found")
		}

		return b.ForEach(func(k, v []byte) error {
			idea, decodeErr := decodeIdea(v)
			if decodeErr != nil {
				return fmt.Errorf("decode idea %d: %w", bytesToU64(k), decodeErr)
			}
			out = append(out, StoredIdea{
				ID:   bytesToU64(k),
				Idea: idea,
			})
			return nil
		})
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].Idea.CreatedAt == out[j].Idea.CreatedAt {
			return out[i].ID > out[j].ID
		}
		return out[i].Idea.CreatedAt > out[j].Idea.CreatedAt
	})

	return out, nil
}

func decodeIdea(raw []byte) (Idea, error) {
	var idea Idea
	if err := json.Unmarshal(raw, &idea); err != nil {
		return Idea{}, err
	}
	return idea, nil
}

func u64ToBytes(v uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, v)
	return buf
}

func bytesToU64(b []byte) uint64 {
	if len(b) < 8 {
		return 0
	}
	return binary.BigEndian.Uint64(b)
}

func migrateLegacyBuckets(tx *bolt.Tx, ideasBucket *bolt.Bucket) error {
	return tx.ForEach(func(name []byte, b *bolt.Bucket) error {
		bucketName := string(name)
		if bucketName == ideasBucketKey || bucketName == metaBucketKey {
			return nil
		}

		return b.ForEach(func(k, v []byte) error {
			idea := Idea{
				Title:      string(k),
				Raw:        string(v),
				Elaborated: "",
				CreatedAt:  time.Now().Unix(),
			}

			payload, err := json.Marshal(idea)
			if err != nil {
				return err
			}

			nextID, seqErr := ideasBucket.NextSequence()
			if seqErr != nil {
				return seqErr
			}

			return ideasBucket.Put(u64ToBytes(nextID), payload)
		})
	})
}
