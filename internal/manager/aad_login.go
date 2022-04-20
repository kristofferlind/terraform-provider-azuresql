package manager

import (
	"context"
	"database/sql"
	"fmt"
)

func (manager *Manager) GetAADLogin(context context.Context, username string) (*DBLogin, error) {
	var login DBLogin

	statement := fmt.Sprintf(`
    SELECT name, principal_id, default_database_name
    FROM master.sys.server_principals
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

func (manager *Manager) CreateAADLogin(context context.Context, username string, defaultDatabase string) error {
	statement := fmt.Sprintf(`
    CREATE LOGIN [%[1]s] FROM EXTERNAL PROVIDER WITH DEFAULT_DATABASE = [%[2]s]
  `, username, defaultDatabase)

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

func (manager *Manager) UpdateAADLogin(context context.Context, username string, defaultDatabase string) error {
	statement := fmt.Sprintf(`
    ALTER LOGIN [%[1]s] WITH DEFAULT_DATABASE = [%[2]s]
  `, username, defaultDatabase)

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
