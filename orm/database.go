package orm

import (
	"log"
	"reflect"

	"github.com/arangodb/go-driver"
)

func (dc *DatabaseConnection) FindByQuery(query string, bindVars map[string]interface{}, docType DocumentInterface) ([]DocumentInterface, error) {

	var docs []DocumentInterface
	cursor, err := dc.currentDatabase.Query(dc.currentContext, query, bindVars)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	for {
		var docMap map[string]interface{}
		newDoc := reflect.New(reflect.TypeOf(docType)).Interface().(DocumentInterface)
		meta, err := cursor.ReadDocument(dc.currentContext, &docMap)

		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			panic(err)
		}

		LoadFromMap(newDoc, docMap)
		newDoc.SetMeta(meta)

		docs = append(docs, newDoc)
	}

	return docs, nil
}

func (dc *DatabaseConnection) DB(database string) *DatabaseConnection {

	dc.currentDatabase = dc.GetDB(database)
	dc.currentCollection = nil

	return dc
}

func (dc *DatabaseConnection) GetDB(database string) driver.Database {

	db, err := dc.httpClient.Database(dc.currentContext, database)
	if err != nil {
		panic(err)
	}

	return db
}

func (dc *DatabaseConnection) Collection(collectionName string) *Collection {
	collection, err := dc.currentDatabase.Collection(dc.currentContext, collectionName)
	if err != nil {
		panic(err)
	}

	currentCollection := &Collection{
		Name:               collectionName,
		collection:         collection,
		databaseConnection: dc,
	}

	return currentCollection
}
