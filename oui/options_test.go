package oui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_sanitizeVersion(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Parallel()
		ver := "db1"
		result := sanitizeVersion(ver)
		assert.Equal(t, ver, result)
	})
	t.Run("starting number", func(t *testing.T) {
		t.Parallel()
		ver := "1oui"
		result := sanitizeVersion(ver)
		assert.Equal(t, "oui_1oui", result)
	})
	t.Run("extra chars", func(t *testing.T) {
		t.Parallel()
		ver := "oui&db*123"
		result := sanitizeVersion(ver)
		assert.Equal(t, "oui__db__123", result)
	})
	t.Run("semver", func(t *testing.T) {
		t.Parallel()
		ver := "2.0.3"
		result := sanitizeVersion(ver)
		assert.Equal(t, "oui_2__0__3", result)
	})
	t.Run("numbers only", func(t *testing.T) {
		t.Parallel()
		ver := "12345"
		result := sanitizeVersion(ver)
		assert.Equal(t, "oui_12345", result)
	})
}
