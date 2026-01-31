package si

import (
	"fmt"
	"strings"
)

type S[T Modeler] struct {
	q    *Q[T]
	sets []SetConf
}

type SetConf struct {
	column string
	value  any
}

func (s *S[T]) Do(db DB) error {
	query := s.buildSet()
	log(query, s.q.args)
	_, err := db.Exec(query, s.q.args...)
	if err != nil {
		return fmt.Errorf("si.set: execute query: %w", err)
	}
	return nil
}

// WithDeleted will ignore the deleted timestamp.
func (s *S[T]) WithDeleted() *S[T] {
	s.q = s.q.WithDeleted()
	return s
}

func (s *S[T]) Join(f func(t T) *JoinConf) *S[T] {
	s.q = s.q.Join(f)
	return s
}

func (s *S[T]) Where(column, op string, value any) *S[T] {
	s.q = s.q.Where(column, op, value)
	return s
}

func (s *S[T]) OrWhere(column, op string, value any) *S[T] {
	s.q.OrWhere(column, op, value)
	return s
}

func (s *S[T]) WhereF(f func(s *Q[T]) *Q[T]) *S[T] {
	s.q = s.q.WhereF(f)
	return s
}

func (s *S[T]) OrWhereF(f func(q *Q[T]) *Q[T]) *S[T] {
	s.q = s.q.OrWhereF(f)
	return s
}

func (s *S[T]) Set(column string, value any) *S[T] {
	s.sets = append(s.sets, SetConf{column: column, value: value})
	return s
}

func (s *S[T]) buildSet() string {
	t := new(T)
	table := (*t).GetTable()

	// Update
	query := "UPDATE " + table

	// Join
	for _, jf := range s.q.joins {
		j := jf(*t)
		if config.useDeletedAt && !s.q.withDeleted {
			j.Condition = append(j.Condition, filter{Column: j.Table + ".deleted_at", Operation: "IS", Value: nil})
		}
		condition := s.q.buildFilters(j.Condition)
		query += fmt.Sprintf(" %s JOIN %s ON%s", j.JoinType, j.Table, condition)
	}

	// Set
	query += " SET "
	var list []string
	for _, set := range s.sets {
		if _, ok := set.value.(Raw); ok {
			list = append(list, fmt.Sprintf("%s = %s", set.column, set.value))
		} else {
			s.q.argsCounter += 1
			list = append(list, fmt.Sprintf("%s = $%d", set.column, s.q.argsCounter))
			s.q.args = append(s.q.args, set.value)
		}
	}
	query += strings.Join(list, ",")

	// Filter
	if len(s.q.filters) != 0 {
		filterSql := s.q.buildFilters(s.q.filters)
		query += fmt.Sprintf(" WHERE%s", filterSql)
	}

	return query
}
