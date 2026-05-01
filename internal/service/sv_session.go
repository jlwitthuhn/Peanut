// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"context"
	"errors"
	"peanut/internal/data"
	"peanut/internal/keynames/sessionkeys"
	"peanut/internal/logger"
	"peanut/internal/security"
	"peanut/internal/security/passhash"
)

type SessionService interface {
	CountUsersWithValidSession(ctx context.Context) (int64, error)
	CreateSession(ctx context.Context, username string, plainPassword string) (string, error)
	DestroySession(ctx context.Context, sessionId string) error
	GetLoggedInUserIdBySessionId(ctx context.Context, sessionId string) (string, error)
	GetString(ctx context.Context, sessionId string, key string) (string, error)
}

func NewSessionService(sessionDao data.SessionDao, sessionStringDao data.SessionStringDao, userDao data.UserDao) SessionService {
	return &sessionServiceImpl{sessionDao: sessionDao, sessionStringDao: sessionStringDao, userDao: userDao}
}

type sessionServiceImpl struct {
	sessionDao       data.SessionDao
	sessionStringDao data.SessionStringDao
	userDao          data.UserDao
}

func (this *sessionServiceImpl) CountUsersWithValidSession(ctx context.Context) (int64, error) {
	return this.sessionDao.CountValidDedupeByUser(ctx)
}

func (this *sessionServiceImpl) CreateSession(ctx context.Context, username string, plainPassword string) (string, error) {
	userRow, userErr := this.userDao.SelectRowByName(ctx, username)
	if userErr != nil {
		return "", userErr
	}
	if passhash.ValidatePassword(plainPassword, userRow.Password) == false {
		return "", errors.New("Invalid password")
	}

	newSessionId := security.GenerateSessionId()
	err := this.sessionDao.InsertRow(ctx, newSessionId, userRow.Id)
	if err != nil {
		return "", err
	}

	newCsrfToken := security.GenerateCsrfToken()
	err = this.sessionStringDao.UpsertString(ctx, newSessionId, sessionkeys.CsrfToken, newCsrfToken)
	if err != nil {
		return "", err
	}

	logger.Info(ctx, "User logger in:", userRow.Id)

	return newSessionId, nil
}

func (this *sessionServiceImpl) DestroySession(ctx context.Context, sessionId string) error {
	err := this.sessionDao.DeleteRowById(ctx, sessionId)
	if err != nil {
		return err
	}
	return nil
}

func (this *sessionServiceImpl) GetLoggedInUserIdBySessionId(ctx context.Context, sessionId string) (string, error) {
	sessionRow, sessionErr := this.sessionDao.SelectValidRowBySessionId(ctx, sessionId)
	if sessionErr != nil {
		return "", sessionErr
	}
	if sessionRow == nil {
		return "", nil
	}
	return sessionRow.UserId, nil
}

func (this *sessionServiceImpl) GetString(ctx context.Context, sessionId string, key string) (string, error) {
	row, err := this.sessionStringDao.SelectRow(ctx, sessionId, key)
	if err != nil {
		return "", err
	}
	return row.Value, nil
}
