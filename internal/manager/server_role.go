package manager

import (
	"context"
	"database/sql"
	"fmt"
)

type DBServerRole struct {
	Name string
	Role string
}

func (manager *Manager) GetServerRole(context context.Context, username string, role string) (*DBServerRole, error) {
	var user DBServerRole
	var isRoleMember int

	statement := fmt.Sprintf(`
    SELECT IS_SRVROLEMEMBER ('%[2]s', '%[1]s')
  `, username, role)

	err := manager.queryRow(
		context,
		statement,
		"",
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

	user.Name = username

	return &user, nil
}

func (manager *Manager) GrantServerRole(context context.Context, username string, role string) error {
	statement := fmt.Sprintf(`
    ALTER SERVER ROLE %[2]s ADD MEMBER [%[1]s]
  `, username, role)

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

func (manager *Manager) RevokeServerRole(context context.Context, username string, role string) error {
	statement := fmt.Sprintf(`
    ALTER SERVER ROLE %[2]s DROP MEMBER [%[1]s]
  `, username, role)

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
