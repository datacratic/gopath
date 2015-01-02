// Copyright (c) 2014 Datacratic. All rights reserved.

package path

import (
	"fmt"
	"reflect"
	"strconv"
)

// Apply applies the given context to the object using the path.
func (path P) Apply(obj interface{}, ctx *Context) (err error) {
	return apply(reflect.ValueOf(obj), P{}, path, ctx)
}

func apply(obj reflect.Value, head, tail P, ctx *Context) (err error) {
	if err := ensure(head, tail, obj, ctx); err != nil {
		return err
	}

	ctx.push(obj)

	if len(tail) == 0 {
		var cont bool
		cont, err = ctx.Fn(head, ctx)
		ctx.stop = !cont

	} else {
		switch obj.Kind() {

		case reflect.Interface, reflect.Ptr:
			err = applyToPtr(obj, head, tail, ctx)

		case reflect.Struct:
			err = applyToStruct(obj, head, tail[0], tail[1:], ctx)

		case reflect.Array, reflect.Slice:
			err = applyToSlice(obj, head, tail[0], tail[1:], ctx)

		case reflect.Map:
			err = applyToMap(obj, head, tail[0], tail[1:], ctx)

		case reflect.Func:
			err = applyToFunc(obj, head, tail[0], tail[1:], ctx)

		case reflect.Chan:
			err = applyToChan(obj, head, tail[0], tail[1:], ctx)

		default:
			err = fmt.Errorf("invalid kind '%s' in at '%s'", obj, head)
		}

	}

	ctx.pop()
	return
}

func applyToPtr(obj reflect.Value, head, tail P, ctx *Context) error {

	if result := obj.MethodByName(tail[0]); result.Kind() != reflect.Invalid {
		return apply(result, append(head, tail[0]), tail[1:], ctx)
	}

	return apply(obj.Elem(), head, tail, ctx)
}

func applyToFunc(obj reflect.Value, head P, mid string, tail P, ctx *Context) error {
	if mid != "()" {
		return fmt.Errorf("missing required '()' pathing component at '%s'", head)
	}

	if !isGetter(obj) {
		return fmt.Errorf("invalid return signature for function '%s' at '%s'", mid, head)
	}

	result, err := callGetter(obj)

	if err != nil {
		return err
	}

	return apply(result, append(head, mid), tail, ctx)
}

func applyToStruct(obj reflect.Value, head P, mid string, tail P, ctx *Context) error {
	if mid != "*" {
		result := obj.FieldByName(mid)

		if result.Kind() == reflect.Invalid {
			result = obj.MethodByName(mid)
		}

		if result.Kind() == reflect.Invalid {
			return fmt.Errorf("no field '%s' in type '%s' at '%s'", mid, obj.Type(), head)
		}

		return apply(result, append(head, mid), tail, ctx)
	}

	typ := obj.Type()

	for i := 0; i < typ.NumField() && !ctx.stop; i++ {
		if err := applyToStruct(obj, head, typ.Field(i).Name, tail, ctx); err != nil && err != ErrMissing {
			return err
		}
	}

	return nil
}

func applyToSlice(obj reflect.Value, head P, mid string, tail P, ctx *Context) error {
	if mid != "*" {
		index, err := strconv.ParseInt(mid, 10, 32)
		if err != nil {
			return fmt.Errorf("invalid index '%s' at '%s' -> %s", mid, head, err)
		}

		if int(index) < 0 {
			return fmt.Errorf("invalid index %d < 0 at '%s'", index, head)
		}

		result, err := ensureSliceIndex(head, mid, obj, int(index), ctx)
		if err != nil {
			return err
		}

		return apply(result, append(head, mid), tail, ctx)
	}

	for i := 0; i < obj.Len() && !ctx.stop; i++ {
		if err := applyToSlice(obj, head, strconv.Itoa(i), tail, ctx); err != nil && err != ErrMissing {
			return err
		}
	}

	return nil
}

func applyToMap(obj reflect.Value, head P, mid string, tail P, ctx *Context) error {
	if mid != "*" {
		result, err := ensureMapKey(head, mid, obj, reflect.ValueOf(mid), ctx)
		if err != nil {
			return err
		}
		return apply(result, append(head, mid), tail, ctx)
	}

	if key := obj.Type().Key(); key.Kind() != reflect.String {
		return fmt.Errorf("unsupported key type '%s' for map '%s' at '%s'", key, mid, head)
	}

	keys := obj.MapKeys()

	for i := 0; i < len(keys) && !ctx.stop; i++ {
		if err := applyToMap(obj, head, keys[i].String(), tail, ctx); err != nil && err != ErrMissing {
			return err
		}
	}

	return nil
}

func applyToChan(obj reflect.Value, head P, mid string, tail P, ctx *Context) error {
	if dir := obj.Type().ChanDir(); dir != reflect.RecvDir && dir != reflect.BothDir {
		return fmt.Errorf("invalid channel direction '%s' at '%s'", dir, head)
	}

	switch mid {

	case "1":
		result, ok := obj.Recv()
		if !ok {
			return ErrMissing
		}

		return apply(result, append(head, mid), tail, ctx)

	case "*":
		for !ctx.stop {
			result, ok := obj.Recv()
			if !ok {
				return nil
			}

			if err := apply(result, append(head, mid), tail, ctx); err != nil && err != ErrMissing {
				return err
			}
		}

		return nil

	default:
		return fmt.Errorf("invalid channel component '%s' at '%s'", mid, head)

	}
}
