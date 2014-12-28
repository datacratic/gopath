// Copyright (c) 2014 Datacratic. All rights reserved.

package path

import (
	"bytes"
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
