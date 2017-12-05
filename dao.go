package kip

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Dao is the combination of `Collection` definition plus a `Database`
type Dao struct {
	Collection *Collection
	Database   *Database
}

/**
 * Create a new item for the existing collection
 */
func (i *Dao) Create() *Item {

	c := i.Collection.OnCreate
	if nil == c {
		panic("Mandatory callback `OnCreate` is needed for `" + i.Collection.Name + "`")
	}

	return &Item{
		Dao:     i,
		Value:   c(),
		saved:   false,
		updated: false,
	}
}

func (d *Dao) Insert(o *Item) error {

	// TODO: Check if already inserted?

	db := d.Database.Clone()
	defer db.Close()

	err := db.C(d.Collection.Name).Insert(o.Value)

	// TODO: Update inserted field?

	return err
}

func (d *Dao) update(selector interface{}, update interface{}) error {

	db := d.Database.Clone()
	defer db.Close()

	return db.C(d.Collection.Name).Update(selector, update)
}

/**
 * FindById is a particular case of FindOne
 */
func (i *Dao) FindById(id interface{}) (*Item, error) {
	return i.FindOne(bson.M{"_id": id})
}

/**
 * Returned values:
 *  - *Item   -> All works
 *  - nil     -> Item not found
 *  - panic() -> Some kind of uncontrolled error happened
 */
func (i *Dao) FindOne(query bson.M) (*Item, error) {
	item := i.Create()

	collection := i.Collection.Name

	db := i.Database.Clone()
	defer db.Close()

	err := db.C(collection).Find(query).One(item.Value)

	if mgo.ErrNotFound == err {
		return nil, nil
	}

	if nil != err {
		return nil, err
	}

	item.saved = true
	item.updated = true
	return item, nil
}

func (d *Dao) Find(query interface{}) *Query {
	return &Query{
		dao:      d,
		selector: query,
	}
}

// Delete will remove all items that match with the query
func (d *Dao) Delete(query bson.M) (n int, err error) {

	db := d.Database.Clone()
	defer db.Close()

	info, err := db.C(d.Collection.Name).RemoveAll(query)

	return info.Removed, err
}
