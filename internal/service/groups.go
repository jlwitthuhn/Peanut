// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"database/sql"
	"net/http"
	"peanut/internal/data"
	"peanut/internal/data/datasource"
	"peanut/internal/logger"
)

type GroupService interface {
	CreateGroup(tx *sql.Tx, name string, desc string, systemOwned bool) error
	EnrollUserInGroup(r *http.Request, tx *sql.Tx, userId string, groupName string) error
}

func NewGroupService() GroupService {
	return &groupServiceImpl{}
}

type groupServiceImpl struct{}

func (*groupServiceImpl) CreateGroup(tx *sql.Tx, name string, desc string, systemOwned bool) error {
	err := data.GroupDaoInst().InsertRow(tx, name, desc, systemOwned)
	return err
}

func (*groupServiceImpl) EnrollUserInGroup(r *http.Request, tx *sql.Tx, userId string, groupName string) error {
	shouldCommit := false
	if tx == nil {
		newTx, txErr := datasource.PostgresHandle().BeginTx(r.Context(), nil)
		if txErr != nil {
			return txErr
		}
		tx = newTx
		shouldCommit = true
		defer tx.Rollback()
	}
	groupDao := data.GroupDaoInst()
	groupRow, groupErr := groupDao.SelectRowByName(tx, groupName)
	if groupErr != nil {
		return groupErr
	}
	groupMembershipDao := data.GroupMembershipDaoInst()
	err := groupMembershipDao.InsertRow(tx, userId, groupRow.Id)
	if err != nil {
		return err
	}
	if shouldCommit {
		commitErr := tx.Commit()
		if commitErr != nil {
			logger.Error(r, "Committing transaction in GroupService/EnrollUserInGroup failed with error: %v", commitErr)
		}
	}
	return nil
}
