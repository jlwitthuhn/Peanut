// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"context"
	"errors"
	"peanut/internal/data"
	"peanut/internal/security/passhash"
)

type UserService interface {
	CountUsers(ctx context.Context) (int64, error)
	CreateUser(ctx context.Context, name string, email string, plainPassword string) (string, error)
	GetUserRowById(ctx context.Context, id string) (*data.UserRow, error)
	GetUserRowsAll(ctx context.Context) ([]data.UserRow, error)
	GetUserRowsLikeName(ctx context.Context, namePattern string) ([]data.UserRow, error)
	IsEmailTaken(ctx context.Context, email string) (bool, error)
	IsNameTaken(ctx context.Context, username string) (bool, error)
}

func NewUserService(sessionDao data.SessionDao, userDao data.UserDao) UserService {
	return &userServiceImpl{sessionDao: sessionDao, userDao: userDao}
}

type userServiceImpl struct {
	sessionDao data.SessionDao
	userDao    data.UserDao
}

func (this *userServiceImpl) CountUsers(ctx context.Context) (int64, error) {
	return this.userDao.CountRows(ctx)
}

func (this *userServiceImpl) CreateUser(ctx context.Context, name string, email string, plainPassword string) (string, error) {
	nameTaken, nameErr := this.IsNameTaken(ctx, name)
	if nameErr != nil {
		return "", nameErr
	}
	if nameTaken {
		return "", errors.New("User name is already taken.")
	}
	emailTaken, emailErr := this.IsEmailTaken(ctx, email)
	if emailErr != nil {
		return "", emailErr
	}
	if emailTaken {
		return "", errors.New("User email is already taken.")
	}

	hashedPassword := passhash.GenerateDefaultPhcString(plainPassword)

	newId, insertErr := this.userDao.InsertRow(ctx, name, email, hashedPassword)
	if insertErr != nil {
		return "", insertErr
	}

	return newId, nil
}

func (this *userServiceImpl) GetUserRowById(ctx context.Context, id string) (*data.UserRow, error) {
	return this.userDao.SelectRowById(ctx, id)
}

func (this *userServiceImpl) GetUserRowsAll(ctx context.Context) ([]data.UserRow, error) {
	return this.userDao.SelectRowsAll(ctx)
}

func (this *userServiceImpl) GetUserRowsLikeName(ctx context.Context, namePattern string) ([]data.UserRow, error) {
	return this.userDao.SelectRowsLikeName(ctx, namePattern)
}

func (this *userServiceImpl) IsEmailTaken(ctx context.Context, email string) (bool, error) {
	count, err := this.userDao.CountRowsByEmail(ctx, email)
	if err != nil {
		return true, err
	}
	return count > 0, nil
}

func (this *userServiceImpl) IsNameTaken(ctx context.Context, username string) (bool, error) {
	count, err := this.userDao.CountRowsByName(ctx, username)
	if err != nil {
		return true, err
	}
	return count > 0, nil
}
