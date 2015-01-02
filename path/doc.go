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
strings as keys can be traversed. The returned value will be used to
dereference the rest of the path.

A wildcard component, denoted by the '*' character, is also available when
using the GetAll to return all the values that match the path pattern.

For channels a wildcard component can be provided to read all values until the
channel is closed or a count which to indicate the number of values to read.

Functions contains special handling whereby to path through a function the '()'
component needs to be provided and the function must take no input arguments and
output a single value with an optional error. To set, the function must take a
single input argument and return at most a single error argument. The errors
returned by function calls will be reported as errors from the pathing function.

A translation mechanism is available to convert JSON paths into paths usable by
gopath. This is accomplished by creating an alias table using the JSONAliases
function which is then used to translate paths using the Path.Translate
function.

*/
package path
