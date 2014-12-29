// Copyright (c) 2014 Datacratic. All rights reserved.

package path

import (
	"reflect"
)

// Context contains the state of the path crawl.
type Context struct {

	// Fn is the function to be called when the path is fully matched. The P
	// argument is the fully qualified path where wildcard are replaced by their
	// actual component name. The Context argument contains the current state of
	// the path crawler and can be used to retrieve various information. The
	// bool return value indicates whether we should continue to explore the
	// path for additional matchin values.
	Fn func(P, *Context) (bool, error)

	// CreateIfMissing will prevent Missing errors from being returned by
	// filling in any missing components with their zero values. Arrays will be
	// expanded until they are large enough for the desired index. Channels will
	// be created with a default configuration. Note that interfaces, and
	// unsafe.Pointers will cause an error to be returned.
	CreateIfMissing bool

	stop   bool
	values []reflect.Value
}

// Value returns the current value being tracked by the path crawler.
func (ctx *Context) Value() reflect.Value {
	return ctx.values[len(ctx.values)-1]
}

// Parent returns the parent value being tracked by the path crawler. This can
// be useful to modify the values in maps which are not addresable.
func (ctx *Context) Parent() (parent reflect.Value) {
	if len(ctx.values) > 1 {
		parent = ctx.values[len(ctx.values)-2]
	}
	return
}

func (ctx *Context) push(value reflect.Value) {
	ctx.values = append(ctx.values, value)
}

func (ctx *Context) pop() {
	ctx.values = ctx.values[:len(ctx.values)-1]
}
