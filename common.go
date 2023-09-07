package si

import (
	"errors"
	"fmt"
	"github.com/gofrs/uuid"
	"reflect"
	"regexp"
	"strings"
	"time"
)

// DB interface is based on the `sql.DB` and `sql.Tx.
type DB interface {
	Exec(query string, args ...any) (any, error)
	Query(query string, args ...any) (Rows, error)
}

type Rows interface {
	Next() bool
	Scan(dest ...any) error
	Close() error
}

type Model struct {
	ID        *uuid.UUID `si:"id"`
	CreatedAt time.Time  `si:"created_at"`
	UpdatedAt *time.Time `si:"updated_at"`
	DeletedAt *time.Time `si:"deleted_at"`
}

type Modeler interface {
	GetModel() Model
	GetTable() string
}

var (
	config                secretIngredientConfig
	ResourceNotFoundError = errors.New("resource not found")
)

type secretIngredientConfig struct {
	logger  func(a ...any)
	configs map[string]any // 'any' must be `ModelConfig[T]`
}

type ModelConfig[T Modeler] struct {
	Table string
}

// SetLogger will set a logger function for debugging all queries.
func SetLogger(f func(a ...any)) {
	config.logger = f
}

// Query will start a query.
// Main starting point for retrieving objects.
func Query[T Modeler](db DB) *QueryBuilder[T] {
	return &QueryBuilder[T]{
		db:      db,
		filters: []filter{},
		orderBy: []orderBy{},
	}
}

// Save a model to the database.
func Save[T Modeler](db DB, m *T) error {
	return save[T](db, m)
}

func log(s ...any) {
	if config.logger != nil {
		config.logger(s)
	}
}

type typeInfo struct {
	Columns []string
	Names   []string
	Values  []any
}

func getTypeInfo(obj any) typeInfo {
	result := typeInfo{}
	refType := reflect.TypeOf(obj).Elem()
	refVal := reflect.ValueOf(obj).Elem()
	for i := 0; i < refVal.NumField(); i++ {
		fieldType := refType.Field(i)
		fieldVal := refVal.Field(i)
		if !fieldType.IsExported() {
			continue
		}

		if i == 0 {
			for j := 0; j < fieldType.Type.NumField(); j++ {
				modelType := fieldType.Type.Field(j)
				modelVal := fieldVal.Field(j)
				result = appendResult(modelType, modelVal, result)
			}
			continue
		}
		result = appendResult(fieldType, fieldVal, result)
	}

	return result
}

func appendResult(t reflect.StructField, v reflect.Value, ti typeInfo) typeInfo {
	column := getColumnName(t)
	value := v.Addr().Interface()
	ti.Columns = append(ti.Columns, column)
	ti.Names = append(ti.Names, t.Name)
	ti.Values = append(ti.Values, value)
	return ti
}

func toSnakeCase(str string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func getColumnNameString(t reflect.Type, fieldName string) string {
	field, ok := t.Elem().FieldByName(fieldName)
	if !ok {
		panic(fmt.Sprintf("'%s' is not a field on model '%s'", fieldName, t.Elem().Name()))
	}
	return getColumnName(field)
}

func getColumnName(field reflect.StructField) string {
	if siTag, ok := field.Tag.Lookup("si"); ok {
		return siTag
	}
	return toSnakeCase(field.Name)
}

func getRelationFieldName(f reflect.Type, t reflect.Type, fieldName string, fieldOnTo bool) string {
	field, ok := f.Elem().FieldByName(fieldName)
	if !ok {
		panic(fmt.Sprintf("'%s' is not a field on model '%s'", fieldName, t.Elem().Name()))
	}
	referenceField, ok := field.Tag.Lookup("si")
	if !ok {
		s := f.String()
		if fieldOnTo {
			s = t.String()
		}
		result := strings.Split(s, ".")[1] + "ID"
		return strings.ToUpper(result[:1]) + result[1:]
	}
	return referenceField
}
