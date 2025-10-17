package si

import (
	"reflect"

	"github.com/google/uuid"
)

// HasOne is a relationship where there are ONE other objects (T) that points to this one (F).
// Example:
// | F    |     | T    |
// |------|     |------|
// |      |     | ID   |
// | ID   | <-- | F_ID |
func HasOne[F, T Modeler](model F, refFieldName, fieldName string, relationDataFunc func(f *F) *RelationData[T]) *Relation[F, T] {
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
		relationType: hasOneConf[F, T]{
			idColumn:  column,
			refColumn: refColumn,
			idValue: func(t T) uuid.UUID {
				tVal := reflect.ValueOf(t)
				field := tVal.FieldByName(relationFieldName)
				return field.Interface().(uuid.UUID)
			},
		},
	}
}

type hasOneConf[F, T Modeler] struct {
	idColumn  string
	refColumn string
	idValue   func(T) uuid.UUID
}

func (h hasOneConf[F, T]) collectID(f F) uuid.UUID {
	return *f.GetModel().ID
}

func (h hasOneConf[F, T]) groupBy(t T) uuid.UUID {
	return h.idValue(t)
}

func (h hasOneConf[F, T]) queryColumn() string {
	return h.idColumn
}

func (h hasOneConf[F, T]) joinColumns() (string, string) {
	return "id", h.refColumn
}
