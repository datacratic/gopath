// Copyright (c) 2014 Datacratic. All rights reserved.

package path

import (
	"reflect"
)

// Read sets the given dest object to the content of the path applied to the
// given obj object. Returns ErrInvalidType if the type of the value can't be
// converted or assigned to dest. Panics if dest is can't be set.
func (path P) Read(obj, dest interface{}) (err error) {
	value := reflect.ValueOf(dest)

	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	if !value.CanSet() {
		panic("dest must be setable")
	}

	fn := func(_ P, ctx *Context) (bool, error) {

		for result := ctx.Value(); ; result = result.Elem() {

			if result.Type().ConvertibleTo(value.Type()) {
				result = result.Convert(value.Type())
			}

			if result.Type().AssignableTo(value.Type()) {
				value.Set(result)
				return false, nil
			}

			if result.Kind() != reflect.Interface && result.Kind() != reflect.Ptr {
				return false, ErrInvalidType
			}
		}

	}

	return path.Apply(obj, &Context{Fn: fn})
}

// ReadAll appends to the given dest slice to content of the path applied to the
// given obj object. Returns ErrInvalidType if the type of the value can't be
// converted or assigned to the member of dest. Panics if dest is not a slice.
func (path P) ReadAll(obj, dest interface{}) (interface{}, error) {
	value := reflect.ValueOf(dest)

	if value.Kind() != reflect.Slice {
		panic("can only read slice values")
	}

	elem := value.Type().Elem()

	fn := func(_ P, ctx *Context) (bool, error) {

		for result := ctx.Value(); ; result = result.Elem() {

			if result.Type().ConvertibleTo(elem) {
				result = result.Convert(elem)
			}

			if result.Type().AssignableTo(elem) {
				value = reflect.Append(value, result)
				return true, nil
			}

			if result.Kind() != reflect.Interface && result.Kind() != reflect.Ptr {
				return false, ErrInvalidType
			}
		}

	}

	return value, path.Apply(obj, &Context{Fn: fn})
}
