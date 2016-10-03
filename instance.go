package kip

import "gopkg.in/mgo.v2/bson"

type Instance struct {
	Collection *Collection
	Database   *Database
}

/**
 * Create a new item for the existing collection
 */
func (i *Instance) Create() *Item {

	c := i.Collection.OnCreate
	if nil == c {
		panic("Mandatory callback `OnCreate` is needed for `" + i.Collection.Name + "`")
	}

	return &Item{
		Instance: i,
		Value:    c(),
	}
}

/**
 * FindById is a particular case of FindOne
 */
func (i *Instance) FindById(id interface{}) *Item {
	return i.FindOne(bson.M{"_id": id})
}

/**
 * Returned values:
 *  - *Item   -> All works
 *  - nil     -> Item not found
 *  - panic() -> Some kind of uncontrolled error happened
 */
func (i *Instance) FindOne(query bson.M) *Item {
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

func (i *Instance) Find(query bson.M) *Query {
	collection := i.Collection.Name
	return &Query{
		mgo_query: i.Database.C(collection).Find(query),
	}
}
