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
	CreateForum(ctx context.Context, sectionId string, name string, ordering float32, visibility string) (string, error)
	CreateSection(ctx context.Context, name string, ordering float32) (string, error)
	DeleteForum(ctx context.Context, id string) error
	DeleteSection(ctx context.Context, id string) error
	GetAllForumRows(ctx context.Context) ([]data.ForumRow, error)
	GetAllSectionRows(ctx context.Context) ([]data.ForumSectionRow, error)
	GetForumRowById(ctx context.Context, id string) (*data.ForumRow, error)
	GetSectionRowById(ctx context.Context, id string) (*data.ForumSectionRow, error)
	UpdateForumUserConfig(ctx context.Context, id string, sectionId string, name string, ordering float32, visibility string) error
	UpdateSectionUserConfig(ctx context.Context, id string, name string, ordering float32) error
}

func NewForumsService(forumsDao data.ForumsDao, forumSectionsDao data.ForumSectionsDao, systemLogDao data.SystemLogDao) ForumsService {
	return &forumsServiceImpl{forumsDao: forumsDao, forumSectionsDao: forumSectionsDao, systemLogDao: systemLogDao}
}

type forumsServiceImpl struct {
	forumsDao        data.ForumsDao
	forumSectionsDao data.ForumSectionsDao
	systemLogDao     data.SystemLogDao
}

func (this *forumsServiceImpl) CreateForum(ctx context.Context, sectionId string, name string, ordering float32, visibility string) (string, error) {
	if middleutil.ContextHasPermission(ctx, perms.Admin_Forums_Structure_Edit) == false {
		return "", errors.New("permission denied")
	}

	userId, ok := ctx.Value(contextkeys.UserId).(string)
	if !ok {
		return "", errors.New("cannot create forum: no user id in context")
	}

	newId, err := this.forumsDao.InsertRow(ctx, sectionId, name, ordering, visibility)
	if err != nil {
		return "", err
	}

	err = this.systemLogDao.InsertRow(ctx, userId, newId, fmt.Sprintf("Created forum: %s", name))
	if err != nil {
		return "", err
	}

	return newId, nil
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

func (this *forumsServiceImpl) DeleteForum(ctx context.Context, id string) error {
	if middleutil.ContextHasPermission(ctx, perms.Admin_Forums_Structure_Edit) == false {
		return errors.New("permission denied")
	}

	userId, ok := ctx.Value(contextkeys.UserId).(string)
	if !ok {
		return errors.New("cannot delete forum: no user id in context")
	}

	err := this.forumsDao.UpdateVisibilityById(ctx, id, "Deleted")
	if err != nil {
		return err
	}

	err = this.systemLogDao.InsertRow(ctx, userId, id, fmt.Sprintf("Deleted forum: %s", id))
	if err != nil {
		return err
	}

	return nil
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

func (this *forumsServiceImpl) GetAllForumRows(ctx context.Context) ([]data.ForumRow, error) {
	result, err := this.forumsDao.SelectRowAll(ctx)
	return result, err
}

func (this *forumsServiceImpl) GetForumRowById(ctx context.Context, id string) (*data.ForumRow, error) {
	result, err := this.forumsDao.SelectRowById(ctx, id)
	return result, err
}

func (this *forumsServiceImpl) GetAllSectionRows(ctx context.Context) ([]data.ForumSectionRow, error) {
	result, err := this.forumSectionsDao.SelectRowAll(ctx)
	return result, err
}

func (this *forumsServiceImpl) GetSectionRowById(ctx context.Context, id string) (*data.ForumSectionRow, error) {
	result, err := this.forumSectionsDao.SelectRowById(ctx, id)
	return result, err
}

func (this *forumsServiceImpl) UpdateSectionUserConfig(ctx context.Context, id string, name string, ordering float32) error {
	if middleutil.ContextHasPermission(ctx, perms.Admin_Forums_Structure_Edit) == false {
		return errors.New("permission denied")
	}

	userId, ok := ctx.Value(contextkeys.UserId).(string)
	if !ok {
		return errors.New("cannot update forum section: no user id in context")
	}

	err := this.forumSectionsDao.UpdateUserConfigById(ctx, id, name, ordering)
	if err != nil {
		return err
	}

	err = this.systemLogDao.InsertRow(ctx, userId, id, fmt.Sprintf("Updated forum section user config: %s", name))
	if err != nil {
		return err
	}

	return nil
}

func (this *forumsServiceImpl) UpdateForumUserConfig(ctx context.Context, id string, sectionId string, name string, ordering float32, visibility string) error {
	if middleutil.ContextHasPermission(ctx, perms.Admin_Forums_Structure_Edit) == false {
		return errors.New("permission denied")
	}

	userId, ok := ctx.Value(contextkeys.UserId).(string)
	if !ok {
		return errors.New("cannot update forum: no user id in context")
	}

	err := this.forumsDao.UpdateUserConfigById(ctx, id, sectionId, name, ordering, visibility)
	if err != nil {
		return err
	}

	err = this.systemLogDao.InsertRow(ctx, userId, id, fmt.Sprintf("Updated forum user config: %s", name))
	if err != nil {
		return err
	}

	return nil
}
