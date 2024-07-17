package si

import (
	"github.com/google/uuid"
	"reflect"
)

// BelongsTo is a relationship where the object in question (F) has a reference to the other object (T)
// Example:
// | F    |     | T  |
// |------|     |----|
// | ID   |     |    |
// | T_ID | --> | ID |
func BelongsTo[F, T Modeler](model F, refFieldName, fieldName string, relationDataFunc func(f *F) *RelationData[T]) *Relation[F, T] {
	fromType := reflect.TypeOf(new(F))
	toType := reflect.TypeOf(new(T))
	relationFieldName := getRelationFieldName(fromType, toType, fieldName, true)
	refColumn := getColumnNameString(fromType, refFieldName)
	idField := func(f F) uuid.UUID {
		val := reflect.ValueOf(f)
		field := val.FieldByName(relationFieldName)
		return field.Interface().(uuid.UUID)
	}

	return &Relation[F, T]{
		model:        model,
		query:        Query[T](),
		get:          Query[T]().Where("id", "=", idField(model)),
		relationData: relationDataFunc,
		relationType: belongsToConf[F, T]{
			refColumn: refColumn,
			idField:   idField,
		},
	}
}

type belongsToConf[F, T Modeler] struct {
	refColumn string
	idField   func(F) uuid.UUID
}

func (b belongsToConf[F, T]) collectID(f F) uuid.UUID {
	return b.idField(f)
}

func (b belongsToConf[F, T]) groupBy(t T) uuid.UUID {
	return *t.GetModel().ID
}

func (b belongsToConf[F, T]) queryColumn() string {
	return "id"
}

func (b belongsToConf[F, T]) joinColumns() (string, string) {
	return b.refColumn, "id"
}
