package util_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thatmattlove/oui/v2/internal/util"
)

func Test_removeComments(t *testing.T) {
	t.Run("removeComments - inline", func(t *testing.T) {
		assert.Equal(t, "stuff", util.RemoveComments("stuff # with comment"))
	})
	t.Run("removeComments - full line", func(t *testing.T) {
		assert.Equal(t, "", util.RemoveComments("# line with comment"))
	})
	t.Run("removeComments - escaped", func(t *testing.T) {
		assert.Equal(t, `escaped \# comment`, util.RemoveComments(`escaped \# comment`))
	})
	t.Run("removeComments - empty line", func(t *testing.T) {
		assert.Equal(t, "", util.RemoveComments("#"))
	})
	t.Run("removeComments - end comment", func(t *testing.T) {
		r := util.RemoveComments(`00:00:D1	Adaptec		Adaptec, Inc.	# "Nodem" product`)
		e := `00:00:D1	Adaptec		Adaptec, Inc.`
		assert.Equal(t, e, r)
	})
}

func Test_splitTabs(t *testing.T) {
	t.Run("splitTabs 1", func(t *testing.T) {
		r := util.SplitTabs(`one	two	three`)
		e := []string{"one", "two", "three"}
		assert.Equal(t, e, r)
	})
	t.Run("splitTabs 2", func(t *testing.T) {
		r := util.SplitTabs(`one	two		three`)
		e := []string{"one", "two", "three"}
		assert.Equal(t, e, r)
	})
	t.Run("splitTabs 3", func(t *testing.T) {
		r := util.SplitTabs(`one		two		three`)
		e := []string{"one", "two", "three"}
		assert.Equal(t, e, r)
	})
	t.Run("splitTabs 4", func(t *testing.T) {
		r := util.SplitTabs(`one		two		three	`)
		e := []string{"one", "two", "three"}
		assert.Equal(t, e, r)
	})
}

func Test_pathExists(t *testing.T) {
	t.Run("pathExists - not exists", func(t *testing.T) {
		r := util.PathExists("/this/path/does/not/exist")
		assert.False(t, r)
	})
	t.Run("pathExists - exists", func(t *testing.T) {
		tf, err := os.CreateTemp(os.TempDir(), "oui-test-*")
		if err != nil {
			panic(err)
		}
		defer os.Remove(tf.Name())

		r := util.PathExists(tf.Name())
		assert.True(t, r)
	})
}

func Test_SplitSlice(t *testing.T) {
	t.Run("equal parts", func(t *testing.T) {
		t.Parallel()
		original := []int{1, 2, 3, 4, 5, 6}
		max := 2
		expected := [][]int{
			{1, 2},
			{3, 4},
			{5, 6},
		}
		result := util.SplitSlice(original, max)
		assert.True(t, reflect.DeepEqual(expected, result))
	})
	t.Run("unequal parts", func(t *testing.T) {
		t.Parallel()
		original := []int{1, 2, 3, 4, 5, 6}
		max := 4
		expected := [][]int{
			{1, 2, 3, 4},
			{5, 6},
		}
		result := util.SplitSlice(original, max)
		assert.True(t, reflect.DeepEqual(expected, result))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		original := []string{}
		max := 3
		expected := [][]string{}
		result := util.SplitSlice(original, max)
		assert.Equal(t, expected, result)
	})

}
