// Copyright (c) 2014 Datacratic. All rights reserved.

package path

import (
	"fmt"
	"sort"
	"strconv"
	"testing"
)

type Interface interface {
	B() int
}

type Struct struct {
	A int `json:"alice"`
	Z int `json:"zebra"`
}

func (s *Struct) B() int {
	return s.A
}

func (s *Struct) C() (int, error) {
	if s.A == 0 {
		return 0, fmt.Errorf("blah")
	}
	return s.A, nil
}

func TestPathGet(t *testing.T) {
	mapVal := map[string]int{
		"A": 1,
		"B": 10,
	}
	getInt(t, "map", "A", mapVal, 1)
	getInt(t, "map", "B", mapVal, 10)
	getAllInt(t, "map", "*", mapVal, []int{1, 10})
	getFail(t, "map", "C", mapVal)
	getFail(t, "map", "A.B", mapVal)

	getInt(t, "mapPtr", "A", &mapVal, 1)
	getAllInt(t, "mapPtr", "*", &mapVal, []int{1, 10})
	getFail(t, "mapPtr", "C", &mapVal)

	arrayVal := []int{1, 10}
	getInt(t, "array", "0", arrayVal, 1)
	getInt(t, "array", "1", arrayVal, 10)
	getInt(t, "array", "*", arrayVal, 1)
	getAllInt(t, "array", "*", arrayVal, []int{1, 10})
	getFail(t, "array", "A", arrayVal)
	getFail(t, "array", "-1", arrayVal)
	getFail(t, "array", "2", arrayVal)
	getFail(t, "array", "1.A", arrayVal)

	structVal := Struct{100, 1000}
	getInt(t, "struct", "A", structVal, 100)
	getAllInt(t, "struct", "*", structVal, []int{100, 1000})
	getFail(t, "struct", "B", structVal)
	getFail(t, "struct", "C", structVal)
	getFail(t, "struct", "D", structVal)
	getFail(t, "struct", "A.B", structVal)

	getInt(t, "structPtr", "A", &structVal, 100)
	getInt(t, "structPtr", "B", &structVal, 100)
	getInt(t, "structPtr", "C", &structVal, 100)
	getAllInt(t, "structPtr", "*", &structVal, []int{100, 1000})
	getFail(t, "structPtr", "C", &Struct{0, 0})
	getFail(t, "structPtr", "D", &structVal)
	getFail(t, "struct", "A.B", structVal)

	var interfaceVal Interface
	interfaceVal = &Struct{10, 100}
	getInt(t, "interface", "B", interfaceVal, 10)
	getAllInt(t, "interface", "*", interfaceVal, []int{10, 100})

	compound := map[string][]*Struct{
		"X": []*Struct{&Struct{1, 0}},
		"Y": []*Struct{&Struct{2, 0}, &Struct{20, 0}},
		"Z": []*Struct{&Struct{3, 0}, &Struct{30, 0}, &Struct{300, 0}},
	}

	getInt(t, "compound", "X.0.A", compound, 1)
	getInt(t, "compound", "Y.1.A", compound, 20)
	getInt(t, "compound", "Z.2.B", compound, 300)
	getAllInt(t, "compound", "Y.*.A", compound, []int{2, 20})
	getAllInt(t, "compound", "*.0.A", compound, []int{1, 2, 3})
	getAllInt(t, "compound", "*.*.A", compound, []int{1, 2, 3, 20, 30, 300})
	getInterface(t, "compound", "Y", compound)
	getInterface(t, "compound", "Z.2", compound)

	getFail(t, "compound", "X.1.A", compound)
	getFail(t, "compound", "Y.1.D", compound)
	getFail(t, "compound", "W.0.A", compound)
}

func getFail(t *testing.T, title string, path string, obj interface{}) {
	_, err := New(path).Get(obj)
	if err == nil {
		t.Errorf("FAIL(%s): %s -> expected failure", title, path)
	}
}

func getInterface(t *testing.T, title string, path string, obj interface{}) (result interface{}) {
	result, err := New(path).Get(obj)
	if err != nil {
		t.Errorf("FAIL(%s): %s -> %s", title, path, err)
	}
	return
}

func getInt(t *testing.T, title string, path string, obj interface{}, exp int) {
	result := getInterface(t, title, path, obj)
	if result != nil {
		if val := result.(int); val != exp {
			t.Errorf("FAIL(%s): %s -> exp %d got %d", title, path, exp, val)
		}
	}
}

func getAllInterface(t *testing.T, title string, path string, obj interface{}) (result []interface{}) {
	result, err := New(path).GetAll(obj)
	if err != nil {
		t.Errorf("FAIL(%s): %s -> %s", title, path, err)
	}
	return
}

func getAllInt(t *testing.T, title string, path string, obj interface{}, exp []int) {
	result := getAllInterface(t, title, path, obj)
	if result == nil {
		return
	}

	if len(result) != len(exp) {
		t.Errorf("FAIL(%s): %s -> invalid length %d < %d", title, path, len(result), len(exp))
		return
	}

	var intResult []int
	for _, val := range result {
		intResult = append(intResult, val.(int))
	}

	sort.Ints(exp)
	sort.Ints(intResult)

	for i := 0; i < len(intResult); i++ {
		if intResult[i] != exp[i] {
			t.Errorf("FAIL(%s): %s -> exp %d got %d", title, path, exp[i], intResult)
		}
	}
}

func BenchmarkPathNewEmpty(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		New("")
	}
}

func BenchmarkPathNewSimple(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		New("A")
	}
}

func BenchmarkPathNewComplex(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		New("Aasd.*.BJDS.*.*.*.ADJK")
	}
}

func BenchmarkPathGetMap(b *testing.B) {
	n := 26
	obj := make(map[string]int)
	path := make([]P, 26)

	for i := 0; i < n; i++ {
		key := string('A' + i)
		obj[key] = i
		path[i] = New(key)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		path[i%n].Get(obj)
	}
}

func BenchmarkPathGetAllMap(b *testing.B) {
	n := 1000
	obj := make(map[string]int)
	path := New("*")

	for i := 0; i < n; i++ {
		obj[strconv.Itoa(i)] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N/n; i++ {
		path.GetAll(obj)
	}
}

func BenchmarkPathGetArray(b *testing.B) {
	n := 100
	obj := make([]int, n)
	path := make([]P, n)

	for i := 0; i < n; i++ {
		obj[i] = i
		path[i] = New(strconv.Itoa(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		path[i%n].Get(obj)
	}
}

func BenchmarkPathGetAllArray(b *testing.B) {
	obj := make([]int, b.N)
	path := New("*")

	for i := 0; i < b.N; i++ {
		obj[i] = i
	}

	b.ResetTimer()
	path.GetAll(obj)
}

type BenchStruct struct {
	A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Q, R, S, T, U, V, W, X, Y, Z int
}

func BenchmarkPathGetStruct(b *testing.B) {
	n := 26
	obj := BenchStruct{}
	path := make([]P, 26)

	for i := 0; i < n; i++ {
		path[i] = New(string('A' + i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		path[i%n].Get(obj)
	}
}

func BenchmarkPathGetAllStruct(b *testing.B) {
	obj := BenchStruct{}
	path := New("*")

	b.ResetTimer()
	for i := 0; i < b.N/26; i++ {
		path.GetAll(obj)
	}
}

func BenchmarkPathGetStructPtr(b *testing.B) {
	n := 26
	obj := &BenchStruct{}
	path := make([]P, 26)

	for i := 0; i < n; i++ {
		path[i] = New(string('A' + i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		path[i%n].Get(obj)
	}
}

func BenchmarkPathGetAllStructPtr(b *testing.B) {
	obj := &BenchStruct{}
	path := New("*")

	b.ResetTimer()
	for i := 0; i < b.N/26; i++ {
		path.GetAll(obj)
	}
}
