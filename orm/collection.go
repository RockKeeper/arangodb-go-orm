package orm

import (
	"fmt"
	"log"
	"reflect"

	"github.com/arangodb/go-driver"
)

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

func (c *Collection) FindAll(doc DocumentInterface) ([]DocumentInterface, error) {
	return c.databaseConnection.FindByQuery(fmt.Sprintf("FOR d IN %s RETURN d", c.Name), nil, doc)
}

func GetType(arr interface{}) reflect.Type {
	return reflect.TypeOf(arr).Elem()
}
