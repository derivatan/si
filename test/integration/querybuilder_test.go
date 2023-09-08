//go:build integration

package integration

import (
	"fmt"
	"github.com/derivatan/si"
	"github.com/gofrs/uuid"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	db := DB(t)
	name := "Pink Floyd"
	Seed(db, []artist{
		{Name: name},
	})

	list, err := si.Query[artist](db).Get()

	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, name, list[0].Name)
}

func TestFirst(t *testing.T) {
	db := DB(t)
	name := "Michael Jackson"
	Seed(db, []artist{
		{Name: name},
		{Name: "Stevie Wonder"},
	})

	obj, err := si.Query[artist](db).OrderBy("name", true).First()

	assert.NoError(t, err)
	assert.NotNil(t, obj)
	assert.Equal(t, name, obj.Name)
}

func TestFind(t *testing.T) {
	db := DB(t)
	name := "Portishead"
	Seed(db, []artist{
		{Name: name},
	})

	obj, err := si.Query[artist](db).Where("name", "=", name).Find()

	assert.NoError(t, err)
	assert.NotNil(t, obj)
	assert.Equal(t, name, obj.Name)
}

func TestFindWithID(t *testing.T) {
	db := DB(t)
	id := uuid.FromStringOrNil("12341234-1234-1234-1234-123412341234")
	name := "Rammstein"
	Seed(db, []artist{
		{Model: si.Model{ID: &id}, Name: name},
		{Name: "Dream Theater"},
	})

	for _, i2 := range si.Query[artist](db).MustGet() {
		fmt.Println(i2)
	}
	obj, err := si.Query[artist](db).Find(id)

	assert.NoError(t, err)
	assert.NotNil(t, obj)
	assert.Equal(t, name, obj.Name)
}

func TestWithWrongNumberOfResults(t *testing.T) {
	db := DB(t)
	Seed(db, []artist{
		{Name: "Eminem"},
		{Name: "The Beatles"},
	})

	obj, err := si.Query[artist](db).Find()

	assert.Nil(t, obj)
	assert.Error(t, err)
}

func TestGetArtists(t *testing.T) {
	db := DB(t)
	Seed(db, []artist{
		{Name: "Beethoven"},
		{Name: "Mozart"},
	})

	rows, err := si.Query[artist](db).Get()

	assert.NoError(t, err)
	assert.Len(t, rows, 2)
}

func TestGetWhere(t *testing.T) {
	db := DB(t)
	wantedName := "Second"
	Seed(db, []artist{
		{Name: "First"},
		{Name: wantedName},
		{Name: "Third"},
	})

	rows, err := si.Query[artist](db).Where("name", "=", wantedName).Get()

	assert.NoError(t, err)
	assert.Len(t, rows, 1)
	assert.Equal(t, rows[0].Name, wantedName)
}

// Test aggregations (with new functions)

// Test types

// Test all query function

// Test save

// Test relations
