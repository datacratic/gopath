// Copyright (c) 2014 Datacratic. All rights reserved.

package path

import (
	"reflect"
)

var errorType = reflect.TypeOf((*error)(nil)).Elem()

func isSetterFor(fn reflect.Value, value reflect.Value) bool {
	typ := fn.Type()

	return fn.Kind() == reflect.Func && !fn.IsNil() &&
		(typ.NumIn() == 1 && typ.In(0) == value.Type()) &&
		(typ.NumOut() == 0 || (typ.NumOut() == 1 && typ.Out(0) == errorType))
}

func callSetter(fn reflect.Value, value reflect.Value) (err error) {
	results := fn.Call([]reflect.Value{value})

	if len(results) > 0 && !results[0].IsNil() {
		err = results[0].Interface().(error)
	}

	return
}

func isGetter(fn reflect.Value) bool {
	typ := fn.Type()

	return fn.Kind() == reflect.Func && !fn.IsNil() &&
		(typ.NumIn() == 0) &&
		(typ.NumOut() == 1 || (typ.NumOut() == 2 && typ.Out(1) == errorType))
}

func callGetter(fn reflect.Value) (result reflect.Value, err error) {
	results := fn.Call([]reflect.Value{})

	result = results[0]
	if len(results) > 1 && !results[1].IsNil() {
		err = results[1].Interface().(error)
	}

	return
}
