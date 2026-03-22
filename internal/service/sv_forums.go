// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"net/http"
	"peanut/internal/data"
)

type ForumsService interface {
	CreateSection(req *http.Request, name string, ordering int) error
}

func NewForumsService(forumSectionsDao data.ForumSectionsDao) ForumsService {
	return &forumsServiceImpl{forumSectionsDao: forumSectionsDao}
}

type forumsServiceImpl struct {
	forumSectionsDao data.ForumSectionsDao
}

func (this *forumsServiceImpl) CreateSection(req *http.Request, name string, ordering int) error {
	err := this.forumSectionsDao.InsertRow(req, name, ordering)
	return err
}
