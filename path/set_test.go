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
	U Interface
	V interface{}

	PI *int
	PS *GetStruct
	PA []*int

	AU []Interface
	AV []interface{}

	MI map[string]int
	MS map[string]GetStruct
	MA map[string][]int
	MM map[string]map[string]int
	MU map[string]Interface
	MV map[string]interface{}

	MPI map[string]*int
	MPS map[string]*GetStruct

	C  chan int
	MC map[string]chan int

	F  func(i int)
	MF map[string]func(i int)
}

func (s *SetStruct) Fn(i int) {
	s.I = i
}

func (s *SetStruct) FnErr(i int) error {
	return fmt.Errorf("BOOM")
}

func (s *SetStruct) FnNilErr(i int) error {
	s.I = i
	return nil
}

func TestSet(t *testing.T) {
	setFail(t, "sx", "I", SetStruct{}, intPtr(10))
	setFail(t, "sx", "I", SetStruct{}, 10.0)

	n := 16

	s0 := new(SetStruct)
	setObj(t, "s0", "I", s0, 10)
	setObj(t, "s0", "I", s0, 10)
	setObj(t, "s0", "S.A", s0, 10)

	setFail(t, "s0", "U.A", s0, 10)
	setObj(t, "s0", "U", s0, new(GetStruct))
	setObj(t, "s0", "U.A", s0, 10)
	setObj(t, "s0", "U", s0, new(GetStruct))

	setFail(t, "s0", "V.A", s0, 10)
	setObj(t, "s0", "V", s0, new(GetStruct))
	setObj(t, "s0", "V.A", s0, 10)
	setObj(t, "s0", "V", s0, new(GetStruct))

	for i := 0; i < n; i++ {
		setObj(t, "s0", fmt.Sprintf("A.%d", i), s0, 10)
		setObj(t, "s0", fmt.Sprintf("A.%d", i), s0, 11)
	}

	setObj(t, "s0", "PI", s0, intPtr(12))
	setObj(t, "s0", "PI", s0, intPtr(13))
	setObj(t, "s0", "PS.A", s0, 12)
	setObj(t, "s0", "PS.A", s0, 13)
	for i := 0; i < n; i++ {
		setObj(t, "s0", fmt.Sprintf("PA.%d", i), s0, intPtr(10))
		setObj(t, "s0", fmt.Sprintf("PA.%d", i), s0, intPtr(11))
	}

	for i := 0; i < n; i++ {
		setFail(t, "s0", fmt.Sprintf("AU.%d.A", i), s0, 10)
		setObj(t, "s0", fmt.Sprintf("AU.%d", i), s0, new(GetStruct))
		setObj(t, "s0", fmt.Sprintf("AU.%d.A", i), s0, 10)
		setObj(t, "s0", fmt.Sprintf("AU.%d", i), s0, new(GetStruct))

		setFail(t, "s0", fmt.Sprintf("AV.%d.A", i), s0, 10)
		setObj(t, "s0", fmt.Sprintf("AV.%d", i), s0, new(GetStruct))
		setObj(t, "s0", fmt.Sprintf("AV.%d.A", i), s0, 10)
		setObj(t, "s0", fmt.Sprintf("AV.%d", i), s0, new(GetStruct))
	}

	setObj(t, "s0", "MI.X", s0, 14)
	setObj(t, "s0", "MI.X", s0, 15)
	setObj(t, "s0", "MI.Y", s0, 16)

	// This may look odd but a struct in a map is not addressable which means
	// that it can't be modified through reflection. Unlike map[string]int, it's
	// also too far dissociated from the map for it to detect that the thing in
	// in a map and have set fudge the results.
	setFail(t, "s0", "MS.X.A", s0, 14)
	setFail(t, "s0", "MS.X.A", s0, 15)
	setFail(t, "s0", "MS.Y.A", s0, 16)

	for i := 0; i < n; i++ {
		setObj(t, "s0", fmt.Sprintf("MA.X.%d", i), s0, 14)
		setObj(t, "s0", fmt.Sprintf("MA.Y.%d", i), s0, 15)
	}

	setObj(t, "s0", "MM.X.Z", s0, 14)
	setObj(t, "s0", "MM.X.Z", s0, 14)
	setObj(t, "s0", "MM.Y.Z", s0, 14)

	setFail(t, "s0", "MU.X.A", s0, 10)
	setObj(t, "s0", "MU.X", s0, new(GetStruct))
	setObj(t, "s0", "MU.X.A", s0, 10)
	setObj(t, "s0", "MU.X", s0, new(GetStruct))

	setFail(t, "s0", "MV.X.A", s0, 10)
	setObj(t, "s0", "MV.X", s0, new(GetStruct))
	setObj(t, "s0", "MV.X.A", s0, 10)
	setObj(t, "s0", "MV.X", s0, new(GetStruct))

	setObj(t, "s0", "MPI.X", s0, intPtr(17))
	setObj(t, "s0", "MPI.X", s0, intPtr(18))
	setObj(t, "s0", "MPI.Y", s0, intPtr(19))

	setObj(t, "s0", "MPS.X.A", s0, 17)
	setObj(t, "s0", "MPS.X.A", s0, 18)
	setObj(t, "s0", "MPS.Y.A", s0, 19)

	setObjFn(t, "s0", "Fn", s0, 20)
	setFail(t, "s0", "FnErr", s0, 21)
	setObjFn(t, "s0", "FnNilErr", s0, 22)

	fn := func(i int) { s0.I = i }

	setFail(t, "s0", "F", s0, 23)
	New("F").Set(s0, fn)
	setObjFn(t, "s0", "F", s0, 24)

	setFail(t, "s0", "MF.X", s0, 25)
	New("MF.X").Set(s0, fn)
	setObjFn(t, "s0", "MF.X", s0, 26)

	// setChan(t, "s0", "C", s0, 30)
	// setChan(t, "s0", "C", s0, 31)
	// setChan(t, "s0", "MC.X", s0, 32)
	// setChan(t, "s0", "MC.X", s0, 33)
	// setChan(t, "s0", "MC.Y", s0, 33)
}

func intPtr(i int) *int {
	p := new(int)
	*p = i
	return p
}

func setObj(t *testing.T, title string, path string, obj, exp interface{}) {
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

func setObjFn(t *testing.T, title string, path string, obj, exp interface{}) {
	if err := New(path).Set(obj, exp); err != nil {
		t.Errorf("FAIL(%s): set %s -> %s", title, path, err)
		return
	}

	if value, err := New("I").Get(obj); err != nil {
		t.Errorf("FAIL(%s): get %s -> %s", title, path, err)

	} else if value != exp {
		t.Errorf("FAIL(%s): get %s -> %v != %v", title, path, value, exp)
	}
}

func setChan(t *testing.T, title string, path string, obj, exp interface{}) {
	if err := New(path).Set(obj, exp); err != nil {
		t.Errorf("FAIL(%s): send %s -> %s", title, path, err)
		return
	}

	if value, err := New(path + ".1").Get(obj); err != nil {
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
