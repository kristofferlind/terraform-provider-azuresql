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
    USE %[1]s
    SELECT IS_ROLEMEMBER ('%[3]s', '%[2]s')
  `, database, username, role)

	err := manager.queryRow(
		context,
		statement,
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
    USE %[1]s
    ALTER ROLE %[3]s ADD MEMBER %[2]s
  `, database, username, role)

	err := manager.execute(
		context,
		statement,
	)

	if err != nil {
		return err
	}

	return nil
}

func (manager *Manager) RemoveRole(context context.Context, username string, role string, database string) error {
	statement := fmt.Sprintf(`
    USE %[1]s
    ALTER ROLE %[3]s DROP MEMBER %[2]s
  `, database, username, role)

	err := manager.execute(
		context,
		statement,
	)

	if err != nil {
		return err
	}

	return nil
}
