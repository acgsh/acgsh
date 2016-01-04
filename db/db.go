package db

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"strconv"
	"time"
)

const max_post_id = 99999999

var db *bolt.DB

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
func ui64tob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

func Init() {
	var err error
	log.Println("Open DB: acgsh_bolt.db")
	db, err = bolt.Open("acgsh_bolt.db", 0600, &bolt.Options{Timeout: 1 * time.Second})

	if err != nil {
		log.Fatal(err)
	}

	//Initialise all buckets
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Posts"))
		if err != nil {
			log.Printf("DB create Posts bucket: %s", err)
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte("Publishers"))
		if err != nil {
			log.Printf("DB create Publisher bucket: %s", err)
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte("PublishersReplyPosts"))
		if err != nil {
			log.Printf("DB create PublishersReplyPosts bucket: %s", err)
			return err
		}
		return nil
	})

	// Check stats
	/*	go func() {
		// Grab the initial stats.
		prev := db.Stats()

		for {
			// Wait for 10s.
			time.Sleep(10 * time.Second)

			// Grab the current stats and diff them.
			stats := db.Stats()
			diff := stats.Sub(&prev)

			// Encode stats to JSON and print to STDERR.
			json.NewEncoder(os.Stderr).Encode(diff)

			// Save stats for the next loop.
			prev = stats
		}
	}()*/

	log.Println("DB initialised successfully.")
}

func Close() {

	db.Close()
	log.Println("DB closed.")

}

func AddPublishersIfNotExist(names *[]string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("Publishers"))
		if bucket == nil {
			return fmt.Errorf("Bucket Publisher not found!")
		}

		for _, name := range *names {
			//fmt.Printf("Value [%d] is [%s]\n", index, name)
			v := bucket.Get([]byte(name))
			if v == nil {
				err := bucket.Put([]byte(name), []byte{})
				if err != nil {
					return err
				}
			}
		}

		return nil
	})

	return err
}

type SyncData struct {
	Max    int64 `json:"max"`
	Latest int64 `json:"latest"`
	Since  int64 `json:"since"`
}

func newSyncData() *SyncData {
	return &SyncData{
		Max:    max_post_id,
		Latest: -1,
		Since:  -1,
	}
}

func GetPublishers() (map[string]SyncData, error) {
	publishers := make(map[string]SyncData)
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("Publishers"))
		if bucket == nil {
			return fmt.Errorf("Bucket Publishers not found!")
		}

		bucket.ForEach(func(k, v []byte) error {
			//fmt.Printf("key=%s, value=%s\n", k, v)
			var sd *SyncData
			sd = newSyncData()
			if len(v) == 0 {
				publishers[string(k)] = *sd
				return nil
			}

			err := json.Unmarshal(v, sd)
			if err != nil {
				log.Println(err)
				log.Printf("Error: DB GetPublishers Unmarshal: %s", v)
				sd = newSyncData()
				publishers[string(k)] = *sd
				return nil
			}

			publishers[string(k)] = *sd
			return nil
		})

		return nil
	})
	return publishers, err
}

func UpdatePublishers(publishers *map[string]SyncData) error {
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("Publishers"))
		if bucket == nil {
			return fmt.Errorf("Bucket Publisher not found!")
		}

		for name, sd := range *publishers {
			//fmt.Printf("Value [%d] is [%s]\n", index, name)
			jsonData, err := json.Marshal(sd)
			if err != nil {
				log.Println(err)
				log.Printf("Error: DB UpdatePublishers Marshal: %+v\n", sd)
				continue
			}
			//log.Printf("%s\n", jsonData)

			err = bucket.Put([]byte(name), jsonData)
			if err != nil {
				log.Println(err)
				log.Printf("Error: DB UpdatePublishers bucket.Put: %s : %s\n", name, jsonData)
				continue
			}
		}

		return nil
	})

	return err
}

func DeletePublishers(names *[]string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		/*		bucket, err := tx.CreateBucketIfNotExists("fsef")
				if err != nil {
					return err
				}

				err = bucket.Put(key, value)
				if err != nil {
					return err
				}*/
		return nil
	})

	return err
}

func AddPosts(posts *ShPosts) error {
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("Posts"))
		if bucket == nil {
			return fmt.Errorf("Bucket Posts not found!")
		}

		for _, sp := range *posts {
			//fmt.Printf("Value [%d] is [%s]\n", index, name)
			jsonData, err := json.Marshal(sp)
			if err != nil {
				log.Println(err)
				log.Printf("Error: DB addPosts Marshal: %+v\n", sp)
				continue
			}
			//			log.Printf("%s\n", jsonData)

			bs := ui64tob(sp.Time)
			bs = append(bs, ":"...)
			bs = append(bs, sp.N...)
			bs = append(bs, ":"...)
			bs = append(bs, strconv.FormatInt(sp.K, 10)...)
			err = bucket.Put(bs, jsonData)
			if err != nil {
				log.Println(err)
				log.Printf("Error: DB addPosts bucket.Put: %s : %s\n", bs, jsonData)
				continue
			}
			//			fmt.Printf("Value [%s] is [%s]\n", bs, jsonData)
		}

		return nil
	})

	return err
}

// todo: Check if reply post was reply to acgsh post in DB.
func AddPublishersReplyPosts(posts *ShPubReplyPosts) error {
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("PublishersReplyPosts"))
		if bucket == nil {
			return fmt.Errorf("Bucket PublishersReplyPosts not found!")
		}

		for _, sp := range *posts {
			jsonData, err := json.Marshal(sp)
			if err != nil {
				log.Println(err)
				log.Printf("Error: DB AddPublishersReplyPosts Marshal: %+v\n", sp)
				continue
			}
			bs := []byte(sp.N)
			bs = append(bs, ":"...)
			bs = append(bs, strconv.FormatInt(sp.K, 10)...)

			newData := []byte{}
			v := bucket.Get(bs)
			if v == nil {
				newData = append(newData, "["...)
			} else if len(v) >= 2 {
				newData = append(newData, v[:(len(v)-1)]...)
				newData = append(newData, ","...)
			} else {
				newData = append(newData, "["...)
			}
			newData = append(newData, jsonData...)
			newData = append(newData, "]"...)

			err = bucket.Put(bs, newData)
			if err != nil {
				log.Println(err)
				log.Printf("Error: DB AddPublishersReplyPosts bucket.Put: %s : %s\n", bs, newData)
				continue
			}
		}

		return nil
	})

	return err
}

func GetPosts(idx, n uint) ([]byte, error) {
	if n < 1 {
		return []byte{}, fmt.Errorf("Invalid n")
	}
	buf := []byte("[")
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("Posts"))
		if bucket == nil {
			return fmt.Errorf("Bucket Posts not found!")
		}

		cur := bucket.Cursor()

		comma := []byte(",")
		i := uint(0)
		for k, v := cur.Last(); k != nil; k, v = cur.Prev() {
			if i >= idx {
				if i != idx {
					buf = append(buf, comma...)
				}
				buf = append(buf, v...)
			}
			i = i + 1
			if i >= idx+n {
				break
			}
		}

		return nil
	})
	buf = append(buf, "]"...)
	return buf, err
}

func GetPublishersReplyPosts(name, k string) ([]byte, error) {
	var data []byte
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("PublishersReplyPosts"))
		if bucket == nil {
			return fmt.Errorf("Bucket PublishersReplyPosts not found!")
		}

		bs := []byte(name)
		bs = append(bs, ":"...)
		bs = append(bs, k...)
		value := bucket.Get(bs)
		if value == nil {
			value = []byte("[]")
		}

		data = make([]byte, len(value))
		copy(data, value)

		return nil
	})

	return data, err
}
