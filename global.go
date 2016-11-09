package kip

var global *Kip = nil

func init() {
	global = NewKip()
}

func Define(c *Collection) {
	global.Define(c)
}

func NewDao(name string, db *Database) *Instance {
	return global.NewDao(name, db)
}
