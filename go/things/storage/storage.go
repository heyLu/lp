package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Storage interface {
	Query(ctx context.Context, namespace string, kind string, numFields int, args ...any) (Rows, error)
	Insert(ctx context.Context, namespace string, kind string, args ...any) (*Metadata, error)
	Close() error
}

type Rows interface {
	Next() bool
	Scan(args ...any) (*Metadata, error)
	Close() error
}

type Metadata struct {
	Namespace    string
	Kind         string
	Tags         []string
	DateCreated  time.Time
	DateModified time.Time
	ID           int64
}

func NewDBStorage(ctx context.Context, dsn string) (Storage, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	_, err = db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS things (namespace TEXT, kind TEXT, tags TEXT, date_created INT, date_modified INT, id INT, value1 TEXT, value2 TEXT, value3 TEXT, value4 TEXT, value5 TEXT, value6 TEXT, value7 TEXT, value8 TEXT, value9 TEXT)")
	if err != nil {
		return nil, err
	}

	return &dbStorage{db: db}, nil
}

type dbStorage struct {
	db *sql.DB
}

type dbRows struct {
	rows *sql.Rows
}

func (dbr *dbRows) Next() bool {
	return dbr.rows.Next()
}

func (dbr *dbRows) Scan(args ...any) (*Metadata, error) {
	var metadata Metadata
	var tags string
	var dateCreated int64
	var dateModified int64
	scanArgs := []any{
		&metadata.Namespace,
		&metadata.Kind,
		&tags,
		&dateCreated,
		&dateModified,
		&metadata.ID,
	}
	scanArgs = append(scanArgs, args...)

	err := dbr.rows.Scan(scanArgs...)
	if err != nil {
		return nil, err
	}

	metadata.Tags = strings.Split(tags, ",")
	metadata.DateCreated = time.Unix(dateCreated, 0)
	metadata.DateModified = time.Unix(dateModified, 0)

	return &metadata, nil
}

func (dbr *dbRows) Close() error {
	return dbr.rows.Close()
}

func (dbs *dbStorage) Query(ctx context.Context, namespace string, kind string, numFields int, args ...any) (Rows, error) {
	queryArgs := []any{
		namespace,
		kind,
	}
	queryArgs = append(queryArgs, args...)

	fields := ""
	for i := 0; i < numFields; i++ {
		fields += fmt.Sprintf(", value%d", i+1)
	}

	conditions := ""
	for i := range args {
		conditions += fmt.Sprintf(" AND value%d = ?", i+1)
	}

	query := "SELECT namespace, kind, tags, date_created, date_modified, id" + fields + " FROM things WHERE namespace = ? AND kind = ?" + conditions
	rows, err := dbs.db.QueryContext(ctx, query, queryArgs...)
	if err != nil {
		return nil, err
	}

	return &dbRows{rows: rows}, nil
}

func (dbs *dbStorage) Insert(ctx context.Context, namespace string, kind string, args ...interface{}) (*Metadata, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("no values to insert")
	}

	metadata := Metadata{
		Namespace:    namespace,
		Kind:         kind,
		Tags:         nil,
		DateCreated:  time.Now().UTC().Truncate(time.Second),
		DateModified: time.Unix(0, 0),
		ID:           time.Now().UTC().Unix(),
	}
	queryArgs := []any{
		metadata.Namespace,
		metadata.Kind,
		strings.Join(metadata.Tags, ","),
		metadata.DateCreated.Unix(),
		metadata.DateModified.Unix(),
		metadata.ID,
	}
	queryArgs = append(queryArgs, args...)

	fields := ""
	values := "?, ?, ?, ?, ?, ?"
	for i := range args {
		fields += fmt.Sprintf(", value%d", i+1)
		values += ", ?"
	}

	query := "INSERT INTO things (namespace, kind, tags, date_created, date_modified, id" + fields + ") VALUES (" + values + ")"
	res, err := dbs.db.ExecContext(ctx, query, queryArgs...)
	if err != nil {
		return nil, err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if n != 1 {
		return nil, fmt.Errorf("expected %d changes, but %d changes happened", 1, n)
	}

	return &Metadata{}, nil
}

func (dbs *dbStorage) Close() error {
	return dbs.db.Close()
}
