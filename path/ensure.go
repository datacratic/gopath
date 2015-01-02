// Copyright (c) 2014 Datacratic. All rights reserved.

package path

import (
	"fmt"
	"reflect"
)

func isNillable(obj reflect.Value) bool {
	switch obj.Kind() {
	case reflect.Ptr, reflect.Func, reflect.Interface, reflect.Chan, reflect.Map, reflect.Slice:
		return true
	default:
		return false
	}
}

func ensure(head, tail P, obj reflect.Value, ctx *Context) error {
	if !isNillable(obj) || !obj.IsNil() {
		return nil
	}

	// If we're at the end of the path then a nil value is not invalid.
	if len(tail) == 0 {
		return nil
	}

	if !ctx.CreateIfMissing {
		return ErrMissing
	}

	if !obj.CanSet() {
		return fmt.Errorf("unable to ensure '%s' at '%s'", obj, head)
	}

	switch obj.Kind() {

	case reflect.Ptr:
		obj.Set(reflect.New(obj.Type().Elem()))

	case reflect.Map:
		obj.Set(reflect.MakeMap(obj.Type()))

	case reflect.Slice:
		obj.Set(reflect.MakeSlice(obj.Type(), 0, 0))

	case reflect.Chan:
		obj.Set(reflect.MakeChan(obj.Type(), 1))

	case reflect.Interface:
		obj.Set(reflect.Zero(obj.Type()))

	default:
		return fmt.Errorf("unable to create '%s' at '%s'", obj, head)
	}

	return nil
}

func zero(head P, mid string, typ reflect.Type) (value reflect.Value, err error) {
	switch typ.Kind() {

	case reflect.Ptr:
		value = reflect.New(typ.Elem())

	case reflect.Map:
		value = reflect.MakeMap(typ)

	case reflect.Slice:
		value = reflect.MakeSlice(typ, 0, 0)

	case reflect.Chan:
		value = reflect.MakeChan(typ, 1)

	case reflect.Func, reflect.Interface, reflect.UnsafePointer:
		value = reflect.Zero(typ)

	case reflect.Invalid:
		err = fmt.Errorf("unable to create '%s' at '%s'", typ, append(head, mid))

	default:
		value = reflect.Zero(typ)
	}

	return
}

func ensureMapKey(head P, mid string, obj, key reflect.Value, ctx *Context) (value reflect.Value, err error) {
	value = obj.MapIndex(key)

	if value.Kind() != reflect.Invalid {
		if !ctx.CreateIfMissing {
			return
		}

		if !isNillable(value) || !value.IsNil() {
			return
		}
	}

	if !ctx.CreateIfMissing {
		err = ErrMissing
		return
	}

	if value, err = zero(head, mid, obj.Type().Elem()); err != nil {
		return
	}

	obj.SetMapIndex(key, value)
	return
}

func ensureSliceIndex(
	head P, mid string, obj reflect.Value, index int, ctx *Context) (value reflect.Value, err error) {

	if index < 0 {
		err = fmt.Errorf("invalid index %d < 0 at '%s'", index, head)
		return
	}

	if index < obj.Len() {
		value = obj.Index(index)
		return
	}

	if !ctx.CreateIfMissing {
		err = ErrMissing
		return
	}

	expanded := obj

	for i := obj.Len(); i <= index; i++ {
		if value, err = zero(head, mid, obj.Type().Elem()); err != nil {
			return
		}

		expanded = reflect.Append(expanded, value)
	}

	// Value is not addresable so we need to grabb the value from the array.
	value = expanded.Index(index)

	if obj.CanSet() {
		obj.Set(expanded)

	} else if parent := ctx.Parent(); parent.Kind() == reflect.Map {
		parent.SetMapIndex(reflect.ValueOf(head.Last()), expanded)

	} else {
		err = fmt.Errorf("value is not addreseable at '%s'", append(head, mid))
	}

	return
}
