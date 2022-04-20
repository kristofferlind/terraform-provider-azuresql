package manager

import (
	"context"
	"database/sql"
	"fmt"
)

type DBLogin struct {
	LoginName       string
	PrincipalID     string
	DefaultDatabase string
}

func (manager *Manager) GetLogin(context context.Context, username string) (*DBLogin, error) {
	var login DBLogin

	statement := fmt.Sprintf(`
    SELECT name, principal_id, default_database_name
    FROM master.sys.sql_logins
    WHERE name = '%[1]s'
  `, username)

	err := manager.queryRow(
		context,
		statement,
		"",
		func(row *sql.Row) error {
			return row.Scan(&login.LoginName, &login.PrincipalID, &login.DefaultDatabase)
		},
	)
	if err != nil {
		return nil, err
	}

	return &login, nil
}

func (manager *Manager) CreateLogin(context context.Context, username string, password string, defaultDatabase string) error {
	statement := fmt.Sprintf(`
    CREATE LOGIN [%[1]s] WITH PASSWORD = '%[2]s', DEFAULT_DATABASE = [%[3]s]
  `, username, password, defaultDatabase)

	err := manager.execute(
		context,
		statement,
		"",
	)

	if err != nil {
		return err
	}

	return nil
}

func (manager *Manager) UpdateLogin(context context.Context, username string, password string, defaultDatabase string) error {
	statement := fmt.Sprintf(`
    ALTER LOGIN [%[1]s] WITH PASSWORD = '%[2]s', DEFAULT_DATABASE = [%[3]s]
  `, username, password, defaultDatabase)

	err := manager.execute(
		context,
		statement,
		"",
	)
	if err != nil {
		return err
	}

	return nil
}

func (manager *Manager) DeleteLogin(context context.Context, username string) error {
	statement := fmt.Sprintf(`
    DROP LOGIN [%[1]s]
  `, username)

	err := manager.execute(
		context,
		statement,
		"",
	)
	if err != nil {
		return err
	}

	return nil
}
