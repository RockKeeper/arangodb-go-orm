package debug

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func Dump(object interface{}) {
	fmt.Println("------------ Dump Start -------------------")
	fmt.Printf("%v", object)
	fmt.Println("")
	fmt.Println("------------ Dump End ---------------------")
	fmt.Println("")
}

func DumpJson(object interface{}) {
	bytes, _ := jsoniter.MarshalIndent(object, "", "    ")
	fmt.Println("------------ Dump Start -------------------")
	fmt.Println(string(bytes))
	fmt.Println("------------ Dump End ---------------------")
	fmt.Println("")
}
