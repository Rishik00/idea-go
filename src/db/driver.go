package db

import (
	"errors"
	"fmt"

	bolt "go.etcd.io/bbolt" 
)

var db *bolt.DB

func InitDB() error {
	var err error
	db, err = bolt.Open("my.db", 0600, nil)
	if err != nil {
		return err
	}

	fmt.Println("DB initialized at my.db")
	return nil
}

func CheckBucket(bucketname string) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketname))
		if b == nil {
			return errors.New("bucket does not exist")
		}
		return nil
	})
}

func AddBucket(bucketname string) error {
	return db.Update(func (tx *bolt.Tx) error{

		_, err := tx.CreateBucketIfNotExists([]byte(bucketname))
		if err != nil {
			return fmt.Errorf("failed to create bucket %q: %w", bucketname, err)
		}

		return nil
	}) 
}

func AddIdea(bucketname, idea, description string) (string, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketname))
		if b == nil {
			return fmt.Errorf("bucket %q does not exist", bucketname)
		}
		return b.Put([]byte(idea), []byte(description))
	})

	return "Done, added title and description to DB", err
}

func deleteIdea(bucket, idea string) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("bucket %q doesnt exist", bucket)
		}
		return b.Delete([]byte(idea))
	})
}

func IdeasPerBucket(bucket string) ([]string, error) {
	var bucketOfIdeas []string

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("bucket %q doesnt exist", bucket)
		}
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("key=%s, value=%s\n", k, v)
			bucketOfIdeas = append(bucketOfIdeas, string(k))
		}

		return nil    
	})	

	return bucketOfIdeas, err
}

func ShowExistingBuckets() ([]string, error) {
	var existingBuckets []string

	err := db.View(func(tx *bolt.Tx) error {
		tx.ForEach(func(name []byte, b *bolt.Bucket) error {

			existingBuckets = append(existingBuckets, string(name))
			return nil
		})
		return nil
	})

	return existingBuckets, err		
}

func DontJustShowThemAll(bucket string, num int) error {
	// TODO: implement a filtered or formatted idea viewer
	return nil
}

func CloseDB() {
	if db != nil {
		db.Close()
	}
}
