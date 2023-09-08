//go:build integration

package integration

import (
	"database/sql"
	"fmt"
	"github.com/derivatan/si"
	"testing"
)

var (
	db *sql.DB
)

func TestMain(m *testing.M) {
	var err error
	db, err = sql.Open("postgres", "host=database port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")
	if err != nil {
		panic(err)
	}

	si.SetLogger(func(a ...any) {
		//fmt.Println(a...)
	})
	m.Run()
}

func DB(t *testing.T) si.DB {
	tx, err := db.Begin()
	if err != nil {
		t.Fatal("Failed to create database transaction")
	}
	t.Cleanup(func() {
		err := tx.Rollback()
		if err != nil {
			t.Fatal("Failed to rollback database transaction")
		}
	})
	return si.WrapDB(tx)
}

func Seed[T si.Modeler](tx si.DB, list []T) {
	for _, elem := range list {
		err := si.Save[T](tx, &elem)
		if err != nil {
			panic(fmt.Errorf("Fialed to seed '%T': %w", elem, err))
		}
	}
}
