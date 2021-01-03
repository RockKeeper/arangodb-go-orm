package main

import (
	orm "gopkg.in/RockKeeper/arangodb-go-orm.v1/orm"
	debug "gopkg.in/RockKeeper/go-dump.v1"
)

type Foo struct {
	orm.Document
	Name string `json:"name"`
}

func main() {
	creds := &orm.DatabaseConnectionData{
		Host:     "http://localhost:8529",
		Username: "root",
		Password: "root",
		Database: "testdb",
	}

	db, _ := orm.NewDatabaseConnection(creds)

	// db.Collection("members").FindByID("peter", &foo)
	// orm.DumpJson(foo)
	var doc Foo
	docs, _ := db.Collection("members").FindAll(doc)
	debug.DumpJson(docs)

}
