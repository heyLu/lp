package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Storage interface {
	Query(ctx context.Context, namespace string, conditions ...Condition) (Rows, error)
	Insert(ctx context.Context, row *Row) error
	Close() error
}

type Rows interface {
	Next() bool
	Scan(row *Row) error
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

func (dbr *dbRows) Close() error {
	return dbr.rows.Close()
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

func Kind(kind string) Condition       { return Condition{expr: "kind = ?", args: []any{kind}} }
func Summary(summary string) Condition { return Condition{expr: "summary = ?", args: []any{summary}} }

func Gt(field string, val any) Condition { return Condition{expr: field + " > ?", args: []any{val}} }
func Lt(field string, val any) Condition { return Condition{expr: field + " < ?", args: []any{val}} }
func Match(field string, val string) Condition {
	return Condition{expr: field + " LIKE concat('%', ?, '%')", args: []any{val}}
}

func (dbs *dbStorage) Query(ctx context.Context, namespace string, conditions ...Condition) (Rows, error) {
	query := "SELECT namespace, kind, id, summary, content, ref, number, float, bool, time, jsonb(fields_json), tags, date_created, date_modified FROM things_v2 WHERE namespace = ?"
	queryArgs := []any{namespace}

	for _, condition := range conditions {
		query += " AND " + condition.expr
		args := condition.args
		if len(args) > 0 {
			queryArgs = append(queryArgs, args...)
		}
	}

	rows, err := dbs.db.QueryContext(ctx, query+" ORDER BY date_created DESC", queryArgs...)
	if err != nil {
		return nil, err
	}

	return &dbRows{rows: rows}, nil
}

func (dbs *dbStorage) Insert(ctx context.Context, row *Row) error {
	if row.Namespace == "" {
		return fmt.Errorf("namespace cannot be empty")
	}
	if row.Kind == "" {
		return fmt.Errorf("kind cannot be empty")
	}
	if row.Summary == "" {
		return fmt.Errorf("summary cannot be empty")
	}

	row.ID = time.Now().Unix()
	row.DateModified = time.Now().UTC().Truncate(time.Second)

	tags := row.Tags
	tags = append(tags, tagsFromString(row.Summary)...)
	if row.Content.Valid {
		tags = append(tags, tagsFromString(row.Content.String)...)
	}
	slices.Sort(tags)
	tags = slices.Compact(tags)

	var fieldsJSON []byte
	if row.Fields != nil {
		var err error
		fieldsJSON, err = json.Marshal(row.Fields)
		if err != nil {
			return err
		}
	}

	res, err := dbs.db.ExecContext(ctx, `INSERT INTO things_v2 (namespace, kind, id, summary, content, ref, number, float, bool, time, fields_json, tags, date_created, date_modified) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		row.Namespace, row.Kind, row.ID, row.Summary,
		row.Content, row.Ref, row.Number, row.Float, row.Bool, row.Time, fieldsJSON,
		strings.Join(tags, ","), row.DateCreated.Unix(), row.DateModified.Unix(),
	)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n != 1 {
		return fmt.Errorf("expected %d changes, but %d changes happened", 1, n)
	}

	return nil
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

func (dbr *dbRows) Scan(row *Row) error {
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
