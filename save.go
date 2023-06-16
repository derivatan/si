package si

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

func save[T Modeler](m *T) error {
	now := time.Now()
	// Updated at
	reflect.ValueOf(m).Elem().Field(0).Field(2).Set(reflect.ValueOf(&now))

	if (*m).GetModel().ID == nil {
		// Created at
		reflect.ValueOf(m).Elem().Field(0).Field(1).Set(reflect.ValueOf(now))

		return insert[T](m)
	} else {
		return update[T](m)
	}
}

func insert[T Modeler](m *T) error {
	ti := getTypeInfo(m)

	query, parameters := buildInsert[T](ti)
	log(query, parameters)

	rows, err := config.db.Query(query, parameters...)
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
	conf := getModelConf[T]()
	var values []string
	var parameters []any

	for i := 1; i < len(ti.Values); i++ {
		values = append(values, fmt.Sprintf("$%d", i))
		parameters = append(parameters, ti.Values[i])
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) RETURNING id",
		conf.Table,
		strings.Join(ti.Columns[1:], ","),
		strings.Join(values, ","),
	)
	return query, parameters
}

func update[T Modeler](m *T) error {
	ti := getTypeInfo(m)
	query, parameters := buildUpdate[T](ti)
	log(query, parameters)

	_, err := config.db.Exec(query, parameters...)
	if err != nil {
		return fmt.Errorf("si.update: execute query: %w", err)
	}
	return nil
}

func buildUpdate[T Modeler](ti typeInfo) (string, []any) {
	conf := getModelConf[T]()

	var columns []string
	var parameters []any
	for i := 1; i < len(ti.Columns); i++ {
		columns = append(columns, fmt.Sprintf("%s=$%d", ti.Columns[i], i))
		parameters = append(parameters, ti.Values[i])
	}

	parameters = append(parameters, ti.Values[0])
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d",
		conf.Table,
		strings.Join(columns, ","),
		len(ti.Columns),
	)
	return query, parameters
}
