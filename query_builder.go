package si

import (
	"fmt"
	"github.com/gofrs/uuid"
	"reflect"
	"strings"
)

type QueryBuilder[T Modeler] struct {
	db DB

	withs       []func(m T, r []T) error
	withDeleted bool

	filters []filter
	orderBy []orderBy
	take    int
	skip    int

	selects      []string
	selectValues []any
	groupValues  func() (any, []any)

	havings  []filter
	groupBys []string
}

///////////////
// Executors //
///////////////

// Get will Execute the query and return a list of the result.
func (q *QueryBuilder[T]) Get() ([]T, error) {
	query, values := q.buildSelect()
	log(query, values)
	rows, err := q.db.Query(query, values...)
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
			if q.groupValues != nil {
				obj, scans := q.groupValues()
				err = rows.Scan(scans...)
				// obj -> Add this onto a list???, how to return the list?
			} else {
				err = rows.Scan(q.selectValues...)
			}
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
func (q *QueryBuilder[T]) First() (*T, error) {
	q.take = 1
	result, err := q.Get()
	if err != nil {
		return nil, fmt.Errorf("si.first: %w", err)
	}
	return &result[0], nil
}

// Find will return the one element in the query result.
// This will be successful IFF there was one result.
// The variadic parameter `id` is used to make it optional. If present, only the first element is used.
func (q *QueryBuilder[T]) Find(id ...uuid.UUID) (*T, error) {
	if len(id) >= 1 {
		q = q.Where("id", "=", id[0])
	}
	result, err := q.Get()
	if err != nil {
		return nil, fmt.Errorf("si.find: %w", err)
	}

	if len(result) != 1 {
		return nil, ResourceNotFoundError
	}
	return &result[0], nil
}

// MustGet is same as Get, but will panic on error.
func (q *QueryBuilder[T]) MustGet() []T {
	result, err := q.Get()
	if err != nil {
		panic(err)
	}
	return result
}

// MustFirst is same as First, but will panic on error.
func (q *QueryBuilder[T]) MustFirst() *T {
	result, err := q.First()
	if err != nil {
		panic(err)
	}
	return result
}

// MustFind is same as Find, but will panic on error.
func (q *QueryBuilder[T]) MustFind(id ...uuid.UUID) *T {
	result, err := q.Find(id...)
	if err != nil {
		panic(err)
	}
	return result
}

func (q *QueryBuilder[T]) executeWith(results []T) error {
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

func (q *QueryBuilder[T]) Select(selects []string, values ...any) *QueryBuilder[T] {
	if len(q.selects) > 0 {
		log("Select values are already set. Ignoring new values.")
		return q
	}
	q.selects = selects
	q.selectValues = values
	return q
}

func (q *QueryBuilder[T]) GroupSelect(selects []string, scanArgs func() (any, []any)) *QueryBuilder[T] {
	if len(q.selects) > 0 {
		log("Select values are already set. Ignoring new values.")
		return q
	}
	q.selects = selects
	q.groupValues = scanArgs
	return q
}

// Where adds a condition, separated by `AND`
func (q *QueryBuilder[T]) Where(column, op string, value any) *QueryBuilder[T] {
	q.filters = append(q.filters, filter{Column: column, Operation: op, Value: value, Separator: "AND"})
	return q
}

// OrWhere adds a condition, separated by `OR`
func (q *QueryBuilder[T]) OrWhere(column, op string, value any) *QueryBuilder[T] {
	q.filters = append(q.filters, filter{Column: column, Operation: op, Value: value, Separator: "OR"})
	return q
}

// WhereF add a condition in parentheses, separated by `AND`
func (q *QueryBuilder[T]) WhereF(f func(q *QueryBuilder[T]) *QueryBuilder[T]) *QueryBuilder[T] {
	subQ := &QueryBuilder[T]{}
	subQ = f(subQ)
	q.filters = append(q.filters, filter{Separator: "AND", Sub: subQ.filters})
	return q
}

// OrWhereF add a condition in parentheses, separated by `OR`
func (q *QueryBuilder[T]) OrWhereF(f func(q *QueryBuilder[T]) *QueryBuilder[T]) *QueryBuilder[T] {
	subQ := &QueryBuilder[T]{}
	subQ = f(subQ)
	q.filters = append(q.filters, filter{Separator: "OR", Sub: subQ.filters})
	return q
}

type orderBy struct {
	Column    string
	Ascending bool
}

// OrderBy adds an order to the query.
func (q *QueryBuilder[T]) OrderBy(column string, asc bool) *QueryBuilder[T] {
	q.orderBy = append(q.orderBy, orderBy{column, asc})
	return q
}

// Take will limit the result to the given number.
func (q *QueryBuilder[T]) Take(number int) *QueryBuilder[T] {
	q.take = number
	return q
}

// Skip will remove the first `number`of the result.
func (q *QueryBuilder[T]) Skip(number int) *QueryBuilder[T] {
	q.skip = number
	return q
}

func (q *QueryBuilder[T]) GroupBy(field string) *QueryBuilder[T] {
	q.groupBys = append(q.groupBys, field)
	return q
}

func (q *QueryBuilder[T]) Having(condition any) *QueryBuilder[T] {
	q.havings = append(q.havings) // TODO.
	return q
}

// With will retrieve a relation, while getting the main object(s).
func (q *QueryBuilder[T]) With(f func(m T, r []T) error) *QueryBuilder[T] {
	q.withs = append(q.withs, f)
	return q
}

// WithDeleted will ignore the deleted timestamp.
func (q *QueryBuilder[T]) WithDeleted() *QueryBuilder[T] {
	q.withDeleted = true
	return q
}

func (q *QueryBuilder[T]) buildSelect() (string, []any) {
	var filterValues []any
	paramCounter := 1
	specialSelect := len(q.selects) > 0
	query := "SELECT "

	// Select
	if specialSelect {
		query += strings.Join(q.selects, ",")
	} else {
		typeInfo := getTypeInfo(new(T))
		query += strings.Join(typeInfo.Columns, ",")
	}

	// From
	query += fmt.Sprintf(" FROM %s", (*new(T)).GetTable())

	// With Deleted
	if !q.withDeleted {
		otherFilters := q.filters
		q.filters = []filter{{Column: "deleted_at", Operation: "IS", Value: nil}}
		if len(otherFilters) > 0 {
			q.filters = append(q.filters, filter{
				Separator: "AND",
				Sub:       otherFilters,
			})
		}
	}

	// Where
	if len(q.filters) > 0 {
		var filterSql string
		filterSql, filterValues, paramCounter = q.buildFilters(q.filters, paramCounter)
		query += fmt.Sprintf(" WHERE%s", filterSql)
	}

	// Group By
	if len(q.groupBys) > 0 && specialSelect {
		query += " GROUP BY " + strings.Join(q.groupBys, ", ")
	}

	// Having
	if len(q.havings) > 0 && len(q.groupBys) > 0 && specialSelect {
		filterSql, havingFilterValues, havingParamCounter := q.buildFilters(q.havings, paramCounter)
		query += fmt.Sprintf(" HAVING%s", filterSql)
		filterValues = append(filterValues, havingFilterValues)
		paramCounter += havingParamCounter
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

	return query, filterValues
}

func (q *QueryBuilder[T]) buildFilters(filters []filter, paramCounter int) (string, []any, int) {
	var query string
	var filterValues []any
	for i, filter := range filters {
		if i != 0 {
			query += fmt.Sprintf(" %s", filter.Separator)
		}
		if filter.Sub != nil {
			subSql, subFilterValues, pC := q.buildFilters(filter.Sub, paramCounter)
			paramCounter = pC
			filterValues = append(filterValues, subFilterValues...)
			query += fmt.Sprintf(" (%s)", subSql)
			continue
		}

		if filter.Operation == "IS" && filter.Value == nil {
			query += fmt.Sprintf(" %s IS NULL", filter.Column)
			continue
		}
		parameters := []string{}
		if filter.Operation == "IN" {
			for _, elem := range filter.Value.([]string) {
				filterValues = append(filterValues, elem)
				parameters = append(parameters, fmt.Sprintf("$%d", paramCounter))
				paramCounter += 1
			}
			query += fmt.Sprintf(" %s IN (%s)", filter.Column, strings.Join(parameters, ","))
			continue
		} else {
			filterValues = append(filterValues, filter.Value)
		}
		query += fmt.Sprintf(" %s %s $%d", filter.Column, filter.Operation, paramCounter)
		paramCounter += 1
	}
	return query, filterValues, paramCounter
}
