// Copyright (c) 2014 Datacratic. All rights reserved.

package path

import (
	"reflect"
	"strings"
)

// Translate replaces any path component that have an alias within the given
// aliases map.
func (path P) Translate(aliases map[string]string) (result P) {
	for _, item := range path {
		if alias, ok := aliases[item]; ok {
			result = append(result, alias)
		} else {
			result = append(result, item)
		}
	}
	return
}

// JSONAliases crawls the given type and returns an alias map from JSON names to
// struct field names.
func JSONAliases(typ reflect.Type) map[string]string {
	aliases := make(map[string]string)
	jsonAliases(typ, aliases)
	return aliases
}

func jsonAliases(typ reflect.Type, aliases map[string]string) {
	switch typ.Kind() {

	case reflect.Interface, reflect.Ptr, reflect.Map, reflect.Array, reflect.Slice:
		jsonAliases(typ.Elem(), aliases)

	case reflect.Struct:

		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)

			list := strings.Split(field.Tag.Get("json"), ",")
			if len(list) > 0 {
				aliases[list[0]] = field.Name
			}

			jsonAliases(field.Type, aliases)
		}
	}
}
