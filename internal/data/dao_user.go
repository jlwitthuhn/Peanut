// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package data

import (
	"net/http"
	"peanut/internal/logger"
	"time"
)

type UserRow struct {
	Id          string
	DisplayName string
	Email       string
	Password    string
	Created     time.Time
	Updated     time.Time
}

type UserDao interface {
	CreateDBObjects(req *http.Request) error
	CountRows(req *http.Request) (int64, error)
	CountRowsByEmail(req *http.Request, name string) (int64, error)
	CountRowsByName(req *http.Request, name string) (int64, error)
	InsertRow(req *http.Request, name string, email string, hashedPassword string) (string, error)
	SelectRowsAll(req *http.Request) ([]UserRow, error)
	SelectRowByName(req *http.Request, name string) (*UserRow, error)
}

func NewUserDao() UserDao {
	return &userDaoImpl{}
}

type userDaoImpl struct{}

var sqlCreateTableUsers = `
	CREATE TABLE users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		display_name VARCHAR(255) UNIQUE NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		_created TIMESTAMP WITH TIME ZONE NOT NULL,
		_updated TIMESTAMP WITH TIME ZONE NOT NULL
	);

	CREATE TRIGGER
		users_trigger_created_updated_before_insert
	BEFORE INSERT ON
		users
	FOR EACH ROW EXECUTE FUNCTION
		fn_created_updated_before_insert();

	CREATE TRIGGER
		users_trigger_created_updated_before_update
	BEFORE UPDATE ON
		users
	FOR EACH ROW EXECUTE FUNCTION
		fn_created_updated_before_update();
`

func (*userDaoImpl) CreateDBObjects(req *http.Request) error {
	sqlh := getSqlExecutorFromRequest(req)
	_, err := sqlh.Exec(sqlCreateTableUsers)
	if err != nil {
		logger.Error(nil, "Got error on UserDao/CreateDBObjects query: ", err)
		return err
	}
	return nil
}

var sqlCountUsers = "SELECT COUNT(*) FROM users;"

func (*userDaoImpl) CountRows(req *http.Request) (int64, error) {
	sqlh := getSqlExecutorFromRequest(req)
	var count int64
	row := sqlh.QueryRow(sqlCountUsers)
	err := row.Scan(&count)
	if err != nil {
		logger.Error(nil, "Got error on UserDao/CountRows query: ", err)
		return 0, err
	}
	return count, nil
}

var sqlCountUsersByEmail = "SELECT COUNT(*) FROM users WHERE email = $1;"

func (*userDaoImpl) CountRowsByEmail(req *http.Request, email string) (int64, error) {
	sqlh := getSqlExecutorFromRequest(req)
	var count int64
	row := sqlh.QueryRow(sqlCountUsersByEmail, email)
	err := row.Scan(&count)
	if err != nil {
		logger.Error(nil, "Got error on UserDao/CountRowsByEmail query: ", err)
		return 0, err
	}
	return count, nil
}

var sqlCountUsersByName = "SELECT COUNT(*) FROM users WHERE display_name = $1;"

func (*userDaoImpl) CountRowsByName(req *http.Request, name string) (int64, error) {
	sqlh := getSqlExecutorFromRequest(req)
	var count int64
	row := sqlh.QueryRow(sqlCountUsersByName, name)
	err := row.Scan(&count)
	if err != nil {
		logger.Error(nil, "Got error on UserDao/CountRowsByName query: ", err)
		return 0, err
	}
	return count, nil
}

var sqlInsertUsersRow = "INSERT INTO users (display_name, email, password) VALUES ($1, $2, $3) RETURNING id"

func (*userDaoImpl) InsertRow(req *http.Request, name string, email string, hashedPassword string) (string, error) {
	sqlh := getSqlExecutorFromRequest(req)
	row := sqlh.QueryRow(sqlInsertUsersRow, name, email, hashedPassword)
	newId := ""
	err := row.Scan(&newId)
	if err != nil {
		logger.Error(nil, "Got error on UserDao/InsertRow query: ", err)
		return "", err
	}
	return newId, nil
}

var sqlSelectUsersRowAll = "SELECT id, display_name, email, password, _created, _updated FROM users ORDER BY _created"

func (*userDaoImpl) SelectRowsAll(req *http.Request) ([]UserRow, error) {
	sqlh := getSqlExecutorFromRequest(req)
	rows, err := sqlh.Query(sqlSelectUsersRowAll)
	if err != nil {
		logger.Error(nil, "Got error on UserDao/SelectRowAll query: ", err)
		return nil, err
	}
	var result []UserRow
	for rows.Next() {
		thisRow := UserRow{}
		err = rows.Scan(&thisRow.Id, &thisRow.DisplayName, &thisRow.Email, &thisRow.Password, &thisRow.Created, &thisRow.Updated)
		if err != nil {
			return nil, err
		}
		result = append(result, thisRow)
	}
	return result, nil
}

var sqlSelectUsersRowByName = "SELECT id, display_name, email, password, _created, _updated FROM users WHERE display_name = $1"

func (*userDaoImpl) SelectRowByName(req *http.Request, name string) (*UserRow, error) {
	sqlh := getSqlExecutorFromRequest(req)
	result := &UserRow{}
	row := sqlh.QueryRow(sqlSelectUsersRowByName, name)
	err := row.Scan(&result.Id, &result.DisplayName, &result.Email, &result.Password, &result.Created, &result.Updated)
	if err != nil {
		logger.Error(nil, "Got error on UserDao/SelectRowByName query: ", err)
		return nil, err
	}
	return result, nil
}
