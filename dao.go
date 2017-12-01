package kip

import (
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

func (i *Dao) Insert(o *Item) error {
	// TODO: Check if already inserted?

	err := i.Database.C(i.Collection.Name).Insert(o.Value)

	// TODO: Update inserted field?

	return err
}

func (d *Dao) update(selector interface{}, update interface{}) error {
	return d.Database.C(d.Collection.Name).Update(selector, update)
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
	err := i.Database.C(collection).Find(query).One(item.Value)

	if nil != err {
		return nil, err
	}

	item.saved = true
	item.updated = true
	return item, nil
}

func (d *Dao) Find(query bson.M) *Query {
	collection := d.Collection.Name
	return &Query{
		dao:       d,
		mgo_query: d.Database.C(collection).Find(query),
	}
}

// Delete will remove all items that match with the query
func (i *Dao) Delete(query bson.M) (n int, err error) {

	name := i.Collection.Name
	c := i.Database.C(name)

	info, err := c.RemoveAll(query)

	return info.Removed, err
}
