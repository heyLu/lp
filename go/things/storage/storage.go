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
	QueryV2(ctx context.Context, namespace string, conditions ...Condition) (Rows, error)
	Insert(ctx context.Context, namespace string, kind string, args ...any) (*Metadata, error)
	Close() error
}

type Rows interface {
	Next() bool
	Scan(args ...any) (*Metadata, error)
	ScanV2(row *Row) error
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

	_, err = db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS things_v2 (namespace TEXT NOT NULL, kind TEXT NOT NULL, id INTEGER NOT NULL, summary TEXT NOT NULL, content TEXT, ref TEXT, number INTEGER, float REAL, bool INTEGER, time INTEGER, fields_json BLOB, tags TEXT NOT NULL, date_created INTEGER NOT NULL, date_modified INTEGER NOT NULL, PRIMARY KEY (namespace, kind, id))")
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

type Option struct {
	Field string
	Op    string
	Value any
}

func (dbs *dbStorage) Query(ctx context.Context, namespace string, kind string, numFields int, args ...any) (Rows, error) {
	queryArgs := []any{
		namespace,
	}

	conditions := ""
	if kind != "" {
		conditions += " AND kind = ?"
		queryArgs = append(queryArgs, kind)
	}

	fields := ""
	for i := 0; i < numFields; i++ {
		fields += fmt.Sprintf(", value%d", i+1)
	}

	for i, arg := range args {
		if option, ok := arg.(Option); ok {
			conditions += fmt.Sprintf(" AND %s %s ?", option.Field, option.Op)
			queryArgs = append(queryArgs, option.Value)
			continue
		}

		conditions += fmt.Sprintf(" AND value%d = ?", i+1)
		queryArgs = append(queryArgs, arg)
	}

	query := "SELECT namespace, kind, tags, date_created, date_modified, id" + fields + " FROM things WHERE namespace = ?" + conditions + " ORDER BY date_created DESC"
	rows, err := dbs.db.QueryContext(ctx, query, queryArgs...)
	if err != nil {
		return nil, err
	}

	return &dbRows{rows: rows}, nil
}

func (dbs *dbStorage) Insert(ctx context.Context, namespace string, kind string, args ...any) (*Metadata, error) {
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

	fields := ""
	values := "?, ?, ?, ?, ?, ?"
	for i, arg := range args {
		fields += fmt.Sprintf(", value%d", i+1)
		values += ", ?"

		if s, ok := arg.(string); ok {
			metadata.Tags = tagsFromString(s)
		}
		if s, ok := arg.(*string); ok && s != nil {
			metadata.Tags = tagsFromString(*s)
		}
	}

	execArgs := []any{
		metadata.Namespace,
		metadata.Kind,
		strings.Join(metadata.Tags, ","),
		metadata.DateCreated.Unix(),
		metadata.DateModified.Unix(),
		metadata.ID,
	}
	execArgs = append(execArgs, args...)

	stmt := "INSERT INTO things (namespace, kind, tags, date_created, date_modified, id" + fields + ") VALUES (" + values + ")"
	res, err := dbs.db.ExecContext(ctx, stmt, execArgs...)
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

func tagsFromString(s string) []string {
	var tags []string
	parts := strings.Split(s, " ")
	for _, part := range parts {
		if len(part) > 0 && part[0] == '#' {
			if tags == nil {
				tags = make([]string, 0, 5)
			}
			tags = append(tags, part)
		}
	}
	return tags
}

func (dbs *dbStorage) Close() error {
	return dbs.db.Close()
}

// v2 sketch

type Condition struct {
	expr string
	args []any
}

func (c Condition) Expr() string { return c.expr }
func (c Condition) Args() []any  { return c.args }

func Kind(kind string) Condition { return Condition{expr: "kind = ?", args: []any{kind}} }

func (dbs *dbStorage) QueryV2(ctx context.Context, namespace string, conditions ...Condition) (Rows, error) {
	query := "SELECT namespace, kind, id, summary, content, ref, number, float, bool, time, jsonb(fields_json), tags, date_created, date_modified FROM things_v2 WHERE namespace = ?"
	queryArgs := []any{namespace}

	for _, condition := range conditions {
		query += " " + condition.expr
		args := condition.args
		if len(args) > 0 {
			queryArgs = append(queryArgs, args...)
		}
	}

	rows, err := dbs.db.QueryContext(ctx, query, queryArgs...)
	if err != nil {
		return nil, err
	}

	return &dbRows{rows: rows}, nil
}

type Row struct {
	Metadata

	Summary string
	Content sql.NullString
	Ref     sql.NullString
	Number  sql.NullInt64
	Float   sql.NullFloat64
	Bool    sql.NullBool
	Time    sql.NullTime

	Fields map[string]any
}

func (dbr *dbRows) ScanV2(row *Row) error {
	var fieldsRaw interface{}
	var tags string
	var dateCreated int64
	var dateModified int64
	err := dbr.rows.Scan(&row.Namespace, &row.Kind, &row.ID, &row.Summary, &row.Content, &row.Ref, &row.Number, &row.Float, &row.Bool, &row.Time, &fieldsRaw, &tags, &dateCreated, &dateModified)
	if err != nil {
		return err
	}

	row.DateCreated = time.Unix(dateCreated, 0).UTC()
	row.DateModified = time.Unix(dateModified, 0).UTC()

	row.Tags = strings.Split(tags, ",")
	if fieldsRaw != nil {
		fields, ok := fieldsRaw.(map[string]any)
		if !ok {
			return fmt.Errorf("invalid 'fields': %#v", fields)
		}
		row.Fields = fields
	}

	return nil
}
