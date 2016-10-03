#Kip

[![Build Status](https://travis-ci.org/fulldump/kip.svg?branch=master)](https://travis-ci.org/fulldump/kip)

Kip is a Object wrapper for MongoDB.

# How to use

## Define

Basic usage:

```go
type User struct {
	Name  string `bson:"name"`
	Age   int    `bson:"age"`
	Email string `bson:"email"`
}

kip.Define(&Collection{
	Name: "Users",
	OnCreate: func() interface{} {
		return &User{
			Id:   bson.NewObjectId(),
			Name: "default name",
		}
	},
})
```

Define indexes:

```go

kip.Define(&Collection{
	Name: "Users",
	OnCreate: func() interface{} {
		return &User{
			Id:   bson.NewObjectId(),
			Name: "default name",
		}
	},
}).EnsureIndex(mgo.Index{
    Key: []string{"email"},
    Unique: true,
    DropDups: true,
    Background: true, // See notes.
    Sparse: true,
})
```


## Create DAO

Definitions can be instantiated as many times as you want :)

```go
users := kip.Create("Users")
users.Database = NewDatabase("localhost", "demo")
```

## CRUD: Create

```go
john := users.Create()
```

## CRUD: Retrieve

Objects can be retrieved in three ways:

* FindOne
* FindById
* Find

### FindOne

Retrieve one item based on a query.

If there is no matching objects, nil is returned.

It will panic if an unexpected error happens.

```go
john := users.FindOne(bson.M{"name": "John"})
```

### FindById

Retrieve one item by `_id`.

It is a particular case of `FindOne` with the query `bson.M{"_id": <id> }`.


### Find

Retrieve a cursor...


## CRUD: Update

```go
TODO!
```

## CRUD: Delete

```go
john.Delete()
```
