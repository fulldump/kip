package kip

import "gopkg.in/mgo.v2/bson"

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
		Dao:   i,
		Value: c(),
	}
}

func (i *Dao) Insert(o *Item) error {
	// TODO: Check if already inserted?

	err := i.Database.C(i.Collection.Name).Insert(o.Value)

	// TODO: Update inserted field?

	return err
}

/**
 * FindById is a particular case of FindOne
 */
func (i *Dao) FindById(id interface{}) *Item {
	return i.FindOne(bson.M{"_id": id})
}

/**
 * Returned values:
 *  - *Item   -> All works
 *  - nil     -> Item not found
 *  - panic() -> Some kind of uncontrolled error happened
 */
func (i *Dao) FindOne(query bson.M) *Item {
	item := i.Create()

	collection := i.Collection.Name
	err := i.Database.C(collection).Find(query).One(item.Value)

	if nil == err {
		return item
	}

	if "not found" == err.Error() {
		return nil
	}

	panic(err)
}

func (i *Dao) Find(query bson.M) *Query {
	collection := i.Collection.Name
	return &Query{
		mgo_query: i.Database.C(collection).Find(query),
	}
}

// Delete will remove all items that match with the query
func (i *Dao) Delete(query bson.M) (n int, err error) {

	name := i.Collection.Name
	c := i.Database.C(name)

	info, err := c.RemoveAll(query)

	return info.Removed, err
}
