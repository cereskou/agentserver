package db

import (
	"strings"

	"ditto.co.jp/agentserver/util"
	"github.com/tidwall/buntdb"
)

//Database -
type Database struct {
	DB  *buntdb.DB
	Map *util.AgentMap
}

//NewDatabase -
func NewDatabase() *Database {
	db, _ := buntdb.Open(":memory:")
	// db.CreateIndex("jobid", "*", buntdb.IndexString)

	return &Database{
		DB:  db,
		Map: nil,
	}
}

//CreateIndex - CreateIndex("jobs", "job.*", String)
func (d *Database) CreateIndex(index, pattern string) error {
	return d.DB.CreateIndex(index, pattern, buntdb.IndexString)
}

//Set -
func (d *Database) Set(key, val string) error {
	return d.DB.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(key, val, nil)
		return err
	})
}

//Get -
func (d *Database) Get(key string) (val string, err error) {
	err = d.DB.View(func(tx *buntdb.Tx) error {
		//値取得
		val, err = tx.Get(key)
		if err != nil {
			return err
		}

		return nil
	})
	return
}

//Fetch - list and delete
func (d *Database) Fetch(index string) (text string, err error) {
	slice := make([]string, 0)
	keys := make([]string, 0)
	err = d.DB.Update(func(tx *buntdb.Tx) error {
		//Iterating
		tx.Ascend(index, func(k, v string) bool {
			slice = append(slice, v)
			keys = append(keys, k)

			return true
		})
		//Delete
		for _, k := range keys {
			if _, err = tx.Delete(k); err != nil {
				return err
			}
		}
		return nil
	})

	text = "[" + strings.Join(slice, ",") + "]"

	return
}

//List -
func (d *Database) List(index string) (text string, err error) {
	slice := make([]string, 0)
	err = d.DB.View(func(tx *buntdb.Tx) error {
		tx.Ascend(index, func(k, v string) bool {
			slice = append(slice, v)

			return true
		})
		return nil
	})

	text = "[" + strings.Join(slice, ",") + "]"

	return
}

//Delete -
func (d *Database) Delete(key string) (err error) {
	return d.DB.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete((key))

		return err
	})
}

//Close -
func (d *Database) Close() error {
	return d.DB.Close()
}
