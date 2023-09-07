//go:build integration

package integration

import (
	"github.com/derivatan/si"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	t.Cleanup(ResetDB)
	name := "Pink FLoyd"
	Seed([]artist{
		{Name: name},
	})

	list, err := si.Query[artist]().Get()

	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, name, list[0].Name)
}

func TestFirst(t *testing.T) {
	t.Cleanup(ResetDB)
	name := "Michael Jackson"
	Seed([]artist{
		{Name: name},
		{Name: "Stevie Wonder"},
	})

	obj, err := si.Query[artist]().OrderBy("name", true).First()

	assert.NoError(t, err)
	assert.NotNil(t, obj)
	assert.Equal(t, name, obj.Name)
}

func TestFind(t *testing.T) {
	t.Cleanup(ResetDB)
	Seed("artists", []map[string]any{
		{"name": "something"},
	})
}

func TestFindWithID(t *testing.T) {

}

func TestWithWrongResult(t *testing.T) {

}

func TestGetArtists(t *testing.T) {
	t.Cleanup(ResetDB)
	Seed("artists", []map[string]any{
		{"name": "Beethoven"},
		{"name": "Mozart"},
	})

	rows, err := si.Query[artist]().Get()

	assert.NoError(t, err)
	assert.Len(t, rows, 2)
}

func TestGetWhere(t *testing.T) {
	t.Cleanup(ResetDB)
	wantedName := "Second"
	Seed("artists", []map[string]any{
		{"name": "First"},
		{"name": wantedName},
		{"name": "Third"},
	})

	rows, err := si.Query[artist]().Where("name", "=", wantedName).Get()

	assert.NoError(t, err)
	assert.Len(t, rows, 1)
	assert.Equal(t, rows[0].Name, wantedName)
}

// Test aggregations (with new functions)

// Test types

// Test all query function

// Test save

// Test relations
