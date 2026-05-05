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

type ForumService interface {
	CreateForum(ctx context.Context, sectionId string, name string, description string, ordering float32, visibility string) (string, error)
	CreateSection(ctx context.Context, name string, ordering float32) (string, error)
	DeleteForum(ctx context.Context, id string) error
	DeleteSection(ctx context.Context, id string) error
	GetAllForumRows(ctx context.Context) ([]data.ForumRow, error)
	GetAllSectionRows(ctx context.Context) ([]data.ForumSectionRow, error)
	GetForumRowById(ctx context.Context, id string) (*data.ForumRow, error)
	GetHomeViewRowsBySectionPublic(ctx context.Context, sectionId string) ([]data.ForumHomeViewRow, error)
	GetHomeViewRowsPublic(ctx context.Context) ([]data.ForumHomeViewRow, error)
	GetSectionRowById(ctx context.Context, id string) (*data.ForumSectionRow, error)
	IsForumReadable(ctx context.Context, id string) (bool, error)
	IsSectionReadable(ctx context.Context, id string) (bool, error)
	UpdateForumUserConfig(ctx context.Context, id string, sectionId string, name string, description string, ordering float32, visibility string) error
	UpdateSectionUserConfig(ctx context.Context, id string, name string, ordering float32) error
}

func NewForumService(forumsDao data.ForumsDao, forumSectionsDao data.ForumSectionsDao, forumHomeDvao data.ForumHomeDvao, systemLogDao data.SystemLogDao) ForumService {
	return &forumServiceImpl{forumsDao: forumsDao, forumSectionsDao: forumSectionsDao, forumHomeDvao: forumHomeDvao, systemLogDao: systemLogDao}
}

type forumServiceImpl struct {
	forumsDao        data.ForumsDao
	forumSectionsDao data.ForumSectionsDao
	forumHomeDvao    data.ForumHomeDvao
	systemLogDao     data.SystemLogDao
}

func (this *forumServiceImpl) CreateForum(ctx context.Context, sectionId string, name string, description string, ordering float32, visibility string) (string, error) {
	if middleutil.ContextHasPermission(ctx, perms.Admin_Forums_Structure_Edit) == false {
		return "", errors.New("permission denied")
	}

	userId, ok := ctx.Value(contextkeys.UserId).(string)
	if !ok {
		return "", errors.New("cannot create forum: no user id in context")
	}

	newId, err := this.forumsDao.InsertRow(ctx, sectionId, name, description, ordering, visibility)
	if err != nil {
		return "", err
	}

	err = this.systemLogDao.InsertRow(ctx, userId, newId, fmt.Sprintf("Created forum: %s", name))
	if err != nil {
		return "", err
	}

	return newId, nil
}

func (this *forumServiceImpl) CreateSection(ctx context.Context, name string, ordering float32) (string, error) {
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

func (this *forumServiceImpl) DeleteForum(ctx context.Context, id string) error {
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

func (this *forumServiceImpl) DeleteSection(ctx context.Context, id string) error {
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

func (this *forumServiceImpl) GetAllForumRows(ctx context.Context) ([]data.ForumRow, error) {
	result, err := this.forumsDao.SelectRowAll(ctx)
	return result, err
}

func (this *forumServiceImpl) GetHomeViewRowsBySectionPublic(ctx context.Context, sectionId string) ([]data.ForumHomeViewRow, error) {
	result, err := this.forumHomeDvao.SelectForumHomeViewRowsBySectionPublic(ctx, sectionId)
	return result, err
}

func (this *forumServiceImpl) GetHomeViewRowsPublic(ctx context.Context) ([]data.ForumHomeViewRow, error) {
	result, err := this.forumHomeDvao.SelectForumHomeViewRowsPublic(ctx)
	return result, err
}

func (this *forumServiceImpl) GetForumRowById(ctx context.Context, id string) (*data.ForumRow, error) {
	result, err := this.forumsDao.SelectRowById(ctx, id)
	return result, err
}

func (this *forumServiceImpl) GetAllSectionRows(ctx context.Context) ([]data.ForumSectionRow, error) {
	result, err := this.forumSectionsDao.SelectRowAll(ctx)
	return result, err
}

func (this *forumServiceImpl) GetSectionRowById(ctx context.Context, id string) (*data.ForumSectionRow, error) {
	result, err := this.forumSectionsDao.SelectRowById(ctx, id)
	return result, err
}

func (this *forumServiceImpl) IsForumReadable(ctx context.Context, id string) (bool, error) {
	forum, err := this.forumsDao.SelectRowById(ctx, id)
	if err != nil {
		return false, err
	}
	if forum == nil {
		return false, nil
	}
	return forum.Visibility == "Public", nil
}

func (this *forumServiceImpl) IsSectionReadable(ctx context.Context, id string) (bool, error) {
	section, err := this.forumSectionsDao.SelectRowById(ctx, id)
	if err != nil {
		return false, err
	}
	if section == nil {
		return false, nil
	}
	return section.Visibility == "Public", nil
}

func (this *forumServiceImpl) UpdateSectionUserConfig(ctx context.Context, id string, name string, ordering float32) error {
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

func (this *forumServiceImpl) UpdateForumUserConfig(ctx context.Context, id string, sectionId string, name string, description string, ordering float32, visibility string) error {
	if middleutil.ContextHasPermission(ctx, perms.Admin_Forums_Structure_Edit) == false {
		return errors.New("permission denied")
	}

	userId, ok := ctx.Value(contextkeys.UserId).(string)
	if !ok {
		return errors.New("cannot update forum: no user id in context")
	}

	err := this.forumsDao.UpdateUserConfigById(ctx, id, sectionId, name, description, ordering, visibility)
	if err != nil {
		return err
	}

	err = this.systemLogDao.InsertRow(ctx, userId, id, fmt.Sprintf("Updated forum user config: %s", name))
	if err != nil {
		return err
	}

	return nil
}
