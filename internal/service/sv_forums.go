// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"errors"
	"net/http"
	"peanut/internal/data"
	"peanut/internal/middleutil"
	"peanut/internal/security/perms"
)

type ForumsService interface {
	CreateSection(req *http.Request, name string, ordering float32) error
	GetAllSectionRows(req *http.Request) ([]data.ForumSectionRow, error)
}

func NewForumsService(forumSectionsDao data.ForumSectionsDao) ForumsService {
	return &forumsServiceImpl{forumSectionsDao: forumSectionsDao}
}

type forumsServiceImpl struct {
	forumSectionsDao data.ForumSectionsDao
}

func (this *forumsServiceImpl) CreateSection(req *http.Request, name string, ordering float32) error {
	if middleutil.RequestHasPermission(req, perms.Admin_Forums_Structure_Edit) == false {
		return errors.New("permission denied")
	}

	err := this.forumSectionsDao.InsertRow(req, name, ordering)
	return err
}

func (this *forumsServiceImpl) GetAllSectionRows(req *http.Request) ([]data.ForumSectionRow, error) {
	result, err := this.forumSectionsDao.SelectRowAll(req)
	return result, err
}
