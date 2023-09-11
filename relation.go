package si

import (
	"github.com/gofrs/uuid"
)

// Relation is the configuration for a relationship between two objects.
// F is `from` and T is `to`, so the relation is defined as seen from ´F´.
type Relation[F, T Modeler] struct {
	model        F
	query        *QueryBuilder[T]
	get          *QueryBuilder[T]
	relationData func(f *F) *RelationData[T]
	relationType relationType[F, T]
}

type RelationData[T Modeler] struct {
	loaded bool
	data   []T
}

func (r *Relation[F, T]) Unload() *Relation[F, T] {
	rd := r.relationData(&r.model)
	rd.loaded = false
	rd.data = nil
	return r
}

func (r *Relation[F, T]) Get(db DB) ([]T, error) {
	rd := r.relationData(&r.model)
	if rd.loaded {
		return rd.data, nil
	}
	return r.innerFilter().Get(db)
}

func (r *Relation[F, T]) First(db DB) (*T, error) {
	rd := r.relationData(&r.model)
	if rd.loaded {
		return &rd.data[0], nil
	}
	return r.innerFilter().First(db)
}

func (r *Relation[F, T]) Find(db DB, id ...uuid.UUID) (*T, error) {
	rd := r.relationData(&r.model)
	if rd.loaded {
		return &rd.data[0], nil
	}
	return r.innerFilter().Find(db, id...)
}

func (r *Relation[F, T]) MustGet(db DB) []T {
	result, err := r.Get(db)
	if err != nil {
		panic(err)
	}
	return result
}

func (r *Relation[F, T]) MustFirst(db DB) *T {
	result, err := r.First(db)
	if err != nil {
		panic(err)
	}
	return result
}

func (r *Relation[F, T]) MustFind(db DB, id ...uuid.UUID) *T {
	result, err := r.Find(db, id...)
	if err != nil {
		panic(err)
	}
	return result
}

func (r *Relation[F, T]) innerFilter() *QueryBuilder[T] {
	result := r.get
	if len(r.query.filters) > 0 {
		result = result.WhereF(func(q *QueryBuilder[T]) *QueryBuilder[T] {
			return r.query
		})
	}
	return result
}

func (r *Relation[F, T]) Where(column, op string, value any) *Relation[F, T] {
	r.query = r.query.Where(column, op, value)
	return r
}

func (r *Relation[F, T]) OrWhere(column, op string, value any) *Relation[F, T] {
	r.query = r.query.OrWhere(column, op, value)
	return r
}

func (r *Relation[F, T]) WhereF(f func(q *QueryBuilder[T]) *QueryBuilder[T]) *Relation[F, T] {
	r.query = r.query.WhereF(f)
	return r
}

func (r *Relation[F, T]) OrWhereF(f func(q *QueryBuilder[T]) *QueryBuilder[T]) *Relation[F, T] {
	r.query = r.query.OrWhereF(f)
	return r
}

func (r *Relation[F, T]) OrderBy(column string, asc bool) *Relation[F, T] {
	r.query = r.query.OrderBy(column, asc)
	return r
}

func (r *Relation[F, T]) Take(number int) *Relation[F, T] {
	r.query = r.query.Take(number)
	return r
}

func (r *Relation[F, T]) Skip(number int) *Relation[F, T] {
	r.query = r.query.Skip(number)
	return r
}

func (r *Relation[F, T]) With(f func(m T, r []T) error) *Relation[F, T] {
	r.query = r.query.With(f)
	return r
}

func (r *Relation[F, T]) WithDeleted() *Relation[F, T] {
	r.query = r.query.WithDeleted()
	return r
}

func (r *Relation[F, T]) Execute(db DB, result []F) error {
	if len(result) < 1 {
		return nil
	}
	var ids []string
	for _, r2 := range result {
		ids = append(ids, r.relationType.collectID(r2).String())
	}

	query := Query[T]().Where(r.relationType.queryColumn(), "IN", ids)
	if len(r.query.filters) > 0 {
		query = query.WhereF(func(q *QueryBuilder[T]) *QueryBuilder[T] {
			return r.query
		})
	}

	related, err := query.Get(db)
	if err != nil {
		return err
	}

	m := map[uuid.UUID][]T{}
	for _, r2 := range related {
		group := r.relationType.groupBy(r2)
		m[r.relationType.groupBy(r2)] = append(m[group], r2)
	}
	for i := range result {
		rd := RelationData[T]{}
		rd.data = m[r.relationType.collectID(result[i])]
		rd.loaded = true

		a := r.relationData(&result[i])
		*a = rd
	}
	return nil
}

type relationType[F, T Modeler] interface {
	collectID(F) uuid.UUID
	groupBy(T) uuid.UUID
	queryColumn() string
}
