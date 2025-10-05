// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"database/sql"
	"errors"
	"peanut/internal/data"
	"peanut/internal/security/passhash"
)

type UserService interface {
	CountUsers(tx *sql.Tx) (int64, error)
	CreateUser(tx *sql.Tx, name string, email string, plainPassword string) (string, error)
	IsEmailTaken(tx *sql.Tx, email string) (bool, error)
	IsNameTaken(tx *sql.Tx, username string) (bool, error)
}

func NewUserService(sessionDao data.SessionDao, userDao data.UserDao) UserService {
	return &userServiceImpl{sessionDao: sessionDao, userDao: userDao}
}

type userServiceImpl struct {
	sessionDao data.SessionDao
	userDao    data.UserDao
}

func (this *userServiceImpl) CountUsers(tx *sql.Tx) (int64, error) {
	return this.userDao.CountRows(tx)
}

func (this *userServiceImpl) CreateUser(tx *sql.Tx, name string, email string, plainPassword string) (string, error) {
	nameTaken, nameErr := this.IsNameTaken(tx, name)
	if nameErr != nil {
		return "", nameErr
	}
	if nameTaken {
		return "", errors.New("User name is already taken.")
	}
	emailTaken, emailErr := this.IsEmailTaken(tx, email)
	if emailErr != nil {
		return "", emailErr
	}
	if emailTaken {
		return "", errors.New("User email is already taken.")
	}

	hashedPassword := passhash.GenerateDefaultPhcString(plainPassword)

	newId, insertErr := this.userDao.InsertRow(tx, name, email, hashedPassword)
	if insertErr != nil {
		return "", insertErr
	}

	return newId, nil
}

func (this *userServiceImpl) IsEmailTaken(tx *sql.Tx, email string) (bool, error) {
	count, err := this.userDao.CountRowsByEmail(tx, email)
	if err != nil {
		return true, err
	}
	return count > 0, nil
}

func (this *userServiceImpl) IsNameTaken(tx *sql.Tx, username string) (bool, error) {
	count, err := this.userDao.CountRowsByName(tx, username)
	if err != nil {
		return true, err
	}
	return count > 0, nil
}
