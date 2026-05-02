// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"context"
	"errors"
	"fmt"
	"peanut/internal/data"
	"peanut/internal/keynames/contextkeys"
	"peanut/internal/middleutil"
	"peanut/internal/security/perms"
)

type ForumsService interface {
	CreateSection(ctx context.Context, name string, ordering float32) (string, error)
	DeleteSection(ctx context.Context, id string) error
	GetAllSectionRows(ctx context.Context) ([]data.ForumSectionRow, error)
	GetSectionRowById(ctx context.Context, id string) (*data.ForumSectionRow, error)
}

func NewForumsService(forumSectionsDao data.ForumSectionsDao, systemLogDao data.SystemLogDao) ForumsService {
	return &forumsServiceImpl{forumSectionsDao: forumSectionsDao, systemLogDao: systemLogDao}
}

type forumsServiceImpl struct {
	forumSectionsDao data.ForumSectionsDao
	systemLogDao     data.SystemLogDao
}

func (this *forumsServiceImpl) CreateSection(ctx context.Context, name string, ordering float32) (string, error) {
	if middleutil.ContextHasPermission(ctx, perms.Admin_Forums_Structure_Edit) == false {
		return "", errors.New("permission denied")
	}

	userId, ok := ctx.Value(contextkeys.UserId).(string)
	if !ok {
		return "", errors.New("cannot create forum section: no user id in context")
	}

	newId, err := this.forumSectionsDao.InsertRow(ctx, name, ordering)
	if err != nil {
		return "", err
	}

	err = this.systemLogDao.InsertRow(ctx, userId, newId, fmt.Sprintf("Created forum section: %s", name))
	if err != nil {
		return "", err
	}

	return newId, nil
}

func (this *forumsServiceImpl) DeleteSection(ctx context.Context, id string) error {
	if middleutil.ContextHasPermission(ctx, perms.Admin_Forums_Structure_Edit) == false {
		return errors.New("permission denied")
	}

	userId, ok := ctx.Value(contextkeys.UserId).(string)
	if !ok {
		return errors.New("cannot delete forum section: no user id in context")
	}

	err := this.forumSectionsDao.UpdateVisibilityById(ctx, id, "Deleted")
	if err != nil {
		return err
	}

	err = this.systemLogDao.InsertRow(ctx, userId, id, fmt.Sprintf("Deleted forum section: %s", id))
	if err != nil {
		return err
	}

	return nil
}

func (this *forumsServiceImpl) GetAllSectionRows(ctx context.Context) ([]data.ForumSectionRow, error) {
	result, err := this.forumSectionsDao.SelectRowAll(ctx)
	return result, err
}

func (this *forumsServiceImpl) GetSectionRowById(ctx context.Context, id string) (*data.ForumSectionRow, error) {
	result, err := this.forumSectionsDao.SelectRowById(ctx, id)
	return result, err
}
