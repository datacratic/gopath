// Copyright (c) 2014 Datacratic. All rights reserved.

package path

import ()

// Get fetches the first value in the given object that matches the
// path. Returns ErrMissing if the path could not be completed due to a nil
// field, a missing array index or a missing map value.
func (path P) Get(obj interface{}) (result interface{}, err error) {
	fn := func(_ P, ctx *Context) (bool, error) {
		result = ctx.Value().Interface()
		return false, nil
	}

	err = path.Apply(obj, &Context{Fn: fn})
	return
}

// GetAll fetches all the values in the given object that matches the
// path. Returns ErrMissing if the path could not be completed due to a nil
// field, a missing array index or a missing map value. Note that missing
// components are ignored if they are encountered after a wildcard path
// component.
func (path P) GetAll(obj interface{}) (result []interface{}, err error) {
	fn := func(_ P, ctx *Context) (bool, error) {
		result = append(result, ctx.Value().Interface())
		return true, nil
	}

	err = path.Apply(obj, &Context{Fn: fn})
	return
}
