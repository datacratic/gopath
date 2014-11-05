// Copyright (c) 2014 Datacratic. All rights reserved.

/*

Package path provides a simple utility to navigate go objects using string
paths.

A path is a series of components seperated by a '.' character. As an example:

    value, err := New("A.B").Get(obj)

will access field B of field A of object obj.

Path supports just about all go constructs with with the following caveats:
Pointers will automatically be dereferenced when accessed. Slices and Arrays can
only be traversed using unsigned integers as path components. Only maps that use
strings as keys can be traversed. Functions must take no arguments and return a
single value and an optional error value. The returned value will be used to
dereference the rest of the path. Channels are not supported at all.

A wildcard component, denoted by the '*' character, is also available when
using the GetAll to return all the values that match the path pattern.

A translation mechanism is available to convert JSON paths into paths usable by
gopath. This is accomplished by creating an alias table using the JSONAliases
function which is then used to translate paths using the Path.Translate
function.

*/
package path
