package kip

import (
	. "gopkg.in/check.v1"
	"gopkg.in/mgo.v2/bson"
)

func (w *World) Test_InstanceFindById_Ok(c *C) {

	// Prepare
	john := w.Users.Create()
	w.Users.Insert(john)

	// Run
	u := w.Users.FindById(john.GetId())

	// Check
	c.Assert(u.Value, DeepEquals, john.Value)
}

func (w *World) Test_InstanceFindById_NotFound(c *C) {

	u := w.Users.FindById("invented id")

	// Check
	c.Assert(u, IsNil)
}

func (w *World) Test_InstanceFindOne_Ok(c *C) {

	// Prepare
	john := w.Users.Create()
	john.Value.(*User).Name = "John Snow"
	w.Users.Insert(john)

	// Run
	u := w.Users.FindOne(bson.M{"name": "John Snow"})

	// Check
	c.Assert(u.Value, DeepEquals, john.Value)
}

func (w *World) Test_InstanceFindOne_Fail(c *C) {

	// Run
	u := w.Users.FindOne(bson.M{"name": "John Snow"})

	// Check
	c.Assert(u, IsNil)
}
