package path

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func JsonSchema(t interface{}) string {
	typ := reflect.TypeOf(t)

	mp := map[string]interface{}{}
	js := jsonSchema(typ, mp)
	if j, err := json.Marshal(js); err != nil {
		return fmt.Sprint(err)
	} else {
		return string(j)
	}
}

func jsonSchema(typ reflect.Type, mp map[string]interface{}) interface{} {

	switch typ.Kind() {

	case reflect.Struct:
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)

			var name string
			list := strings.Split(field.Tag.Get("json"), ",")
			if len(list) > 0 && list[0] != "" {
				name = list[0]
			} else {
				name = field.Name
			}

			mp[name] = jsonSchema(field.Type, map[string]interface{}{})
		}
	case reflect.Ptr:
		return "*" + typ.Elem().Kind().String()
	case reflect.Slice:
		js := jsonSchema(typ.Elem(), map[string]interface{}{})
		if j, err := json.Marshal(js); err != nil {
			return err
		} else {
			return "[]" + strings.Replace(string(j), "\"", "", -1)
		}
	case reflect.Map:
		js := jsonSchema(typ.Elem(), map[string]interface{}{})
		if j, err := json.Marshal(js); err != nil {
			return err
		} else {
			return "map[" + typ.Key().Kind().String() + "]" + strings.Replace(string(j), "\"", "", -1)
		}
	default:
		return typ.Kind().String()
	}
	return mp
}
