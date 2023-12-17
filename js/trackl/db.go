package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/rubenv/sql-migrate"
)

var _ TasksStore = &dbStore{}

type dbStore struct {
	db *sql.DB
}

func newDBStore(driverName, dataSourceName string) (*dbStore, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	migrations := &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{
			{
				Id: "01-init",
				Up: []string{
					"CREATE TABLE tasks (namespace TEXT, id TEXT, icon TEXT, description TEXT, PRIMARY KEY (namespace, id))",
					"CREATE TABLE task_states (namespace TEXT, date TEXT, task_id TEXT, state TEXT, date_updated DATETIME, PRIMARY KEY (namespace, task_id, date), FOREIGN KEY (namespace, task_id) REFERENCES tasks (namespace, id))",
					"CREATE TABLE events (namespace TEXT, id TEXT, icon TEXT, date DATETIME, reference_date DATETIME, PRIMARY KEY (namespace, id))",
				},
				Down: []string{
					"DROP TABLE tasks",
					"DROP TABLE task_states",
					"DROP TABLE events",
				},
			},
		},
	}

	n, err := migrate.Exec(db, driverName, migrations, migrate.Up)
	if err != nil {
		return nil, fmt.Errorf("could not run migrations: %w", err)
	}

	if n != 0 {
		log.Printf("ran %d migrations", n)
	}

	return &dbStore{
		db: db,
	}, nil
}

func (ds *dbStore) Tasks() ([]Task, error) {
	rows, err := ds.db.Query("SELECT id, icon, description, COALESCE(s.state, 'not-done') FROM tasks LEFT JOIN task_states AS s ON id = s.task_id AND julianday(s.date) = julianday(date())", DefaultNamespace)
	if err != nil {
		return nil, fmt.Errorf("could not query tasks: %w", err)
	}
	defer rows.Close()

	tasks := make([]Task, 0, 10)
	taskIDs := make([]string, 0, 10)
	for rows.Next() {
		var task Task

		err := rows.Scan(&task.ID, &task.Icon, &task.Description, &task.State)
		if err != nil {
			return nil, fmt.Errorf("could not read task: %w", err)
		}

		tasks = append(tasks, task)
		taskIDs = append(taskIDs, task.ID)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("could not list tasks: %w", err)
	}

	return tasks, nil
}

func (ds *dbStore) FindTask(id string) (*Task, error) {
	row := ds.db.QueryRow("SELECT id, icon, description, COALESCE(s.state, 'not-done') FROM tasks LEFT JOIN task_states AS s ON tasks.namespace = s.namespace AND id = s.task_id WHERE tasks.namespace = ? AND id = ?", DefaultNamespace, id)

	var task Task
	err := row.Scan(&task.ID, &task.Icon, &task.Description, &task.State)
	if err != nil {
		return nil, fmt.Errorf("could not read task: %w", err)
	}

	return &task, nil
}

func (ds *dbStore) ChangeTaskState(id string, state TaskState) error {
	res, err := ds.db.Exec("INSERT OR REPLACE INTO task_states VALUES (?, ?, ?, ?, DATETIME())", DefaultNamespace, time.Now().Format(time.DateOnly), id, string(state))
	if err != nil {
		return fmt.Errorf("could not update task: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not check rows: %w", err)
	}

	if n != 1 {
		return fmt.Errorf("expected to change 1 row, but changed %d", n)
	}

	return nil
}

func (ds *dbStore) Events() ([]Event, error) {
	rows, err := ds.db.Query("SELECT id, icon, date, reference_date FROM events WHERE namespace = ?", DefaultNamespace)
	if err != nil {
		return nil, fmt.Errorf("could not list events: %w", err)
	}
	defer rows.Close()

	events := make([]Event, 0, 10)
	for rows.Next() {
		var event Event
		err := rows.Scan(&event.ID, &event.Icon, &event.Date, &event.ReferenceDate)
		if err != nil {
			return nil, fmt.Errorf("could not read event: %w", err)
		}

		events = append(events, event)
	}

	return events, nil
}

func (ds *dbStore) Close() error {
	return ds.db.Close()
}
