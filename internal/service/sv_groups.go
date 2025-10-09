// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"net/http"
	"peanut/internal/data"
)

type GroupService interface {
	CreateGroup(req *http.Request, name string, desc string, systemOwned bool) error
	GetAllGroupNames(req *http.Request) ([]string, error)
	GetGroupsByUserId(req *http.Request, userId string) ([]string, error)
	EnrollUserInGroup(req *http.Request, userId string, groupName string) error
}

func NewGroupService(groupDao data.GroupDao, groupMembershipDao data.GroupMembershipDao, multiTableDao data.MultiTableDao) GroupService {
	return &groupServiceImpl{groupDao: groupDao, groupMembershipDao: groupMembershipDao, multiTableDao: multiTableDao}
}

type groupServiceImpl struct {
	groupDao           data.GroupDao
	groupMembershipDao data.GroupMembershipDao
	multiTableDao      data.MultiTableDao
}

func (this *groupServiceImpl) CreateGroup(req *http.Request, name string, desc string, systemOwned bool) error {
	err := this.groupDao.InsertRow(req, name, desc, systemOwned)
	return err
}

func (this *groupServiceImpl) GetAllGroupNames(req *http.Request) ([]string, error) {
	groupRows, err := this.groupDao.SelectRowAll(req)
	if err != nil {
		return nil, err
	}
	var result []string
	for _, row := range groupRows {
		result = append(result, row.Name)
	}
	return result, nil
}

func (this *groupServiceImpl) GetGroupsByUserId(req *http.Request, userId string) ([]string, error) {
	groupNames, err := this.multiTableDao.SelectGroupNamesByUserId(req, userId)
	if err != nil {
		return nil, err
	}
	return groupNames, nil
}

func (this *groupServiceImpl) EnrollUserInGroup(req *http.Request, userId string, groupName string) error {
	groupRow, groupErr := this.groupDao.SelectRowByName(req, groupName)
	if groupErr != nil {
		return groupErr
	}
	err := this.groupMembershipDao.InsertRow(req, userId, groupRow.Id)
	if err != nil {
		return err
	}
	return nil
}
