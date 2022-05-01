package manager

import (
	"context"
	"database/sql"
	"fmt"
)

func (manager *Manager) GetAADLogin(context context.Context, username string) (*DBLogin, error) {
	var login DBLogin

	statement := fmt.Sprintf(`
    SELECT name, principal_id
    FROM master.sys.server_principals
    WHERE name = '%[1]s'
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

func (manager *Manager) CreateAADLogin(context context.Context, username string) error {
	statement := fmt.Sprintf(`
    CREATE LOGIN [%[1]s] FROM EXTERNAL PROVIDER
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
