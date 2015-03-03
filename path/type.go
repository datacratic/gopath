// Copyright (c) 2014 Datacratic. All rights reserved.

package path

import (
	"reflect"
)

// Type returns the type of the first value in the given object that matches the
// path. Returns ErrMissing if the path could not be completed due to a nil
// field, a missing array index or a missing map value.
func (path P) Type(obj interface{}) (result reflect.Type, err error) {
	fn := func(_ P, ctx *Context) (bool, error) {
		result = ctx.Value().Type()
		return false, nil
	}

	err = path.Apply(obj, &Context{Fn: fn})
	return
}
