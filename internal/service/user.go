// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"database/sql"
	"errors"
	"net/http"
	"peanut/internal/data"
	"peanut/internal/logger"
	"peanut/internal/security"
	"peanut/internal/security/passhash"
)

type UserService interface {
	CreateSession(r *http.Request, tx *sql.Tx, username string, plainPassword string) (string, error)
	CreateUser(tx *sql.Tx, name string, email string, plainPassword string) error
	GetLoggedInUserIdBySession(r *http.Request, tx *sql.Tx, sessionId string) (string, error)
	IsEmailTaken(tx *sql.Tx, email string) (bool, error)
	IsNameTaken(tx *sql.Tx, username string) (bool, error)
}

func NewUserService() UserService {
	return &userServiceImpl{}
}

type userServiceImpl struct{}

func (*userServiceImpl) CreateSession(r *http.Request, tx *sql.Tx, username string, plainPassword string) (string, error) {
	userDao := data.UserDaoInst()
	userRow, userErr := userDao.SelectRowByName(tx, username)
	if userErr != nil {
		return "", userErr
	}
	if passhash.ValidatePassword(plainPassword, userRow.Password) == false {
		return "", errors.New("Invalid password")
	}

	sessionDao := data.SessionDaoInst()
	newSessionId := security.GenerateSessionId()
	err := sessionDao.InsertRow(tx, newSessionId, userRow.Id)
	if err != nil {
		return "", err
	}

	logger.Info(r, "User logger in:", userRow.Id)

	return newSessionId, nil
}

func (this *userServiceImpl) CreateUser(tx *sql.Tx, name string, email string, plainPassword string) error {
	nameTaken, nameErr := this.IsNameTaken(tx, name)
	if nameErr != nil {
		return nameErr
	}
	if nameTaken {
		return errors.New("User name is already taken.")
	}
	emailTaken, emailErr := this.IsEmailTaken(tx, email)
	if emailErr != nil {
		return emailErr
	}
	if emailTaken {
		return errors.New("User email is already taken.")
	}

	hashedPassword := passhash.GenerateDefaultPhcString(plainPassword)

	userDao := data.UserDaoInst()
	insertErr := userDao.InsertRow(tx, name, email, hashedPassword)
	if insertErr != nil {
		return insertErr
	}

	return nil
}

func (*userServiceImpl) GetLoggedInUserIdBySession(r *http.Request, tx *sql.Tx, sessionId string) (string, error) {
	sessionDao := data.SessionDaoInst()
	sessionRow, sessionErr := sessionDao.SelectRowBySessionId(tx, sessionId)
	if sessionErr != nil {
		return "", sessionErr
	}
	return sessionRow.UserId, nil
}

func (*userServiceImpl) IsEmailTaken(tx *sql.Tx, email string) (bool, error) {
	userDao := data.UserDaoInst()
	count, err := userDao.CountRowsByEmail(tx, email)
	if err != nil {
		return true, err
	}
	return count > 0, nil
}

func (*userServiceImpl) IsNameTaken(tx *sql.Tx, username string) (bool, error) {
	userDao := data.UserDaoInst()
	count, err := userDao.CountRowsByName(tx, username)
	if err != nil {
		return true, err
	}
	return count > 0, nil
}
