// Copyright (c) 2014 Datacratic. All rights reserved.

package path

import (
	"reflect"
)

func (path P) Read(obj, dest interface{}) (err error) {
	value := reflect.ValueOf(dest)

	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	if !value.CanSet() {
		panic("dest must be setable")
	}

	fn := func(_ P, ctx *Context) (bool, error) {
		result := ctx.Value()

		if result.Type().ConvertibleTo(value.Type()) {
			result = result.Convert(value.Type())
		}

		if !result.Type().AssignableTo(value.Type()) {
			return false, ErrInvalidType
		}

		value.Set(result)
		return false, nil
	}

	return path.Apply(obj, &Context{Fn: fn})
}

func (path P) ReadAll(obj, dest interface{}) (interface{}, error) {
	value := reflect.ValueOf(dest)

	if value.Kind() != reflect.Slice {
		panic("can only read slice values")
	}

	elem := value.Type().Elem()

	fn := func(_ P, ctx *Context) (bool, error) {
		result := ctx.Value()

		if result.Type().ConvertibleTo(elem) {
			result = result.Convert(elem)
		}

		if !result.Type().AssignableTo(value.Type(elem)) {
			return false, ErrInvalidType
		}

		value = reflect.Append(value, result)
		return true, nil
	}

	return value, path.Apply(obj, &Context{Fn: fn})
}
