package main

import (
<<<<<<< HEAD
	"gopkg.in/RockKeeper/arangodb-go-orm.v1/debug"
=======
>>>>>>> 90971f45efcc5afbee64317965121a4e285976e6
	orm "gopkg.in/RockKeeper/arangodb-go-orm.v1/orm"
	debug "gopkg.in/RockKeeper/go-dump.v1"
)

// -----------------------------------------------------------------
// Sample Model
// -----------------------------------------------------------------

type Foo struct {
	orm.Document
	Name string `json:"name"`
}

func (f Foo) GetCollection() string {
	return "members"
}

func (f Foo) FindAll() ([]orm.DocumentInterface, error) {
	return f.Document.FindAll(f)
}

// -----------------------------------------------------------------

func main() {
	creds := &orm.DatabaseConnectionData{
		Host:     "http://localhost:8529",
		Username: "root",
		Password: "root",
		Database: "joiner",
	}

	orm.InitDatabaseConnection(creds)

	var foo Foo
	docs, _ := foo.FindAll()

	debug.DumpJson(docs)
}
