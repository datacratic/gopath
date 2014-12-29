// Copyright (c) 2014 Datacratic. All rights reserved.

package path

import (
	"fmt"
	"testing"
)

type SetStruct struct {
	I int
	S GetStruct
	A []int

	PI *int
	PS *GetStruct
	PA []*int

	MI map[string]int
	MS map[string]GetStruct
	MA map[string][]int
	MM map[string]map[string]int

	MPI map[string]*int
	MPS map[string]*GetStruct
}

func TestSet(t *testing.T) {
	setFail(t, "sx", "I", SetStruct{}, intPtr(10))
	setFail(t, "sx", "I", SetStruct{}, 10.0)

	n := 16

	s0 := new(SetStruct)
	setInt(t, "s0", "I", s0, 10)
	setInt(t, "s0", "I", s0, 10)
	setInt(t, "s0", "S.A", s0, 10)
	for i := 0; i < n; i++ {
		setInt(t, "s0", fmt.Sprintf("A.%d", i), s0, 10)
		setInt(t, "s0", fmt.Sprintf("A.%d", i), s0, 11)
	}

	setInt(t, "s0", "PI", s0, intPtr(12))
	setInt(t, "s0", "PI", s0, intPtr(13))
	setInt(t, "s0", "PS.A", s0, 12)
	setInt(t, "s0", "PS.A", s0, 13)
	for i := 0; i < n; i++ {
		setInt(t, "s0", fmt.Sprintf("PA.%d", i), s0, intPtr(10))
		setInt(t, "s0", fmt.Sprintf("PA.%d", i), s0, intPtr(11))
	}

	setInt(t, "s0", "MI.X", s0, 14)
	setInt(t, "s0", "MI.X", s0, 15)
	setInt(t, "s0", "MI.Y", s0, 16)

	// This may look odd but a struct in a map is not addressable which means
	// that it can't be modified through reflection. Unlike map[string]int, it's
	// also too far dissociated from the map for it to detect that the thing in
	// in a map and have set fudge the results.
	setFail(t, "s0", "MS.X.A", s0, 14)
	setFail(t, "s0", "MS.X.A", s0, 15)
	setFail(t, "s0", "MS.Y.A", s0, 16)

	for i := 0; i < n; i++ {
		setInt(t, "s0", fmt.Sprintf("MA.X.%d", i), s0, 14)
		setInt(t, "s0", fmt.Sprintf("MA.Y.%d", i), s0, 15)
	}

	setInt(t, "s0", "MM.X.Z", s0, 14)
	setInt(t, "s0", "MM.X.Z", s0, 14)
	setInt(t, "s0", "MM.Y.Z", s0, 14)

	setInt(t, "s0", "MPI.X", s0, intPtr(17))
	setInt(t, "s0", "MPI.X", s0, intPtr(18))
	setInt(t, "s0", "MPI.Y", s0, intPtr(19))

	setInt(t, "s0", "MPS.X.A", s0, 17)
	setInt(t, "s0", "MPS.X.A", s0, 18)
	setInt(t, "s0", "MPS.Y.A", s0, 19)
}

func intPtr(i int) *int {
	p := new(int)
	*p = i
	return p
}

func setInt(t *testing.T, title string, path string, obj, exp interface{}) {
	if err := New(path).Set(obj, exp); err != nil {
		t.Errorf("FAIL(%s): set %s -> %s", title, path, err)
		return
	}

	if value, err := New(path).Get(obj); err != nil {
		t.Errorf("FAIL(%s): get %s -> %s", title, path, err)

	} else if value != exp {
		t.Errorf("FAIL(%s): get %s -> %v != %v", title, path, value, exp)
	}
}

func setFail(t *testing.T, title string, path string, obj, value interface{}) {
	if err := New(path).Set(obj, value); err == nil {
		t.Errorf("FAIL(%s): %s -> expected failure", title, path)
	}
}
