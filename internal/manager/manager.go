package manager

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

type Manager struct {
	connectionString string
}

func GetManager(connectionString string) Manager {
	return Manager{
		connectionString: connectionString,
	}
}

func (manager *Manager) execute(context context.Context, command string, params ...interface{}) error {
	connection, err := manager.connect(context)
	if err != nil {
		return err
	}
	defer connection.Close()
	_, err = connection.ExecContext(context, command, params...)
	if err != nil {
		return err
	}
	return nil
}

func (manager *Manager) queryRow(context context.Context, query string, scanner func(*sql.Row) error, params ...interface{}) error {
	connection, err := manager.connect(context)
	if err != nil {
		return err
	}
	defer connection.Close()
	row := connection.QueryRowContext(context, query, params...)
	if row.Err() != nil {
		return row.Err()
	}
	return scanner(row)
}

func (manager *Manager) connect(context context.Context) (*sql.DB, error) {
	db, err := manager.open(context)
	if err == nil {
		return db, nil
	}

	// keep trying if connection fails
	interval := 1 * time.Second
	timeout := 1 * time.Minute

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	timeoutExceeded := time.After(timeout)
	for {
		select {
		case <-timeoutExceeded:
			return nil, fmt.Errorf("failed to connect to database")
		case <-ticker.C:
			db, err := manager.open(context)
			if err == nil {
				return db, nil
			}
		}
	}
}

func (manager *Manager) open(context context.Context) (*sql.DB, error) {
	db, err := sql.Open("sqlserver", manager.connectionString)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		defer db.Close()
		return nil, err
	}
	return db, nil
}
