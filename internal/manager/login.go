package manager

import (
	"context"
	"database/sql"
	"fmt"
)

type DBLogin struct {
	LoginName   string
	PrincipalID string
}

func (manager *Manager) GetLogin(context context.Context, username string) (*DBLogin, error) {
	var login DBLogin

	statement := fmt.Sprintf(`
    SELECT name, principal_id FROM [master].[sys].[sql_logins] WHERE [name] = '%[1]s'
  `, username)

	err := manager.queryRow(
		context,
		statement,
		"",
		func(row *sql.Row) error {
			return row.Scan(&login.LoginName, &login.PrincipalID)
		},
	)
	if err != nil {
		return nil, err
	}

	return &login, nil
}

func (manager *Manager) CreateLogin(context context.Context, username string, password string) error {
	statement := fmt.Sprintf(`
    CREATE LOGIN %[1]s WITH PASSWORD = '%[2]s'
  `, username, password)

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

func (manager *Manager) UpdateLogin(context context.Context, username string, password string) error {
	statement := fmt.Sprintf(`
    ALTER LOGIN %[1]s WITH PASSWORD = '%[2]s'
  `, username, password)

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
    DROP LOGIN %[1]s
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
