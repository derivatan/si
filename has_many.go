package si

import (
	"github.com/google/uuid"
	"reflect"
)

// HasMany is a relationship where there are MULTIPLE other objects (T) that points to this one (F).
// Example:
// | F    |     | T    |
// |------|     |------|
// |      |     | ID   |
// | ID   | <-- | F_ID |
func HasMany[F, T Modeler](model F, refFieldName, fieldName string, relationDataFunc func(f *F) *RelationData[T]) *Relation[F, T] {
	fromType := reflect.TypeOf(new(F))
	toType := reflect.TypeOf(new(T))
	relationFieldName := getRelationFieldName(fromType, toType, fieldName, false)
	column := getColumnNameString(toType, relationFieldName)
	refColumn := getColumnNameString(toType, refFieldName)

	return &Relation[F, T]{
		model:        model,
		query:        Query[T](),
		get:          Query[T]().Where(column, "=", model.GetModel().ID),
		relationData: relationDataFunc,
		relationType: hasManyConf[F, T]{
			idField:   column,
			refColumn: refColumn,
			idValue: func(t T) uuid.UUID {
				tVal := reflect.ValueOf(t)
				return tVal.FieldByName(relationFieldName).Interface().(uuid.UUID)
			},
		},
	}
}

type hasManyConf[F, T Modeler] struct {
	idField   string
	refColumn string
	idValue   func(a T) uuid.UUID
}

func (h hasManyConf[F, T]) collectID(f F) uuid.UUID {
	return *f.GetModel().ID
}

func (h hasManyConf[F, T]) groupBy(t T) uuid.UUID {
	return h.idValue(t) // This should not be the id. It should be the referal object.
}

func (h hasManyConf[F, T]) queryColumn() string {
	return h.idField
}

func (h hasManyConf[F, T]) joinColumns() (string, string) {
	return "id", h.refColumn
}
