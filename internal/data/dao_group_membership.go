// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package data

import (
	"net/http"
	"peanut/internal/logger"
)

type GroupMembershipDao interface {
	CreateDBObjects(req *http.Request) error
	InsertRow(req *http.Request, userId string, groupId string) error
}

func NewGroupMembershipDao() GroupMembershipDao {
	return &groupMembershipDaoImpl{}
}

type groupMembershipDaoImpl struct{}

var sqlCreateTableGroupMembership = `
	CREATE TABLE group_membership (
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		group_id UUID NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
		_created TIMESTAMP WITH TIME ZONE NOT NULL,
		_updated TIMESTAMP WITH TIME ZONE NOT NULL,
		PRIMARY KEY (user_id, group_id)
	);

	CREATE TRIGGER
		group_membership_trigger_created_updated_before_insert
	BEFORE INSERT ON
		group_membership
	FOR EACH ROW EXECUTE FUNCTION
		fn_created_updated_before_insert();

	CREATE TRIGGER
		group_membership_trigger_created_updated_before_update
	BEFORE UPDATE ON
		group_membership
	FOR EACH ROW EXECUTE FUNCTION
		fn_created_updated_before_update();
`

func (*groupMembershipDaoImpl) CreateDBObjects(req *http.Request) error {
	sqlh := getSqlExecutorFromRequest(req)
	_, err := sqlh.Exec(sqlCreateTableGroupMembership)
	if err != nil {
		logger.Error(nil, "Got error on GroupMembershipDao/CreateDBObjects query: ", err)
		return err
	}
	return nil
}

var sqlInsertGroupMembershipRow = "INSERT INTO group_membership (user_id, group_id) VALUES ($1,$2)"

func (*groupMembershipDaoImpl) InsertRow(req *http.Request, userId string, groupId string) error {
	sqlh := getSqlExecutorFromRequest(req)
	_, err := sqlh.Exec(sqlInsertGroupMembershipRow, userId, groupId)
	if err != nil {
		logger.Error(nil, "Got error on GroupMembershipDao/InsertRow query: ", err)
		return err
	}
	return nil
}
