// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"context"
	"peanut/internal/data"
)

type ForumThreadService interface {
	GetForumThreadSummaryViewRowByForumIdPublic(ctx context.Context, forumId string) ([]data.ForumThreadSummaryViewRow, error)
}

func NewForumThreadService(forumThreadSummaryDvao data.ForumThreadSummaryDvao) ForumThreadService {
	return &forumThreadServiceImpl{forumThreadSummaryDvao: forumThreadSummaryDvao}
}

type forumThreadServiceImpl struct {
	forumThreadSummaryDvao data.ForumThreadSummaryDvao
}

func (this *forumThreadServiceImpl) GetForumThreadSummaryViewRowByForumIdPublic(ctx context.Context, forumId string) ([]data.ForumThreadSummaryViewRow, error) {
	result, err := this.forumThreadSummaryDvao.SelectForumThreadSummaryViewRowByForumIdPublic(ctx, forumId)
	return result, err
}
