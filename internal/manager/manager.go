package manager

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/microsoft/go-mssqldb"
	"github.com/microsoft/go-mssqldb/azuread"
)

type Manager struct {
	connectionString         string
	useAzureADAuthentication bool
}

func GetManager(connectionString string) Manager {
	useAzureADAuthentication := false
	if strings.Contains(connectionString, "fedauth") {
		useAzureADAuthentication = true
	}

	return Manager{
		connectionString:         connectionString,
		useAzureADAuthentication: useAzureADAuthentication,
	}
}

func (manager *Manager) execute(context context.Context, command string, database string) error {
	connection, err := manager.connect(context, database)
	if err != nil {
		return err
	}
	defer connection.Close()
	_, err = connection.ExecContext(context, command)
	if err != nil {
		return err
	}
	return nil
}

func (manager *Manager) queryRow(context context.Context, query string, database string, scanner func(*sql.Row) error) error {
	connection, err := manager.connect(context, database)
	if err != nil {
		return err
	}
	defer connection.Close()
	row := connection.QueryRowContext(context, query)
	if row.Err() != nil {
		return row.Err()
	}
	return scanner(row)
}

func (manager *Manager) connect(context context.Context, database string) (*sql.DB, error) {
	// keep trying to connect for a minute
	interval := 1 * time.Second
	timeout := 1 * time.Minute

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

  var collectedErrors []string

	timeoutExceeded := time.After(timeout)
	for {
		select {
		case <-timeoutExceeded:
      timeoutError := fmt.Errorf("Timeout reached, give up trying to connect")
      collectedErrors = append(collectedErrors, timeoutError.Error())
			return nil, fmt.Errorf(strings.Join(collectedErrors, "\n\n"))
		case <-ticker.C:
			db, err := manager.open(context, database)
			if err == nil {
				return db, nil
			} else {
        collectedErrors = append(collectedErrors, err.Error())

        // break early for errors that are unlikely to be fixed by trying again
        if strings.Contains(err.Error(), "Login failed for user") { // wrong credentials
          return nil, fmt.Errorf(strings.Join(collectedErrors, "\n\n"))
        }
      }
		}
	}
}

func (manager *Manager) open(context context.Context, database string) (*sql.DB, error) {
	driver := "sqlserver"
	connectionString := manager.connectionString
	if manager.useAzureADAuthentication {
		driver = azuread.DriverName
	}
	if len(database) > 0 {
		operator := "?"
		if strings.Contains(connectionString, "?") {
			operator = "&"
		}
		connectionString = connectionString + operator + "database=" + database
	}
	db, err := sql.Open(driver, connectionString)
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
