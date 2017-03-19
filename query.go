package kip

import mgo "gopkg.in/mgo.v2"

type Query struct {
	dao       *Dao
	mgo_query *mgo.Query
}

func (q *Query) Limit(n int) *Query {
	q.mgo_query = q.mgo_query.Limit(n)
	return q
}

func (q *Query) Select(selector interface{}) *Query {
	q.mgo_query = q.mgo_query.Select(selector)
	return q
}

func (q *Query) Skip(n int) *Query {
	q.mgo_query = q.mgo_query.Skip(n)
	return q
}

func (q *Query) Snapshot() *Query {
	q.mgo_query = q.mgo_query.Snapshot()
	return q
}

func (q *Query) Sort(fields ...string) *Query {
	q.mgo_query = q.mgo_query.Sort(fields...)
	return q
}

// Finalizers
func (q *Query) All(result interface{}) error {
	return q.mgo_query.All(result)
}

func (q *Query) Count() (n int, err error) {
	return q.mgo_query.Count()
}

func (q *Query) Iter() *mgo.Iter {
	return q.mgo_query.Iter()
}

func (q *Query) ForEach(f func(*Item)) *Query {

	i := q.mgo_query.Iter()

	item := q.dao.Create()
	for i.Next(item.Value) {
		item.saved = true

		f(item)

		item = q.dao.Create()
	}

	i.Close()

	return q
}

func (q *Query) One() {
	// TODO!
	// Already implemented by Instance
}
