package orm

import (
	"fmt"
	"log"
	"reflect"

	"github.com/arangodb/go-driver"
)

func (d *Document) FindAll(doc DocumentInterface) ([]DocumentInterface, error) {

	return d.FindByQuery(fmt.Sprintf("FOR d IN %s RETURN d", doc.GetCollection()), nil, doc)
}

func (d *Document) FindByQuery(query string, bindVars map[string]interface{}, doc DocumentInterface) ([]DocumentInterface, error) {

	var docs []DocumentInterface
	cursor, err := ConnectionManager.GetDB().Query(ConnectionManager.GetContext(), query, bindVars)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	for {
		var docMap map[string]interface{}
		newDoc := reflect.New(reflect.TypeOf(doc)).Interface().(DocumentInterface)
		meta, err := cursor.ReadDocument(ConnectionManager.GetContext(), &docMap)

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

type Collection struct {
	Name               string
	collection         driver.Collection
	databaseConnection *DatabaseConnection
}

func (c *Collection) FindByKey(key string, doc DocumentInterface) error {
	return c.FindByID(fmt.Sprintf("%s/%s", c.Name, key), doc)
}

func (c *Collection) FindByID(id string, doc DocumentInterface) error {

	var docMap map[string]interface{}
	meta, err := c.collection.ReadDocument(c.databaseConnection.currentContext, id, &docMap)
	LoadFromMap(doc, docMap)
	doc.SetMeta(meta)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func (c *Collection) FindByFilter(filter string, doc DocumentInterface) ([]DocumentInterface, error) {

	var docs []DocumentInterface

	query := fmt.Sprintf("FOR d IN %s FILTER @filter RETURN d", c.Name)
	bindVars := map[string]interface{}{
		"collection": c.Name,
		"filter":     filter,
	}

	cursor, err := c.databaseConnection.currentDatabase.Query(c.databaseConnection.currentContext, query, bindVars)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	for {
		var docMap map[string]interface{}
		meta, err := cursor.ReadDocument(c.databaseConnection.currentContext, &docMap)
		LoadFromMap(doc, docMap)
		doc.SetMeta(meta)
		docs = append(docs, doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			log.Fatal(err)
			return nil, err
		}
	}

	return docs, nil
}

// func (c *Collection) FindAll() ([]DocumentInterface, error) {
// 	return c.FindByQuery(fmt.Sprintf("FOR d IN %s RETURN d", c.Name), nil)
// }

// func (c *Collection) FindByQuery(query string, bindVars map[string]interface{}) ([]DocumentInterface, error) {

// 	var docs []DocumentInterface
// 	cursor, err := c.databaseConnection.currentDatabase.Query(c.databaseConnection.currentContext, query, bindVars)
// 	if err != nil {
// 		log.Fatal(err)
// 		return nil, err
// 	}
// 	for {
// 		var docMap map[string]interface{}
// 		newDoc := reflect.New(reflect.TypeOf(doc)).Interface().(DocumentInterface)
// 		meta, err := cursor.ReadDocument(c.databaseConnection.currentContext, &docMap)

// 		if driver.IsNoMoreDocuments(err) {
// 			break
// 		} else if err != nil {
// 			panic(err)
// 		}

// 		LoadFromMap(newDoc, docMap)
// 		newDoc.SetMeta(meta)

// 		docs = append(docs, newDoc)
// 	}

// 	return docs, nil
// }

// func (c *Collection) FindOneByQuery(query string, bindVars map[string]interface{}) (DocumentInterface, error) {
// 	cursor, err := c.databaseConnection.currentDatabase.Query(c.databaseConnection.currentContext, query, bindVars)
// 	if err != nil {
// 		log.Fatal(err)
// 		return nil, err
// 	}
// 	var docMap map[string]interface{}
// 	meta, err := cursor.ReadDocument(c.databaseConnection.currentContext, &docMap)
// 	LoadFromMap(doc, docMap)
// 	doc.SetMeta(meta)
// 	if err != nil {
// 		log.Fatal(err)
// 		return nil, err
// 	}

// 	return doc, nil
// }

func GetType(arr interface{}) reflect.Type {
	return reflect.TypeOf(arr).Elem()
}
