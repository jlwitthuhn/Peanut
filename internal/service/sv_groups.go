// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"context"
	"peanut/internal/data"
)

type GroupService interface {
	CreateGroup(ctx context.Context, name string, desc string, systemOwned bool) error
	GetAllGroupNames(ctx context.Context) ([]string, error)
	GetAllGroupRows(ctx context.Context) ([]data.GroupRow, error)
	GetGroupsByUserId(ctx context.Context, userId string) ([]string, error)
	GetUserRowsByGroupName(ctx context.Context, groupName string) ([]data.UserRow, error)
	EnrollUserInGroup(ctx context.Context, userId string, groupName string) error
}

func NewGroupService(groupDao data.GroupDao, groupMembershipDao data.GroupMembershipDao, multiTableDao data.MultiTableDao) GroupService {
	return &groupServiceImpl{groupDao: groupDao, groupMembershipDao: groupMembershipDao, multiTableDao: multiTableDao}
}

type groupServiceImpl struct {
	groupDao           data.GroupDao
	groupMembershipDao data.GroupMembershipDao
	multiTableDao      data.MultiTableDao
}

func (this *groupServiceImpl) CreateGroup(ctx context.Context, name string, desc string, systemOwned bool) error {
	err := this.groupDao.InsertRow(ctx, name, desc, systemOwned)
	return err
}

func (this *groupServiceImpl) GetAllGroupNames(ctx context.Context) ([]string, error) {
	groupRows, err := this.groupDao.SelectRowAll(ctx)
	if err != nil {
		return nil, err
	}
	var result []string
	for _, row := range groupRows {
		result = append(result, row.Name)
	}
	return result, nil
}

func (this *groupServiceImpl) GetAllGroupRows(ctx context.Context) ([]data.GroupRow, error) {
	return this.groupDao.SelectRowAll(ctx)
}

func (this *groupServiceImpl) GetGroupsByUserId(ctx context.Context, userId string) ([]string, error) {
	groupNames, err := this.multiTableDao.SelectGroupNamesByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	return groupNames, nil
}

func (this *groupServiceImpl) GetUserRowsByGroupName(ctx context.Context, groupName string) ([]data.UserRow, error) {
	userRows, err := this.multiTableDao.SelectUserRowsByGroupName(ctx, groupName)
	if err != nil {
		return nil, err
	}
	return userRows, nil
}

func (this *groupServiceImpl) EnrollUserInGroup(ctx context.Context, userId string, groupName string) error {
	groupRow, groupErr := this.groupDao.SelectRowByName(ctx, groupName)
	if groupErr != nil {
		return groupErr
	}
	err := this.groupMembershipDao.InsertRow(ctx, userId, groupRow.Id)
	if err != nil {
		return err
	}
	return nil
}
