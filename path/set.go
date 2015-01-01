// Copyright (c) 2014 Datacratic. All rights reserved.

package path

import (
	"fmt"
	"reflect"
)

// Set modifies the first value in the given object that matches the path to
// contain the given value. Returns ErrInvalidType if the type of the given
// value doesn't match the type of value at the path location in the given
// object. Returns an error if the object is not addresable and therefore not
// modifiable.
func (path P) Set(obj, value interface{}) error {
	fn := func(p P, ctx *Context) (bool, error) {
		return false, set(p, ctx, reflect.ValueOf(value))
	}

	return path.Apply(obj, &Context{CreateIfMissing: true, Fn: fn})
}

// SetAll modifies all the values in the given object that matches the path to
// contain the given value. Returns ErrInvalidType if the type of the given
// value doesn't match the type of value at the path location in the given
// object. Returns an error if the object is not addresable and therefore not
// modifiable.
func (path P) SetAll(obj, value interface{}) (err error) {
	fn := func(p P, ctx *Context) (bool, error) {
		return true, set(p, ctx, reflect.ValueOf(value))
	}

	return path.Apply(obj, &Context{CreateIfMissing: true, Fn: fn})
}

var errorType = reflect.TypeOf((*error)(nil)).Elem()

func set(path P, ctx *Context, value reflect.Value) error {
	obj := ctx.Value()

	if obj.Kind() == reflect.Chan {
		if obj.Type().Elem() == value.Type() {
			obj.Send(value)
			return nil
		}
	}

	if obj.Kind() == reflect.Func {
		if fn := obj.Type(); fn.NumIn() == 1 && fn.In(0) == value.Type() {
			if fn.NumOut() == 0 || (fn.NumOut() == 1 && fn.Out(0) == errorType) {
				if ret := obj.Call([]reflect.Value{value}); len(ret) > 0 {
					return ret[0].Interface().(error)
				}
			}
		}
	}

	if value.Type() != obj.Type() {
		if value.Type().ConvertibleTo(obj.Type()) {
			value = value.Convert(obj.Type())

		} else if obj.Kind() != reflect.Interface || !value.Type().Implements(obj.Type()) {
			return ErrInvalidType
		}
	}

	if obj.CanSet() {
		obj.Set(value)
		return nil
	}

	if parent := ctx.Parent(); parent.Kind() == reflect.Map {
		parent.SetMapIndex(reflect.ValueOf(path.Last()), value)
		return nil
	}

	return fmt.Errorf("unable to set '%s' at '%s'", obj, path)
}
