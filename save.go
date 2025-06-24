package si

import (
	"fmt"
	"github.com/google/uuid"
	"reflect"
	"slices"
	"strings"
	"time"
)

func save[T Modeler](db DB, m *T, fields []string) error {
	now := time.Now()
	// Updated at
	reflect.ValueOf(m).Elem().Field(0).Field(2).Set(reflect.ValueOf(&now))

	if (*m).GetModel().ID == nil {
		// Created at
		reflect.ValueOf(m).Elem().Field(0).Field(1).Set(reflect.ValueOf(now))

		return insert[T](db, m)
	} else {
		return update[T](db, m, fields)
	}
}

func insert[T Modeler](db DB, m *T) error {
	ti := getTypeInfo(m)

	query, parameters := buildInsert[T](ti)
	log(query, parameters)

	rows, err := db.Query(query, parameters...)
	if err != nil {
		return fmt.Errorf("si.insert: execute query: %w", err)
	}
	defer func() { _ = rows.Close() }()

	rows.Next()
	err = rows.Scan(ti.Values[0])
	if err != nil {
		return fmt.Errorf("si.insert: scan: %w", err)
	}

	return nil
}

func buildInsert[T Modeler](ti typeInfo) (string, []any) {
	var values []string
	var parameters []any
	columns := ti.Columns[1:]

	for i := 1; i < len(ti.Values); i++ {
		values = append(values, fmt.Sprintf("$%d", i))
		parameters = append(parameters, ti.Values[i])
	}

	if ti.Values[0] != nil {
		apa := ti.Values[0].(**uuid.UUID)
		if *apa != nil {
			values = append(values, fmt.Sprintf("$%d", len(ti.Values)))
			parameters = append(parameters, *apa)
			columns = append(columns, "id")
		}
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) RETURNING id",
		(*new(T)).GetTable(),
		strings.Join(columns, ","),
		strings.Join(values, ","),
	)
	return query, parameters
}

func update[T Modeler](db DB, m *T, fields []string) error {
	ti := getTypeInfo(m)
	query, parameters := buildUpdate[T](ti, fields)
	log(query, parameters)

	_, err := db.Exec(query, parameters...)
	if err != nil {
		return fmt.Errorf("si.update: execute query: %w", err)
	}
	return nil
}

func buildUpdate[T Modeler](ti typeInfo, fields []string) (string, []any) {

	var columns []string
	var parameters []any
	var parameterCount = 1
	for i := 1; i < len(ti.Columns); i++ {
		if fields != nil && !slices.Contains(fields, ti.Columns[i]) {
			continue
		}
		columns = append(columns, fmt.Sprintf("%s=$%d", ti.Columns[i], parameterCount))
		parameters = append(parameters, ti.Values[i])
		parameterCount++
	}

	parameters = append(parameters, ti.Values[0])
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d",
		(*new(T)).GetTable(),
		strings.Join(columns, ","),
		parameterCount,
	)
	return query, parameters
}
