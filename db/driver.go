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


func AddIdea(bucketname, idea, description string) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketname))
		if b == nil {
			return fmt.Errorf("bucket %q does not exist", bucketname)
		}
		return b.Put([]byte(idea), []byte(description))
	})
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

func ShowThemAll(bucket string, num int) error{
	return db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("bucket %q doesnt exist", bucket)
		}
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("key=%s, value=%s\n", k, v)
		}

		return nil
	})	
}

func DontJustShowThemAll() error {
	// TODO: implement a filtered or formatted idea viewer
	return nil
}

func CloseDB() {
	if db != nil {
		db.Close()
	}
}


func PrintExistingBuckets() error {
	return db.View(func(tx *bolt.Tx) error {
		tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			fmt.Println("Bucket: ", string(name))
			return nil
		})
		return nil
	})
}