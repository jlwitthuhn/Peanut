// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"errors"
	"net/http"
	"peanut/internal/data"
	"peanut/internal/security/passhash"
)

type UserService interface {
	CountUsers(req *http.Request) (int64, error)
	CreateUser(req *http.Request, name string, email string, plainPassword string) (string, error)
	IsEmailTaken(req *http.Request, email string) (bool, error)
	IsNameTaken(req *http.Request, username string) (bool, error)
}

func NewUserService(sessionDao data.SessionDao, userDao data.UserDao) UserService {
	return &userServiceImpl{sessionDao: sessionDao, userDao: userDao}
}

type userServiceImpl struct {
	sessionDao data.SessionDao
	userDao    data.UserDao
}

func (this *userServiceImpl) CountUsers(req *http.Request) (int64, error) {
	return this.userDao.CountRows(req)
}

func (this *userServiceImpl) CreateUser(req *http.Request, name string, email string, plainPassword string) (string, error) {
	nameTaken, nameErr := this.IsNameTaken(req, name)
	if nameErr != nil {
		return "", nameErr
	}
	if nameTaken {
		return "", errors.New("User name is already taken.")
	}
	emailTaken, emailErr := this.IsEmailTaken(req, email)
	if emailErr != nil {
		return "", emailErr
	}
	if emailTaken {
		return "", errors.New("User email is already taken.")
	}

	hashedPassword := passhash.GenerateDefaultPhcString(plainPassword)

	newId, insertErr := this.userDao.InsertRow(req, name, email, hashedPassword)
	if insertErr != nil {
		return "", insertErr
	}

	return newId, nil
}

func (this *userServiceImpl) IsEmailTaken(req *http.Request, email string) (bool, error) {
	count, err := this.userDao.CountRowsByEmail(req, email)
	if err != nil {
		return true, err
	}
	return count > 0, nil
}

func (this *userServiceImpl) IsNameTaken(req *http.Request, username string) (bool, error) {
	count, err := this.userDao.CountRowsByName(req, username)
	if err != nil {
		return true, err
	}
	return count > 0, nil
}
