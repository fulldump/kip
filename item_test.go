package kip

import (
	. "gopkg.in/check.v1"
	"gopkg.in/mgo.v2/bson"
)

func (w *World) Test_Item_Delete(c *C) {

	john := w.Users.Create()
	john.Save()
	john.Delete()

	// Check
	item := &User{}
	err := w.Database.C(w.Users.Collection.Name).Find(bson.M{
		"_id": john.GetId(),
	}).One(item)

	c.Assert(err.Error(), Equals, "not found")

}

func (w *World) Test_Item_Save(c *C) {

	john := w.Users.Create()
	john.Save()

	// Check
	item := &User{}
	err := w.Database.C(w.Users.Collection.Name).Find(bson.M{
		"_id": john.GetId(),
	}).One(item)

	c.Assert(err, IsNil)

}

func (w *World) Test_Item_Save_Twice(c *C) {

	john := w.Users.Create()

	c.Assert(john.Save(), IsNil)
	c.Assert(john.Save(), NotNil)
}
