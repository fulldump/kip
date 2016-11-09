package kip

import (
	"reflect"
	"strings"
)

type Item struct { // or interface{} ?????
	Dao   *Dao
	Value interface{}
}

func (i *Item) Save() error {
	return i.Dao.Insert(i)
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
