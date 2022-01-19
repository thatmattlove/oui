package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_removeComments(t *testing.T) {
	t.Run("removeComments - inline", func(t *testing.T) {
		assert.Equal(t, "stuff", removeComments("stuff # with comment"))
	})
	t.Run("removeComments - full line", func(t *testing.T) {
		assert.Equal(t, "", removeComments("# line with comment"))
	})
	t.Run("removeComments - escaped", func(t *testing.T) {
		assert.Equal(t, `escaped \# comment`, removeComments(`escaped \# comment`))
	})
	t.Run("removeComments - empty line", func(t *testing.T) {
		assert.Equal(t, "", removeComments("#"))
	})
	t.Run("removeComments - end comment", func(t *testing.T) {
		r := removeComments(`00:00:D1	Adaptec		Adaptec, Inc.	# "Nodem" product`)
		e := `00:00:D1	Adaptec		Adaptec, Inc.`
		assert.Equal(t, e, r)
	})
}

func Test_splitTabs(t *testing.T) {
	t.Run("splitTabs 1", func(t *testing.T) {
		r := splitTabs(`one	two	three`)
		e := []string{"one", "two", "three"}
		assert.Equal(t, e, r)
	})
	t.Run("splitTabs 2", func(t *testing.T) {
		r := splitTabs(`one	two		three`)
		e := []string{"one", "two", "three"}
		assert.Equal(t, e, r)
	})
	t.Run("splitTabs 3", func(t *testing.T) {
		r := splitTabs(`one		two		three`)
		e := []string{"one", "two", "three"}
		assert.Equal(t, e, r)
	})
	t.Run("splitTabs 4", func(t *testing.T) {
		r := splitTabs(`one		two		three	`)
		e := []string{"one", "two", "three"}
		assert.Equal(t, e, r)
	})
}

func Test_pathExists(t *testing.T) {
	t.Run("pathExists - not exists", func(t *testing.T) {
		r := pathExists("/this/path/does/not/exist")
		assert.False(t, r)
	})
	t.Run("pathExists - exists", func(t *testing.T) {
		tf, err := os.CreateTemp(os.TempDir(), "oui-test-*")
		if err != nil {
			panic(err)
		}
		defer os.Remove(tf.Name())

		r := pathExists(tf.Name())
		assert.True(t, r)
	})
}
