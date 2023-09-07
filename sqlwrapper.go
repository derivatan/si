package si

import "database/sql"

type SQLDB interface {
	Query(query string, args ...any) (*sql.Rows, error)
	Exec(query string, args ...any) (sql.Result, error)
}

type WrapDB struct {
	db *sql.DB
}

func NewSQLDB(db *sql.DB) *WrapDB {
	return &WrapDB{
		db: db,
	}
}

func (db *WrapDB) Query(query string, args ...any) (Rows, error) {
	result, err := db.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	return &WrapRows{rows: result}, nil
}

func (db *WrapDB) Exec(query string, args ...any) (any, error) {
	_, err := db.db.Exec(query, args...)
	return nil, err
}

type WrapRows struct {
	rows *sql.Rows
}

func (w *WrapRows) Next() bool {
	return w.rows.Next()
}

func (w *WrapRows) Scan(a ...any) error {
	return w.rows.Scan(a...)
}

func (w *WrapRows) Close() error {
	return w.rows.Close()
}
