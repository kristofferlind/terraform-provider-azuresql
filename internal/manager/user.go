package manager

import (
	"context"
	"database/sql"
	"fmt"
)

type DBLoginUser struct {
	LoginName string
	Database  string
}

func (manager *Manager) GetLoginUser(context context.Context, username string, database string) (*DBLoginUser, error) {
	var loginUser DBLoginUser

	statement := fmt.Sprintf(`
    SELECT name FROM [sys].database_principals
    WHERE name = '%[1]s'
  `, username)

	err := manager.queryRow(
		context,
		statement,
		database,
		func(row *sql.Row) error {
			return row.Scan(&loginUser.LoginName)
		},
	)

	if err != nil {
		return nil, err
	}

	loginUser.Database = database

	return &loginUser, nil
}

func (manager *Manager) CreateUser(context context.Context, username string, password string, database string) error {
	statement := fmt.Sprintf(`
    CREATE USER [%[1]s] WITH PASSWORD = '%[2]s'
  `, username, password)

	err := manager.execute(
		context,
		statement,
		database,
	)

	if err != nil {
		return err
	}

	return nil
}

func (manager *Manager) CreateExternalUser(context context.Context, username string, database string) error {
	statement := fmt.Sprintf(`
    CREATE USER [%[1]s] FROM EXTERNAL PROVIDER'
  `, username)

	err := manager.execute(
		context,
		statement,
		database,
	)

	if err != nil {
		return err
	}

	return nil
}

func (manager *Manager) CreateLoginUser(context context.Context, username string, database string) error {
	statement := fmt.Sprintf(`
    CREATE USER [%[1]s] FROM LOGIN [%[1]s]
  `, username)

	err := manager.execute(
		context,
		statement,
		database,
	)

	if err != nil {
		return err
	}

	return nil
}

func (manager *Manager) DeleteLoginUser(context context.Context, username string, database string) error {
	statement := fmt.Sprintf(`
    DROP USER [%[1]s]
  `, username)

	err := manager.execute(
		context,
		statement,
		database,
	)

	if err != nil {
		return err
	}

	return nil
}
