package si

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/google/uuid"
)

type Q[T Modeler] struct {
	withs       []func(m T, r []T) error
	withDeleted bool

	selects    []string
	selectScan func(scan func(...any))
	joins      []func(t T) *JoinConf

	filters []filter
	orderBy []orderBy
	take    int
	skip    int

	havings  []filter
	groupBys []string

	//Used during building
	argsCounter int
	args        []any
}

///////////////
// Executors //
///////////////

// Get will Execute the query and return a list of the result.
func (q *Q[T]) Get(db DB) ([]T, error) {
	query := q.buildSelect()
	log(query, q.args)
	rows, err := db.Query(query, q.args...)
	if err != nil {
		return nil, fmt.Errorf("si.get: execute query: %w", err)
	}
	defer func() { _ = rows.Close() }()

	result := &[]T{}
	for rows.Next() {
		row := new(T)
		ti := getTypeInfo(row)
		var err error
		if len(q.selects) > 0 {
			q.selectScan(func(scan ...any) {
				err = rows.Scan(scan...)
			})
		} else {
			err = rows.Scan(ti.Values...)
		}
		if err != nil {
			return nil, fmt.Errorf("si.get: scan: %w", err)
		}
		reflect.ValueOf(result).Elem().Set(reflect.Append(reflect.ValueOf(result).Elem(), reflect.ValueOf(row).Elem()))
	}

	err = q.executeWith(*result)
	if err != nil {
		return nil, err
	}
	return *result, nil
}

// First will execute the query and return the first element of the result
func (q *Q[T]) First(db DB) (*T, error) {
	q.take = 1
	result, err := q.Get(db)
	if err != nil {
		return nil, fmt.Errorf("si.first: %w", err)
	}
	return &result[0], nil
}

// Find will return the one element in the query result.
// This will be successful IFF there was one result.
// The variadic parameter `id` is used to make it optional. If present, only the first element is used.
func (q *Q[T]) Find(db DB, id ...uuid.UUID) (*T, error) {
	if len(id) >= 1 {
		q = q.Where("id", "=", id[0])
	}
	result, err := q.Get(db)
	if err != nil {
		return nil, fmt.Errorf("si.find: %w", err)
	}

	if len(result) != 1 {
		return nil, ResourceNotFound()
	}
	return &result[0], nil
}

// MustGet is same as Get, but will panic on error.
func (q *Q[T]) MustGet(db DB) []T {
	result, err := q.Get(db)
	if err != nil {
		panic(err)
	}
	return result
}

// MustFirst is same as First, but will panic on error.
func (q *Q[T]) MustFirst(db DB) *T {
	result, err := q.First(db)
	if err != nil {
		panic(err)
	}
	return result
}

// MustFind is same as Find, but will panic on error.
func (q *Q[T]) MustFind(db DB, id ...uuid.UUID) *T {
	result, err := q.Find(db, id...)
	if err != nil {
		panic(err)
	}
	return result
}

func (q *Q[T]) executeWith(results []T) error {
	for _, with := range q.withs {
		var dummy T
		err := with(dummy, results)
		if err != nil {
			return err
		}
	}
	return nil
}

////////////////////
// Query Builders //
////////////////////

type filter struct {
	Column    string
	Operation string
	Value     any

	Separator string
	Sub       []filter
}

type Raw string

func (q *Q[T]) Select(selects []string, selectScan func(scan func(...any))) *Q[T] {
	if len(q.selects) > 0 {
		log("Select values are already set. Ignoring new values.")
		return q
	}
	q.selects = selects
	q.selectScan = selectScan
	return q
}

type JoinType string

const (
	INNER JoinType = "INNER"
	LEFT           = "LEFT"
	RIGHT          = "RIGHT"
	FULL           = "FULL"
)

type JoinConf struct {
	JoinType JoinType
	Table    string
	Alias    string
	// Extra condition?
	Condition []filter // func(q *Q[T]) *Q[T]
}

// Join adds a join on the query. Can be used with `join` a `Relation` to automate the condition.
func (q *Q[T]) Join(f func(t T) *JoinConf) *Q[T] {
	q.joins = append(q.joins, f)
	return q
}

// Where adds a condition, separated by `AND`
func (q *Q[T]) Where(column, op string, value any) *Q[T] {
	q.filters = append(q.filters, filter{Column: column, Operation: op, Value: value, Separator: "AND"})
	return q
}

// OrWhere adds a condition, separated by `OR`
func (q *Q[T]) OrWhere(column, op string, value any) *Q[T] {
	q.filters = append(q.filters, filter{Column: column, Operation: op, Value: value, Separator: "OR"})
	return q
}

// WhereF add a condition in parentheses, separated by `AND`
func (q *Q[T]) WhereF(f func(q *Q[T]) *Q[T]) *Q[T] {
	subQ := &Q[T]{}
	subQ = f(subQ)
	q.filters = append(q.filters, filter{Separator: "AND", Sub: subQ.filters})
	return q
}

// OrWhereF add a condition in parentheses, separated by `OR`
func (q *Q[T]) OrWhereF(f func(q *Q[T]) *Q[T]) *Q[T] {
	subQ := &Q[T]{}
	subQ = f(subQ)
	q.filters = append(q.filters, filter{Separator: "OR", Sub: subQ.filters})
	return q
}

type orderBy struct {
	Column    string
	Ascending bool
}

// OrderBy adds an order to the query.
func (q *Q[T]) OrderBy(column string, asc bool) *Q[T] {
	q.orderBy = append(q.orderBy, orderBy{column, asc})
	return q
}

// Take will limit the result to the given number.
func (q *Q[T]) Take(number int) *Q[T] {
	q.take = number
	return q
}

// Skip will remove the first `number`of the result.
func (q *Q[T]) Skip(number int) *Q[T] {
	q.skip = number
	return q
}

func (q *Q[T]) GroupBy(field string) *Q[T] {
	q.groupBys = append(q.groupBys, field)
	return q
}

func (q *Q[T]) Having(column, op string, value any) *Q[T] {
	q.havings = append(q.havings, filter{Column: column, Operation: op, Value: value, Separator: "AND"})
	return q
}

func (q *Q[T]) OrHaving(column, op string, value any) *Q[T] {
	q.havings = append(q.havings, filter{Column: column, Operation: op, Value: value, Separator: "OR"})
	return q
}

func (q *Q[T]) HavingF(f func(q *Q[T]) *Q[T]) *Q[T] {
	subQ := &Q[T]{}
	subQ = f(subQ)
	q.havings = append(q.havings, filter{Separator: "AND", Sub: subQ.filters})
	return q
}

func (q *Q[T]) OrHavingF(f func(q *Q[T]) *Q[T]) *Q[T] {
	subQ := &Q[T]{}
	subQ = f(subQ)
	q.havings = append(q.havings, filter{Separator: "OR", Sub: subQ.filters})
	return q
}

// With will retrieve a relation, while getting the main object(s).
func (q *Q[T]) With(f func(m T, r []T) error) *Q[T] {
	q.withs = append(q.withs, f)
	return q
}

// WithDeleted will ignore the deleted timestamp.
func (q *Q[T]) WithDeleted() *Q[T] {
	if !config.useDeletedAt {
		log("WithDeleted does nothing if the logs are disabled.")
	}
	q.withDeleted = true
	return q
}

func (q *Q[T]) buildSelect() string {
	specialSelect := len(q.selects) > 0
	t := new(T)
	table := (*t).GetTable()
	query := "SELECT "

	// Select
	if specialSelect {
		query += strings.Join(q.selects, ",")
	} else {
		var list []string
		for _, c := range getTypeInfo(t).Columns {
			list = append(list, fmt.Sprintf("%s.%s", table, c))
		}
		query += strings.Join(list, ",")
	}

	// From
	query += fmt.Sprintf(" FROM %s", table)

	// Joins
	for _, jf := range q.joins {
		j := jf(*t)
		if config.useDeletedAt && !q.withDeleted {
			j.Condition = append(j.Condition, filter{Column: j.Table + ".deleted_at", Operation: "IS", Value: nil, Separator: "AND"})
		}
		condition := q.buildFilters(j.Condition)
		query += fmt.Sprintf(" %s JOIN %s ON%s", j.JoinType, j.Table, condition)
	}

	// With Deleted
	if config.useDeletedAt && !q.withDeleted {
		otherFilters := q.filters
		q.filters = []filter{{Column: table + ".deleted_at", Operation: "IS", Value: nil}}
		if len(otherFilters) > 0 {
			q.filters = append(q.filters, filter{
				Separator: "AND",
				Sub:       otherFilters,
			})
		}
	}

	// Where
	if len(q.filters) > 0 {
		filterSql := q.buildFilters(q.filters)
		query += fmt.Sprintf(" WHERE%s", filterSql)
	}

	// Group By
	if len(q.groupBys) > 0 && specialSelect {
		query += " GROUP BY " + strings.Join(q.groupBys, ", ")
	}

	// Having
	if len(q.havings) > 0 && len(q.groupBys) > 0 && specialSelect {
		filterSql := q.buildFilters(q.havings)
		query += fmt.Sprintf(" HAVING%s", filterSql)
	}

	// Order by
	if len(q.orderBy) > 0 {
		query += " ORDER BY "
		for i, by := range q.orderBy {
			if i != 0 {
				query += ", "
			}
			query += fmt.Sprintf("%s ", by.Column)
			if by.Ascending {
				query += "asc "
			} else {
				query += "desc"
			}
		}
	}

	// Limit
	if q.take > 0 {
		query += fmt.Sprintf(" LIMIT %d ", q.take)
	}

	// Offset
	if q.skip > 0 {
		query += fmt.Sprintf(" OFFSET %d ", q.skip)
	}

	return query
}

func (q *Q[T]) buildFilters(filters []filter) string {
	var query string
	for i, f := range filters {
		if i != 0 {
			query += fmt.Sprintf(" %s", f.Separator)
		}

		// Handle nested parentheses
		if f.Sub != nil {
			subSql := q.buildFilters(f.Sub)
			query += fmt.Sprintf(" (%s)", subSql)
			continue
		}

		// Handle IS NULL.
		if f.Operation == "IS" && f.Value == nil {
			query += fmt.Sprintf(" %s IS NULL", f.Column)
			continue
		}

		// Handle Raw condition
		if _, ok := f.Value.(Raw); ok {
			query += fmt.Sprintf(" %s %s %s", f.Column, f.Operation, f.Value)
			continue
		}

		// Handle IN list
		parameters := []string{}
		if f.Operation == "IN" {
			for _, elem := range f.Value.([]string) {
				q.args = append(q.args, elem)
				q.argsCounter += 1
				parameters = append(parameters, fmt.Sprintf("$%d", q.argsCounter))
			}
			query += fmt.Sprintf(" %s IN (%s)", f.Column, strings.Join(parameters, ","))
			continue
		}

		q.args = append(q.args, f.Value)
		// Default condition handling.
		q.argsCounter += 1
		query += fmt.Sprintf(" %s %s $%d", f.Column, f.Operation, q.argsCounter)
	}
	return query
}
