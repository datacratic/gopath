// Copyright (c) 2014 Datacratic. All rights reserved.

package path

import (
	"errors"
)

// ErrMissing is an error that indicates that the path could not be found in the
// object due to a nil pointer, missing array index or missing map key.
var ErrMissing = errors.New("unable to path through the object")

// ErrInvalidType indicates that the type of the value to be written did not
// match the destination's type.
var ErrInvalidType = errors.New("type mismatch")

// ErrNil indicates that value is nil
var ErrNil = errors.New("value is nil")
