package oui_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thatmattlove/oui/v2/oui"
)

func Test_New(t *testing.T) {
	mac := "00:50:56:00:b3:3f"
	prefix := "00:50:56:00:00:00/24"
	org := "VMware, Inc."
	registry := "MA-L"

	t.Run("postgres", func(t *testing.T) {
		t.Parallel()
		password := os.Getenv("POSTGRES_PASSWORD")
		require.NotEqual(t, "", password, "missing POSTGRES_PASSWORD environment variable")
		cs := fmt.Sprintf("postgresql://oui:%s@localhost/oui?sslmode=disable", password)
		psql, err := oui.CreatePostgresOption(cs)
		require.NoError(t, err)
		ouidb, err := oui.New(oui.WithVersion("test"), psql)
		require.NoError(t, err)
		t.Cleanup(func() {
			ouidb.Clear()
		})
		t.Run("populate", func(t *testing.T) {
			populated, err := ouidb.Populate()
			require.NoError(t, err)
			assert.NotZero(t, populated)
		})
		t.Run("count", func(t *testing.T) {
			count, err := ouidb.Count()
			require.NoError(t, err)
			assert.NotZero(t, count)
		})
		t.Run("query", func(t *testing.T) {
			matches, err := ouidb.Find(mac)
			require.NoError(t, err)
			assert.NotZero(t, len(matches))
			match := matches[0]
			assert.Equal(t, org, match.Org)
			assert.Equal(t, prefix, match.Prefix)
			assert.Equal(t, registry, match.Registry)
		})
	})
	t.Run("sqlite", func(t *testing.T) {
		t.Parallel()
		d := t.TempDir()
		f := filepath.Join(d, "oui.db")
		sqlite, err := oui.CreateSQLiteOption(f)
		require.NoError(t, err)
		ouidb, err := oui.New(oui.WithVersion("test"), sqlite)
		require.NoError(t, err)
		t.Run("populate", func(t *testing.T) {
			populated, err := ouidb.Populate()
			require.NoError(t, err)
			assert.NotZero(t, populated)
		})
		t.Run("count", func(t *testing.T) {
			count, err := ouidb.Count()
			require.NoError(t, err)
			assert.NotZero(t, count)
		})
		t.Run("query", func(t *testing.T) {
			matches, err := ouidb.Find(mac)
			require.NoError(t, err)
			assert.NotZero(t, len(matches))
			match := matches[0]
			assert.Equal(t, org, match.Org)
			assert.Equal(t, prefix, match.Prefix)
			assert.Equal(t, registry, match.Registry)
		})
	})
}
