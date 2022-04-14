package manager

import (
	"context"
	"database/sql"
	"fmt"
)

type DBUserRole struct {
	Name     string
	Database string
	Role     string
}

func (manager *Manager) GetUserWithRole(context context.Context, username string, role string, database string) (*DBUserRole, error) {
	var user DBUserRole
	var isRoleMember int

	statement := fmt.Sprintf(`
    SELECT IS_ROLEMEMBER ('%[2]s', '%[1]s')
  `, username, role)

	err := manager.queryRow(
		context,
		statement,
		database,
		func(row *sql.Row) error {
			return row.Scan(&isRoleMember)
		},
	)

	if err != nil {
		return nil, err
	}

	if isRoleMember == 1 {
		user.Role = role
	}

	user.Database = database
	user.Name = username

	return &user, nil
}

func (manager *Manager) AddRole(context context.Context, username string, role string, database string) error {
	statement := fmt.Sprintf(`
    ALTER ROLE %[2]s ADD MEMBER [%[1]s]
  `, username, role)

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

func (manager *Manager) RemoveRole(context context.Context, username string, role string, database string) error {
	statement := fmt.Sprintf(`
    ALTER ROLE %[2]s DROP MEMBER [%[1]s]
  `, username, role)

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
