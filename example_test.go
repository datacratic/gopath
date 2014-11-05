// Copyright (c) 2014 Datacratic. All rights reserved.

package path_test

import (
	"github.com/datacratic/gopath"

	"fmt"
	"reflect"
)

type Object struct {
	Data map[string]interface{} `json:"blob"`
}

func Example_Path() {
	obj := &Object{
		Data: map[string]interface{}{
			"foo":  "bar",
			"blah": []string{"bleh", "bloh"},
		},
	}

	// Here we use a path to access the value of the key 'foo' in the map Data
	// of the object obj.
	bar, err := path.New("Data.foo").Get(obj)
	fmt.Printf("get(Data.foo): %s, %v\n", bar, err)

	// Here we access the key blah and fetch the first element of the array.
	bleh, err := path.New("Data.blah.0").Get(obj)
	fmt.Printf("get(Data.blah.0): %s, %v\n", bleh, err)

	// We can also fetch all the elements of an array or map by using the *
	// wildcard.
	array, err := path.New("Data.blah.*").GetAll(obj)
	fmt.Printf("get(Data.blah.*): %s, %v\n", array, err)

	// To use JSON paths we first need to create an alias table of the type we
	// which to access with a JSON path.
	aliases := path.JSONAliases(reflect.TypeOf(obj))

	// We can then translate our JSON path into a path that is usable by the
	// gopath package.
	p := path.New("blob.foo")
	fmt.Printf("%s -> %s\n", p, p.Translate(aliases))

	// Output:
	// get(Data.foo): bar, <nil>
	// get(Data.blah.0): bleh, <nil>
	// get(Data.blah.*): [bleh bloh], <nil>
	// blob.foo -> Data.foo
}
