// Copyright (c) 2014 Datacratic. All rights reserved.

package path

import (
	"reflect"
	"testing"
)

func TestPathTranslate(t *testing.T) {
	obj := struct {
		A *struct {
			B int `json:"bob"`
		} `json:"alice,omitempty"`
		C []*struct {
			D map[string]*struct {
				E int `json:"eve"`
			} `json:"dan"`
		} `json:"charlie"`
	}{}

	aliases := JSONAliases(reflect.TypeOf(obj))

	translate(t, aliases, "alice.bob", "A.B")
	translate(t, aliases, "charlie.dan.eve", "C.D.E")
	translate(t, aliases, "charlie.*.eve", "C.*.E")
	translate(t, aliases, "zebra.alice.wall", "zebra.A.wall")
}

func translate(t *testing.T, aliases map[string]string, old, exp string) {
	path := New(old).Translate(aliases).String()
	if path != exp {
		t.Errorf("FAIL: %s -> %s != %s ", old, path, exp)
	}
}
