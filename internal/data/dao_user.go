// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package data

import (
	"database/sql"
	"peanut/internal/data/datasource"
	"peanut/internal/logger"
	"sync"
)

type UserDao interface {
	CreateDBObjects(tx *sql.Tx) error
}

var userDaoInstance UserDao
var userDaoInstanceOnce sync.Once

func UserDaoInst() UserDao {
	userDaoInstanceOnce.Do(func() {
		userDaoInstance = &userDaoImpl{}
	})
	return userDaoInstance
}

type userDaoImpl struct{}

var sqlCreateTableUsers = `
	CREATE TABLE users (
		id BIGSERIAL PRIMARY KEY,
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

func (*userDaoImpl) CreateDBObjects(tx *sql.Tx) error {
	sqlh := selectExecutor(datasource.PostgresHandle(), tx)
	_, err := sqlh.Exec(sqlCreateTableUsers)
	if err != nil {
		logger.Error(nil, "Got error on UserDao/CreateDBObjects query: ", err)
		return err
	}
	return nil
}
