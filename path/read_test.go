// Copyright (c) 2014 Datacratic. All rights reserved.

package path

import (
	"encoding/json"
	"testing"
)

func TestRead(t *testing.T) {
	var obj struct {
		I  int
		S  GetStruct
		PS *GetStruct

		M   map[string]int
		MS  map[string]GetStruct
		MPS map[string]*GetStruct
		A   []GetStruct
	}

	obj.I = 123
	obj.S = GetStruct{1, 1}
	obj.PS = &GetStruct{2, 2}

	obj.M = map[string]int{"A": 1, "B": 2}

	sx, sy := &GetStruct{11, 11}, &GetStruct{21, 21}
	obj.MS = map[string]GetStruct{"I": *sx, "J": *sy}
	obj.MPS = map[string]*GetStruct{"X": sx, "Y": sy}

	readInt(t, "I", &obj, 123)
	readStruct(t, "S", &obj, obj.S)
	readPtr(t, "PS", &obj, obj.PS)

	readInt(t, "M.A", &obj, 1)
	readInt(t, "M.B", &obj, 2)

	readStruct(t, "MS.I", &obj, *sx)
	readStruct(t, "MS.J", &obj, *sy)

	readPtr(t, "MPS.X", &obj, sx)
	readPtr(t, "MPS.Y", &obj, sy)

	a := []GetStruct{{4, 4}}
	obj.A = []GetStruct{{4, 4}}
	readArr(t, "A", &obj, a)
}

func TestReadJson(t *testing.T) {
	body := `{ "a":{ "x": 10 }}`

	var obj interface{}
	if err := json.Unmarshal([]byte(body), &obj); err != nil {
		t.Fatal(err)
	}

	readInt(t, "a.x", obj, 10)
}

func read(t *testing.T, path string, obj, dest interface{}) {
	if err := New(path).Read(obj, dest); err != nil {
		t.Errorf("FAIL(%s): read into '%T': %s", path, dest, err)
	}
}

func readFail(t *testing.T, path string, obj, dest interface{}) {
	if err := New(path).Read(obj, dest); err == nil {
		t.Errorf("FAIL(%s): read into '%T': expected error", path, dest)
	}
}

func readInt(t *testing.T, path string, obj interface{}, exp int) {
	var dest int
	if read(t, path, obj, &dest); dest != exp {
		t.Errorf("FAIL(%s): readInt into '%T': %d != %d", path, dest, dest, exp)
	}
}

func readStruct(t *testing.T, path string, obj interface{}, exp GetStruct) {
	var dest GetStruct
	if read(t, path, obj, &dest); dest.A != exp.A {
		t.Errorf("FAIL(%s): readStruct into '%T': %d != %d", path, dest, dest.A, exp.A)
	}
}

func readPtr(t *testing.T, path string, obj interface{}, exp *GetStruct) {
	var dest *GetStruct
	if read(t, path, obj, &dest); dest != exp {
		t.Errorf("FAIL(%s): readPtr into '%T': %d != %d", path, dest, dest, exp)
	}
}

func readArr(t *testing.T, path string, obj interface{}, exp []GetStruct) {
	var dest []GetStruct
	if read(t, path, obj, &dest); dest[0].A != exp[0].A {
		t.Errorf("FAIL(%s): readArr into '%T': %d != %d", path, dest, dest, exp)
	}
}

func readAll(t *testing.T, path string, obj, dest interface{}) {
	if err := New(path).ReadAll(obj, dest); err != nil {
		t.Errorf("FAIL(%s): read into '%T': %s", path, dest, err)
	}
}

func readIntArr(t *testing.T, path string, obj interface{}, exp []int) {
	var dest []int
	readAll(t, path, obj, &dest)
	if dest[0] != exp[0] {
		t.Errorf("FAIL(%s): readArr into '%T': %d != %d", path, dest, dest, exp)
	}

}

func TestReadJsonArray(t *testing.T) {
	body := `{ "a": [ 1, 2, 3] }`

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(body), &obj); err != nil {
		t.Fatal(err)
	}

	readIntArr(t, "a.*", obj, []int{1, 2, 3})
}
