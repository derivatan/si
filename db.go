package si

// DB is based on `sql.DB`, but generalized with an implementation independent version of `Rows`.
type DB interface {
	Query(query string, args ...any) (Rows, error)
	Exec(query string, args ...any) (any, error)
}

type Rows interface {
	Next() bool
	Scan(dest ...any) error
	Close() error
}
