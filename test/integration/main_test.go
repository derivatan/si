//go:build integration

package integration

import (
	"database/sql"
	"fmt"
	"github.com/derivatan/si"
	"github.com/gofrs/uuid"
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

// DB returns a transaction, that will rollback when the test is finished.
func DB(t *testing.T) si.DB {
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Failed to create database transaction: %v", err)
	}
	t.Cleanup(func() {
		err := tx.Rollback()
		if err != nil {
			t.Fatalf("Failed to rollback database transaction: %v", err)
		}
	})
	return si.WrapDB(tx)
}

func Seed[T si.Modeler](tx si.DB, list []T) []uuid.UUID {
	var result []uuid.UUID
	for _, elem := range list {
		err := si.Save[T](tx, &elem)
		if err != nil {
			panic(fmt.Errorf("Fialed to seed '%T': %w", elem, err))
		}
		result = append(result, *elem.GetModel().ID)
	}
	return result
}
