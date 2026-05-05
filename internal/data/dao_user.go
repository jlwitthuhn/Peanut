// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package data

import (
	"context"
	"database/sql"
	"errors"
	"peanut/internal/logger"
	"time"

	"github.com/lib/pq"
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
	CreateDBObjects(ctx context.Context) error
	CountRows(ctx context.Context) (int64, error)
	CountRowsByEmail(ctx context.Context, name string) (int64, error)
	CountRowsByName(ctx context.Context, name string) (int64, error)
	InsertRow(ctx context.Context, name string, email string, hashedPassword string) (string, error)
	SelectRowById(ctx context.Context, id string) (*UserRow, error)
	SelectRowByName(ctx context.Context, name string) (*UserRow, error)
	SelectRowsAll(ctx context.Context) ([]UserRow, error)
	SelectRowsLikeName(ctx context.Context, namePattern string) ([]UserRow, error)
}

func NewUserDao() UserDao {
	return &userDaoImpl{}
}

type userDaoImpl struct{}

var sqlCreateTableUsers = `
	CREATE TABLE users (
		id UUID PRIMARY KEY DEFAULT uuidv7(),
		display_name VARCHAR(50) UNIQUE NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		password VARCHAR(100) NOT NULL,
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

func (*userDaoImpl) CreateDBObjects(ctx context.Context) error {
	sqlh := getSqlExecutorFromContext(ctx)
	_, err := sqlh.Exec(sqlCreateTableUsers)
	if err != nil {
		logger.Error(ctx, "Got database error on UserDao/CreateDBObjects query: ", err)
		return err
	}
	return nil
}

var sqlCountUsers = "SELECT COUNT(*) FROM users;"

func (*userDaoImpl) CountRows(ctx context.Context) (int64, error) {
	sqlh := getSqlExecutorFromContext(ctx)
	var count int64
	row := sqlh.QueryRow(sqlCountUsers)
	err := row.Scan(&count)
	if err != nil {
		logger.Error(ctx, "Got database error on UserDao/CountRows query: ", err)
		return 0, err
	}
	return count, nil
}

var sqlCountUsersByEmail = "SELECT COUNT(*) FROM users WHERE email = $1;"

func (*userDaoImpl) CountRowsByEmail(ctx context.Context, email string) (int64, error) {
	sqlh := getSqlExecutorFromContext(ctx)
	var count int64
	row := sqlh.QueryRow(sqlCountUsersByEmail, email)
	err := row.Scan(&count)
	if err != nil {
		logger.Error(ctx, "Got database error on UserDao/CountRowsByEmail query: ", err)
		return 0, err
	}
	return count, nil
}

var sqlCountUsersByName = "SELECT COUNT(*) FROM users WHERE display_name = $1;"

func (*userDaoImpl) CountRowsByName(ctx context.Context, name string) (int64, error) {
	sqlh := getSqlExecutorFromContext(ctx)
	var count int64
	row := sqlh.QueryRow(sqlCountUsersByName, name)
	err := row.Scan(&count)
	if err != nil {
		logger.Error(ctx, "Got database error on UserDao/CountRowsByName query: ", err)
		return 0, err
	}
	return count, nil
}

var sqlInsertUsersRow = "INSERT INTO users (display_name, email, password) VALUES ($1, $2, $3) RETURNING id"

func (*userDaoImpl) InsertRow(ctx context.Context, name string, email string, hashedPassword string) (string, error) {
	sqlh := getSqlExecutorFromContext(ctx)
	row := sqlh.QueryRow(sqlInsertUsersRow, name, email, hashedPassword)
	newId := ""
	err := row.Scan(&newId)
	if err != nil {
		logger.Error(ctx, "Got database error on UserDao/InsertRow query: ", err)
		return "", err
	}
	return newId, nil
}

var sqlSelectUsersRowById = "SELECT id, display_name, email, password, _created, _updated FROM users WHERE id = $1"

func (*userDaoImpl) SelectRowById(ctx context.Context, id string) (*UserRow, error) {
	sqlh := getSqlExecutorFromContext(ctx)
	result := &UserRow{}
	row := sqlh.QueryRow(sqlSelectUsersRowById, id)
	err := row.Scan(&result.Id, &result.DisplayName, &result.Email, &result.Password, &result.Created, &result.Updated)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		var pqErr *pq.Error
		ok := errors.As(err, &pqErr)
		if ok {
			if pqErr.Code == "22P02" {
				// INVALID_TEXT_REPRESENTATION
				// This means the input string is not a UUID, so no match exists
				return nil, nil
			}
		}
		logger.Error(ctx, "Got database error on UserDao/SelectRowById query: ", err)
		return nil, err
	}
	return result, nil
}

var sqlSelectUsersRowByName = "SELECT id, display_name, email, password, _created, _updated FROM users WHERE display_name = $1"

func (*userDaoImpl) SelectRowByName(ctx context.Context, name string) (*UserRow, error) {
	sqlh := getSqlExecutorFromContext(ctx)
	result := &UserRow{}
	row := sqlh.QueryRow(sqlSelectUsersRowByName, name)
	err := row.Scan(&result.Id, &result.DisplayName, &result.Email, &result.Password, &result.Created, &result.Updated)
	if err != nil {
		logger.Error(ctx, "Got database error on UserDao/SelectRowByName query: ", err)
		return nil, err
	}
	return result, nil
}

var sqlSelectUsersRowsAll = "SELECT id, display_name, email, password, _created, _updated FROM users ORDER BY _created"

func (*userDaoImpl) SelectRowsAll(ctx context.Context) ([]UserRow, error) {
	sqlh := getSqlExecutorFromContext(ctx)
	rows, err := sqlh.Query(sqlSelectUsersRowsAll)
	if err != nil {
		logger.Error(ctx, "Got database error on UserDao/SelectRowsAll query: ", err)
		return nil, err
	}
	defer rows.Close()

	var result []UserRow
	for rows.Next() {
		thisRow := UserRow{}
		err = rows.Scan(&thisRow.Id, &thisRow.DisplayName, &thisRow.Email, &thisRow.Password, &thisRow.Created, &thisRow.Updated)
		if err != nil {
			logger.Error(ctx, "Got database error on UserDao/SelectRowsAll query: ", err)
			return nil, err
		}
		result = append(result, thisRow)
	}
	return result, nil
}

var sqlSelectUsersRowsLikeName = `
	SELECT
		id, display_name, email, password, _created, _updated
	FROM
		users
	WHERE
	    display_name LIKE $1
	ORDER BY
	    _created
`

func (*userDaoImpl) SelectRowsLikeName(ctx context.Context, namePattern string) ([]UserRow, error) {
	sqlh := getSqlExecutorFromContext(ctx)
	rows, err := sqlh.Query(sqlSelectUsersRowsLikeName, namePattern)
	if err != nil {
		logger.Error(ctx, "Got database error on UserDao/SelectRowsLikeName query: ", err)
		return nil, err
	}
	defer rows.Close()
	var result []UserRow
	for rows.Next() {
		thisRow := UserRow{}
		err = rows.Scan(&thisRow.Id, &thisRow.DisplayName, &thisRow.Email, &thisRow.Password, &thisRow.Created, &thisRow.Updated)
		if err != nil {
			logger.Error(ctx, "Got database error on UserDao/SelectRowsLikeName query: ", err)
			return nil, err
		}
		result = append(result, thisRow)
	}
	return result, nil
}
