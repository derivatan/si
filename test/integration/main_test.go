//go:build integration

package integration

import (
	"database/sql"
	"fmt"
	"github.com/derivatan/si"
	"testing"
)

var (
	db           *sql.DB
	seededTables []string
)

func TestMain(m *testing.M) {
	var err error
	db, err = sql.Open("postgres", "host=database port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")
	if err != nil {
		panic(err)
	}
	tx, _ := db.Begin()

	db.Query()
	tx.Query()
	db.Exec()
	tx.Exec()
	si.SetLogger(func(a ...any) {
		//fmt.Println(a...)
	})
	m.Run()
}

func DB(t *testing.T) DB {
	tx, err := db.Begin()
	if err != nil {
		t.Fatal("Failed to create database transaction")
	}
	t.Cleanup(func() {
		tx.Rollback()
	})
	return si.NewSQLDB(tx)
}

func Seed[T si.Modeler](tx *sql.Tx, list []T) {
	for _, elem := range list {
		err := si.Save[T](tx, &elem)
		if err != nil {
			panic(fmt.Errorf("Fialed to seed '%T': %w", elem, err))
		}
	}
}

func ResetDB() {
	for _, table := range seededTables {
		_, err := db.Exec(fmt.Sprintf("DELETE FROM %s WHERE 1 = 1", table))
		if err != nil {
			panic(fmt.Errorf("reset DB: %w", err))
		}
	}
	seededTables = []string{}
}
