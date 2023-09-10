//go:build integration

package integration

import (
	"fmt"
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

func TestWhereF(t *testing.T) {
	db := DB(t)
	Seed(db, []artist{
		{Name: "Danny Elfman"},
		{Name: "Hans Zimmer"},
		{Name: "John Williams"},
	})

	// WHERE a AND (b OR c)
	rows, err := si.Query[artist](db).Where("name", "ILIKE", "%m%").WhereF(func(q *si.QueryBuilder[artist]) *si.QueryBuilder[artist] {
		return q.Where("name", "ILIKE", "%Zi%").OrWhere("name", "ILIKE", "%Wi%")
	}).Get()

	assert.NoError(t, err)
	assert.Len(t, rows, 2)
}

func TestOrWhereF(t *testing.T) {
	db := DB(t)
	Seed(db, []artist{
		{Name: "Björk"},
		{Name: "Daft Punk"},
		{Name: "The Knife"},
	})

	// WHERE a OR (b AND c)
	rows, err := si.Query[artist](db).Where("name", "ILIKE", "%knife%").OrWhereF(func(q *si.QueryBuilder[artist]) *si.QueryBuilder[artist] {
		return q.Where("name", "ILIKE", "%daft%").Where("name", "ILIKE", "%punk%")
	}).Get()

	assert.NoError(t, err)
	assert.Len(t, rows, 2)
}

func TestOrderBy(t *testing.T) {
	db := DB(t)
	nameA := "Avalanches, The"
	nameB := "Basement Jaxx"
	nameC := "Cure, The"
	nameD := "Deep Purple"
	Seed(db, []artist{
		{Name: nameB},
		{Name: nameC},
		{Name: nameA},
		{Name: nameD},
	})

	rowsAsc, errAsc := si.Query[artist](db).OrderBy("name", true).Get()
	rowsDesc, errDesc := si.Query[artist](db).OrderBy("name", false).Get()

	assert.NoError(t, errAsc)
	assert.NoError(t, errDesc)
	assert.Equal(t, nameA, rowsAsc[0].Name)
	assert.Equal(t, nameB, rowsAsc[1].Name)
	assert.Equal(t, nameC, rowsAsc[2].Name)
	assert.Equal(t, nameD, rowsAsc[3].Name)
	assert.Equal(t, nameD, rowsDesc[0].Name)
	assert.Equal(t, nameC, rowsDesc[1].Name)
	assert.Equal(t, nameB, rowsDesc[2].Name)
	assert.Equal(t, nameA, rowsDesc[3].Name)
}

func TestTakeAndSkip(t *testing.T) {
	db := DB(t)
	name1 := "Detektivbyrån"
	name2 := "Trazan & Banarne"
	name3 := "Electric Banana Band"
	Seed(db, []artist{
		{Name: name1},
		{Name: name2},
		{Name: name3},
	})

	rowsTake, errTake := si.Query[artist](db).OrderBy("name", true).Take(2).Get()
	rowsSkip, errSkip := si.Query[artist](db).OrderBy("name", true).Skip(1).Get()

	assert.NoError(t, errTake)
	assert.NoError(t, errSkip)
	assert.Len(t, rowsTake, 2)
	assert.Equal(t, name1, rowsTake[0].Name)
	assert.Equal(t, name3, rowsTake[1].Name)
	assert.Len(t, rowsSkip, 2)
	assert.Equal(t, name3, rowsSkip[0].Name)
	assert.Equal(t, name2, rowsSkip[1].Name)
}

func TestGroupBy(t *testing.T) {
	db := DB(t)
	Seed(db, []contact{
		{Email: "info@email.com", Phone: 101},
		{Email: "info@email.com", Phone: 103},
		{Email: "support@email.com", Phone: 107},
		{Email: "support@email.com", Phone: 109},
	})

	selects := []string{"email", "SUM(phone)"}
	type result struct {
		email string
		sum   int
	}
	a := func() (any, []any) {
		a := result{}
		return a, []any{&a.email, &a.sum}
	}
	// TODO: How do i get the list in the end?
	_, err := si.Query[contact](db).GroupSelect(selects, a).GroupBy("email").Get()

	fmt.Println("This is the values.")

	assert.NoError(t, err)

	// TODO: test having here aswell.
}

// Test data-types on structs. bool, int, time, duration, json...

// Test save
// Test relations
