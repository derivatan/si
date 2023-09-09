//go:build integration

package integration

import (
	"github.com/derivatan/si"
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
	name := "Ray Charles"
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
	name := "Rammstein"
	Seed(db, []artist{
		{Name: name},
		{Name: "Dream Theater"},
	})
	setupObj, setupErr := si.Query[artist](db).Where("name", "=", name).Find()

	obj, err := si.Query[artist](db).Find(*setupObj.ID)

	assert.NoError(t, setupErr)
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

// Builder functions

func TestSelect(t *testing.T) {
	db := DB(t)
	Seed(db, []artist{
		{Name: "314"},
		{Name: "141"},
		{Name: "271"},
	})

	var numResults, max, min int
	selects := []string{
		"COUNT(1)",
		"MIN(name)",
		"MAX(name)",
	}
	_, err := si.Query[artist](db).Select(selects, &numResults, &min, &max).Get()

	assert.NoError(t, err)
	assert.Equal(t, 3, numResults)
	assert.Equal(t, 141, min)
	assert.Equal(t, 314, max)
}

func TestWhere(t *testing.T) {
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
func TestWhereContains(t *testing.T) {
	db := DB(t)
	name := "Beethoven"
	Seed(db, []artist{
		{Name: name},
		{Name: "Mozart"},
	})

	rows, err := si.Query[artist](db).Where("name", "LIKE", "%ee%").Get()

	assert.NoError(t, err)
	assert.Len(t, rows, 1)
	assert.Equal(t, rows[0].Name, name)
}

func TestOrWhere(t *testing.T) {
	db := DB(t)
	name1 := "Prince"
	name2 := "Queen"
	Seed(db, []artist{
		{Name: name1},
		{Name: name2},
		{Name: "Michael Jackson"},
	})

	rows, err := si.Query[artist](db).
		Where("name", "=", name1).
		OrWhere("name", "=", name2).
		OrderBy("name", true).
		Get()

	assert.NoError(t, err)
	assert.Len(t, rows, 2)
	assert.Equal(t, name1, rows[0].Name)
	assert.Equal(t, name2, rows[1].Name)
}

// Test aggregations (with new functions)

// Test types

// Test all query function

// Test save

// Test relations
