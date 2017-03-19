package kip

import (
	"errors"
	"reflect"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

type Item struct { // or interface{} ?????
	Dao     *Dao
	Value   interface{}
	saved   bool
	updated bool
	updates *bson.M
}

func (i *Item) Save() error {

	if !i.saved {
		insertErr := i.Dao.Insert(i)
		if nil == insertErr {
			i.saved = true
			i.updated = true
		}

		return insertErr
	}

	if !i.updated {
		updateErr := i.Dao.update(i.GetId(), i.updates)
		if nil == updateErr {
			i.updated = true
		}

		return updateErr
	}

	return errors.New("already saved")
}

func (i *Item) Patch(p *Patch) error {

	if "set" == p.Operation {
		if nil == i.updates {
			i.updates = &bson.M{}
		}

		u := *i.updates

		if _, exists := u["$set"]; !exists {
			u["$set"] = bson.M{}
		}

		// TODO: check value type with type
		// TODO: check mapping bson/json field name

		u["$set"].(bson.M)[p.Key] = p.Value
		i.updated = false
		return nil
	} else if "add_to_set" == p.Operation {
		if nil == i.updates {
			i.updates = &bson.M{}
		}

		u := *i.updates

		if _, exists := u["$addToSet"]; !exists {
			u["$addToSet"] = bson.M{}
		}

		c := u["$addToSet"].(bson.M)

		if _, exists := c[p.Key]; !exists {
			c[p.Key] = bson.M{
				"$each": []interface{}{},
			}
		}

		f := c[p.Key].(bson.M)["$each"].([]interface{})

		c[p.Key].(bson.M)["$each"] = append(f, p.Value)

		i.updated = false
		return nil
	}

	return errors.New("invalid operation")
}

func (i *Item) GetId() interface{} {
	return get_id(i.Value)
}

func get_id(item interface{}) interface{} {
	t := reflect.ValueOf(item)

	// Follow pointers
	for reflect.Ptr == t.Kind() {
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.Struct:
		// Traverse all fields and search for tag bson:"_id"
		n := t.NumField()
		for i := 0; i < n; i++ {
			field := t.Type().Field(i)
			if word_in_string("_id", field.Tag.Get("bson")) {
				return t.Field(i).Interface()
			}
		}
		// Fallback: search CI for fieldnames 'id'
		for i := 0; i < n; i++ {
			field := t.Type().Field(i)
			if "id" == strings.ToLower(field.Name) {
				return t.Field(i).Interface()
			}
		}
		return nil

	case reflect.Map:
		return t.MapIndex(reflect.ValueOf("_id")).Interface()
	}

	return nil
}

func word_in_string(w string, s string) bool {
	for _, v := range strings.Split(s, " ") {
		if w == v {
			return true
		}
	}
	return false
}

func (i *Item) Delete() error {
	d := i.Dao
	collection := d.Collection

	return d.Database.C(collection.Name).RemoveId(i.GetId())
}
