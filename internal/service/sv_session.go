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

type SessionService interface {
	CountUsersWithValidSession(tx *sql.Tx) (int64, error)
	CreateSession(r *http.Request, tx *sql.Tx, username string, plainPassword string) (string, error)
	DestroySession(r *http.Request, tx *sql.Tx, sessionId string) error
	GetLoggedInUserIdBySessionId(r *http.Request, tx *sql.Tx, sessionId string) (string, error)
}

func NewSessionService(sessionDao data.SessionDao, userDao data.UserDao) SessionService {
	return &sessionServiceImpl{sessionDao: sessionDao, userDao: userDao}
}

type sessionServiceImpl struct {
	sessionDao data.SessionDao
	userDao    data.UserDao
}

func (this *sessionServiceImpl) CountUsersWithValidSession(tx *sql.Tx) (int64, error) {
	return this.sessionDao.CountValidDedupeByUser(tx)
}

func (this *sessionServiceImpl) CreateSession(r *http.Request, tx *sql.Tx, username string, plainPassword string) (string, error) {
	userRow, userErr := this.userDao.SelectRowByName(tx, username)
	if userErr != nil {
		return "", userErr
	}
	if passhash.ValidatePassword(plainPassword, userRow.Password) == false {
		return "", errors.New("Invalid password")
	}

	newSessionId := security.GenerateSessionId()
	err := this.sessionDao.InsertRow(tx, newSessionId, userRow.Id)
	if err != nil {
		return "", err
	}

	logger.Info(r, "User logger in:", userRow.Id)

	return newSessionId, nil
}

func (this *sessionServiceImpl) DestroySession(r *http.Request, tx *sql.Tx, sessionId string) error {
	err := this.sessionDao.DeleteRowById(tx, sessionId)
	if err != nil {
		return err
	}
	return nil
}

func (this *sessionServiceImpl) GetLoggedInUserIdBySessionId(r *http.Request, tx *sql.Tx, sessionId string) (string, error) {
	sessionRow, sessionErr := this.sessionDao.SelectValidRowBySessionId(tx, sessionId)
	if sessionErr != nil {
		return "", sessionErr
	}
	if sessionRow == nil {
		return "", nil
	}
	return sessionRow.UserId, nil
}
