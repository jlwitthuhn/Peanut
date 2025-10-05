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
	GetGroupsByUserId(r *http.Request, userId string) ([]string, error)
	EnrollUserInGroup(r *http.Request, tx *sql.Tx, userId string, groupName string) error
}

func NewGroupService(groupDao data.GroupDao, groupMembershipDao data.GroupMembershipDao, multiTableDao data.MultiTableDao) GroupService {
	return &groupServiceImpl{groupDao: groupDao, groupMembershipDao: groupMembershipDao, multiTableDao: multiTableDao}
}

type groupServiceImpl struct {
	groupDao           data.GroupDao
	groupMembershipDao data.GroupMembershipDao
	multiTableDao      data.MultiTableDao
}

func (this *groupServiceImpl) CreateGroup(tx *sql.Tx, name string, desc string, systemOwned bool) error {
	err := this.groupDao.InsertRow(tx, name, desc, systemOwned)
	return err
}

func (this *groupServiceImpl) GetGroupsByUserId(r *http.Request, userId string) ([]string, error) {
	groupNames, err := this.multiTableDao.SelectGroupNamesByUserId(r, userId)
	if err != nil {
		return nil, err
	}
	return groupNames, nil
}

func (this *groupServiceImpl) EnrollUserInGroup(r *http.Request, tx *sql.Tx, userId string, groupName string) error {
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
	groupRow, groupErr := this.groupDao.SelectRowByName(tx, groupName)
	if groupErr != nil {
		return groupErr
	}
	err := this.groupMembershipDao.InsertRow(tx, userId, groupRow.Id)
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
