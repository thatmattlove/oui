package main

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_removeComments(t *testing.T) {

	tests := [][]string{
		{"stuff # with comment", "stuff"},
		{"# line with comment", ""},
		{`escaped \# comment`, `escaped \# comment`},
		{"#", ""},
		{`00:00:D1	Adaptec		Adaptec, Inc.	# "Nodem" product`, `00:00:D1	Adaptec		Adaptec, Inc.`},
	}
	for i, p := range tests {
		t.Run(fmt.Sprintf("removeComments %d", i+1), func(t *testing.T) {
			er := p[1]
			r := removeComments(p[0])
			if r != er {
				t.Errorf("'%v' != '%v'", r, er)
			}
		})
	}
}

func Test_splitTabs(t *testing.T) {
	t.Run("splitTabs 1", func(t *testing.T) {
		r := splitTabs(`one	two	three`)
		e := []string{"one", "two", "three"}
		if !reflect.DeepEqual(r, e) {
			t.Errorf("'%v' != '%v'", r, e)
		}
	})
	t.Run("splitTabs 2", func(t *testing.T) {
		r := splitTabs(`one	two		three`)
		e := []string{"one", "two", "three"}
		if !reflect.DeepEqual(r, e) {
			t.Errorf("'%v' != '%v'", r, e)
		}
	})
	t.Run("splitTabs 3", func(t *testing.T) {
		r := splitTabs(`one		two		three`)
		e := []string{"one", "two", "three"}
		if !reflect.DeepEqual(r, e) {
			t.Errorf("'%v' != '%v'", r, e)
		}
	})
	t.Run("splitTabs 4", func(t *testing.T) {
		r := splitTabs(`one		two		three	`)
		e := []string{"one", "two", "three"}
		if !reflect.DeepEqual(r, e) {
			t.Errorf("'%v' != '%v'", r, e)
		}
	})
}
