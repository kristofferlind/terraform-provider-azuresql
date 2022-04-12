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
    USE %[1]s
    SELECT name FROM [sys].database_principals
    WHERE name = '%[2]s'
  `, database, username)

	err := manager.queryRow(
		context,
		statement,
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
    USE %[1]s
    CREATE USER %[2]s WITH PASSWORD = '%[3]s'
  `, database, username, password)

	err := manager.execute(
		context,
		statement,
	)

	if err != nil {
		return err
	}

	return nil
}

func (manager *Manager) CreateExternalUser(context context.Context, username string, database string) error {
	statement := fmt.Sprintf(`
    USE %[1]s
    CREATE USER %[2]s FROM EXTERNAL PROVIDER'
  `, database, username)

	err := manager.execute(
		context,
		statement,
	)

	if err != nil {
		return err
	}

	return nil
}

func (manager *Manager) CreateLoginUser(context context.Context, username string, database string) error {
	statement := fmt.Sprintf(`
    USE %[1]s
    CREATE USER %[2]s FROM LOGIN %[2]s
  `, database, username)

	err := manager.execute(
		context,
		statement,
	)

	if err != nil {
		return err
	}

	return nil
}

func (manager *Manager) DeleteLoginUser(context context.Context, username string, database string) error {
	statement := fmt.Sprintf(`
    USE %[1]s
    DROP USER %[2]s
  `, database, username)

	err := manager.execute(
		context,
		statement,
	)

	if err != nil {
		return err
	}

	return nil
}
