package kip

import (
	"sort"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
)

type Database struct {
	hosts    string
	database string
	session  *mgo.Session
}

var MONGO_DIAL_TIMEOUT = 0 * time.Second
var MONGO_SYNC_TIMEOUT = 3 * time.Second
var MONGO_SOCKET_TIMEOUT = 3 * time.Second

var sessions_by_hosts = map[string]*mgo.Session{}

func NewDatabase(hosts, database string) (*Database, error) {

	hosts = normalize_hosts(hosts)

	// Initialize Database
	db := &Database{
		hosts:    hosts,
		database: database,
		session:  nil,
	}

	// Check if a session already exists
	if previous_session, exists := sessions_by_hosts[hosts]; exists {
		db.session = previous_session
		return db, nil
	}

	// Create new session
	session, err := mgo.DialWithTimeout(hosts, MONGO_DIAL_TIMEOUT)
	if nil != err {
		return nil, err
	}

	// Enable autoreconnect and autoreset
	session.SetMode(mgo.Eventual, true)

	// Set timeouts
	session.SetSyncTimeout(MONGO_SYNC_TIMEOUT)
	session.SetSocketTimeout(MONGO_SOCKET_TIMEOUT)

	sessions_by_hosts[hosts] = session

	db.session = session

	return db, nil
}

func Close(hosts string) {
	if session, exists := sessions_by_hosts[hosts]; exists {
		session.Close()
	}
}

func CloseAll() {
	for i, session := range sessions_by_hosts {
		delete(sessions_by_hosts, i)
		session.Close()
	}
}

func (this *Database) C(collection string) *mgo.Collection {
	return this.session.DB(this.database).C(collection)
}

func normalize_hosts(hosts string) string {

	separator := ","

	parts := strings.Split(hosts, separator)

	for k, v := range parts {
		parts[k] = strings.Trim(v, " ")
	}

	hosts = strings.Join(parts, separator)

	sort.Strings(parts)

	return hosts
}
