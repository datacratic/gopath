// Copyright (c) 2014 Datacratic. All rights reserved.

package path

import (
	"bytes"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

// P represents a path through an object seperated by .
type P []string

// New returns a new P object from a given path string.
func New(path string) P {
	return strings.Split(path, ".")
}

func (path P) String() string {
	buffer := new(bytes.Buffer)

	for i, item := range path {
		buffer.WriteString(item)

		if i < len(path)-1 {
			buffer.WriteString(".")
		}
	}

	return buffer.String()
}

type applyFn func(P, reflect.Value) bool

func applyToPtr(obj reflect.Value, head P, mid string, tail P, fn applyFn) (cont bool, err error) {
	result := obj.MethodByName(mid)
	if result.Kind() != reflect.Invalid {
		cont, err = applyToFunc(result, head, mid, tail, fn)
		return
	}

	cont, err = apply(obj.Elem(), head, append(P{mid}, tail...), fn)
	return
}

func applyToFunc(obj reflect.Value, head P, mid string, tail P, fn applyFn) (cont bool, err error) {
	t := obj.Type()

	if t.NumIn() != 0 {
		err = fmt.Errorf("Invalid arguments for func %s: %s", mid, t)
		return
	}

	var result reflect.Value

	if t.NumOut() == 1 {
		result = obj.Call([]reflect.Value{})[0]

	} else if t.NumOut() == 2 && t.Out(1) == reflect.TypeOf((*error)(nil)).Elem() {
		ret := obj.Call([]reflect.Value{})

		if !ret[1].IsNil() {
			err = ret[1].Interface().(error)
			return
		}

		result = ret[0]

	} else {
		err = fmt.Errorf("Invalid return for func %s: %s", mid, t)
		return
	}

	cont, err = apply(result, append(head, mid), tail, fn)
	return
}

func applyToStruct(obj reflect.Value, head P, mid string, tail P, fn applyFn) (cont bool, err error) {
	if mid != "*" {
		result := obj.FieldByName(mid)
		if result.Kind() == reflect.Invalid {
			result = obj.MethodByName(mid)
		}

		if result.Kind() == reflect.Invalid {
			err = fmt.Errorf("no field named %s in struct", mid)
			return
		}

		cont, err = apply(result, append(head, mid), tail, fn)
		return
	}

	t := obj.Type()

	for i := 0; i < t.NumField(); i++ {
		cont, err = applyToStruct(obj, head, t.Field(i).Name, tail, fn)
		if !cont {
			break
		}
	}

	return
}

func applyToSlice(obj reflect.Value, head P, mid string, tail P, fn applyFn) (cont bool, err error) {
	if mid != "*" {
		index, err := strconv.ParseInt(mid, 10, 32)
		if err != nil {
			return cont, err
		}

		if int(index) >= obj.Len() || int(index) < 0 {
			err = fmt.Errorf("index %d >= slice len %d", index, obj.Len())
			return cont, err
		}

		cont, err = apply(obj.Index(int(index)), append(head, mid), tail, fn)
		return cont, err
	}

	for i := 0; i < obj.Len(); i++ {
		cont, err = applyToSlice(obj, head, strconv.Itoa(i), tail, fn)
		if !cont {
			break
		}
	}

	return
}

func applyToMap(obj reflect.Value, head P, mid string, tail P, fn applyFn) (cont bool, err error) {
	if mid != "*" {
		result := obj.MapIndex(reflect.ValueOf(mid))
		if result.Kind() == reflect.Invalid {
			err = fmt.Errorf("map doesn't contain key %s", mid)
			return
		}

		cont, err = apply(result, append(head, mid), tail, fn)
		return
	}

	if key := obj.Type().Key(); key.Kind() != reflect.String {
		err = fmt.Errorf("Unsupported key type in map: %s", key)
		return
	}

	keys := obj.MapKeys()
	for _, key := range keys {
		cont, err = applyToMap(obj, head, key.String(), tail, fn)
		if !cont {
			break
		}
	}

	return
}

func apply(obj reflect.Value, head, tail P, fn applyFn) (cont bool, err error) {
	if len(tail) == 0 {
		cont = fn(head, obj)
		return
	}

	if obj.Kind() == reflect.Invalid {
		log.Panic("invalid kind in path")
	}

	switch obj.Kind() {

	case reflect.Interface, reflect.Ptr:
		cont, err = applyToPtr(obj, head, tail[0], tail[1:], fn)

	case reflect.Struct:
		cont, err = applyToStruct(obj, head, tail[0], tail[1:], fn)

	case reflect.Array, reflect.Slice:
		cont, err = applyToSlice(obj, head, tail[0], tail[1:], fn)

	case reflect.Map:
		cont, err = applyToMap(obj, head, tail[0], tail[1:], fn)

	default:
		err = fmt.Errorf("premature end of path")
	}

	return
}

// Apply applies to given function to all the element that match the path int he
// given obj.
func (path P) Apply(obj interface{}, fn func(P, interface{}) bool) (err error) {
	_, err = apply(reflect.ValueOf(obj), P{}, path, func(path P, value reflect.Value) bool {
		return fn(path, value.Interface())
	})
	return
}

// Get returns the first element that matches the path in the given object.
func (path P) Get(obj interface{}) (result interface{}, err error) {
	err = path.Apply(obj, func(_ P, value interface{}) bool {
		result = value
		return false
	})
	return
}

// GetAll returns all the elements that matches the path in the given object.
func (path P) GetAll(obj interface{}) (result []interface{}, err error) {
	err = path.Apply(obj, func(_ P, value interface{}) bool {
		result = append(result, value)
		return true
	})
	return
}
