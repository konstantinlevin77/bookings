package driver

import "database/sql"

// DB is a struct that encompasses *sql.DB, in order to provide flexibility to change the database later.
type DB struct {
	SQL *sql.DB
}

// testDB tests database connection by pinging the connection given
func testDB(d *sql.DB) error {
	if err := d.Ping(); err != nil {
		return err
	}
	return nil
}

// NewDB creates a new DB struct using given data source name.
func NewDB(dsn string) (*DB, error) {

	d, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}

	if err = testDB(d); err != nil {
		return nil, err
	}

	db := &DB{SQL: d}
	return db, nil

}
