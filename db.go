package si

// DB interface is based on the `sql.DB` and `sql.Tx`
type DB interface {
	Query(query string, args ...any) (Rows, error)
	Exec(query string, args ...any) (any, error)
}

type Rows interface {
	Next() bool
	Scan(dest ...any) error
	Close() error
}
