package si

import "database/sql"

// This is an example to use SI with the standard `database/sql` library.

// If you want to use another library, you must create a
// similar wrapper, that implements the `si.DB` interface. (in `db.go`)

// SqlDB is implemented by both `sql.DB` and `sql.Tx`.
type SqlDB interface {
	Query(query string, args ...any) (*sql.Rows, error)
	Exec(query string, args ...any) (sql.Result, error)
}

func WrapDB(db SqlDB) DB {
	return &DBWrap{
		db: db,
	}
}

// DB

type DBWrap struct {
	db SqlDB
}

func (db *DBWrap) Query(query string, args ...any) (Rows, error) {
	result, err := db.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	return &RowsWrap{rows: result}, nil
}

func (db *DBWrap) Exec(query string, args ...any) (any, error) {
	_, err := db.db.Exec(query, args...)
	return nil, err
}

// Rows

type RowsWrap struct {
	rows *sql.Rows
}

func (w *RowsWrap) Next() bool {
	return w.rows.Next()
}

func (w *RowsWrap) Scan(a ...any) error {
	return w.rows.Scan(a...)
}

func (w *RowsWrap) Close() error {
	return w.rows.Close()
}
