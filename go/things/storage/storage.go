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
	Query(ctx context.Context, namespace string, kind string, args ...any) (Rows, error)
	Insert(ctx context.Context, namespace string, kind string, args ...any) (int, error)
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
	ID           int
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

func (dbs *dbStorage) Query(ctx context.Context, namespace string, kind string, args ...any) (Rows, error) {
	queryArgs := []any{
		namespace,
		kind,
	}
	queryArgs = append(queryArgs, args...)

	fields := ""
	for i := range args {
		fields += fmt.Sprintf(", value%d", i+1)
	}

	query := "SELECT namespace, kind, tags, date_created, date_modified, id" + fields + " FROM things WHERE namespace = ? AND kind = ?"
	fmt.Println(query)
	rows, err := dbs.db.QueryContext(ctx, query, queryArgs...)
	if err != nil {
		return nil, err
	}

	return &dbRows{rows: rows}, nil
}

func (dbs *dbStorage) Insert(ctx context.Context, namespace string, kind string, args ...interface{}) (int, error) {
	return -1, fmt.Errorf("not implemented")
}
