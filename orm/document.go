package orm

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/arangodb/go-driver"
)

type DocumentInterface interface {
	SetMeta(meta driver.DocumentMeta)
	SetKey(key string)
	GetBaseFieldNameByTag(tag string) string
}

type Document struct {
	ID  string `json:"_id,omitempty"`
	Rev string `json:"-"`
	Key string `json:"_key,omitempty"`

	Meta driver.DocumentMeta `json:"-"`
}

func (d Document) GetBaseFieldNameByTag(tag string) string {
	switch tag {
	case "_id":
		return "ID"
	case "_rev":
		return "Rev"
	case "_key":
		return "Key"
	}

	return ""
}

func (d Document) SetMeta(meta driver.DocumentMeta) {
	d.Meta = meta
}

func (d Document) SetKey(key string) {
	d.Key = key
}

func LoadFromMap(doc DocumentInterface, docMap map[string]interface{}) {

	for field, value := range docMap {
		var fieldName string
		if strings.HasPrefix(field, "_") {
			fieldName = doc.GetBaseFieldNameByTag(field)
		} else {
			fieldName = getFieldName(field, "json", doc)
		}
		SetField(doc, fieldName, value)
	}
}

func SetField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New("Provided value type didn't match obj field type")
	}

	structFieldValue.Set(val)
	return nil
}

func getFieldName(tag, key string, s interface{}) (fieldname string) {

	rt := reflect.TypeOf(s).Elem()
	if rt.Kind() != reflect.Struct {
		panic("bad type")
	}
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		v := strings.Split(f.Tag.Get(key), ",")[0] // use split to ignore tag "options" like omitempty, etc.
		if v == tag {
			return f.Name
		}
	}
	return ""
}
