// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"context"
	"errors"
	"peanut/internal/data"
	"peanut/internal/keynames/contextkeys"
)

type ForumThreadService interface {
	CreateThread(ctx context.Context, forumId string, title string, message string) (string, error)
	GetForumThreadSummaryViewRowByForumIdPublic(ctx context.Context, forumId string) ([]data.ForumThreadSummaryViewRow, error)
}

func NewForumThreadService(forumPostsDao data.ForumPostsDao, forumThreadsDao data.ForumThreadsDao, forumThreadSummaryDvao data.ForumThreadSummaryDvao) ForumThreadService {
	return &forumThreadServiceImpl{forumPostsDao: forumPostsDao, forumThreadsDao: forumThreadsDao, forumThreadSummaryDvao: forumThreadSummaryDvao}
}

type forumThreadServiceImpl struct {
	forumPostsDao          data.ForumPostsDao
	forumThreadsDao        data.ForumThreadsDao
	forumThreadSummaryDvao data.ForumThreadSummaryDvao
}

func (this *forumThreadServiceImpl) CreateThread(ctx context.Context, forumId string, title string, message string) (string, error) {
	userId, ok := ctx.Value(contextkeys.UserId).(string)
	if !ok {
		return "", errors.New("cannot create thread: no user id in context")
	}

	threadId, err := this.forumThreadsDao.InsertRow(ctx, forumId, userId, title, "Public")
	if err != nil {
		return "", err
	}

	_, err = this.forumPostsDao.InsertRow(ctx, threadId, userId, message)
	if err != nil {
		return "", err
	}

	return threadId, nil
}

func (this *forumThreadServiceImpl) GetForumThreadSummaryViewRowByForumIdPublic(ctx context.Context, forumId string) ([]data.ForumThreadSummaryViewRow, error) {
	result, err := this.forumThreadSummaryDvao.SelectForumThreadSummaryViewRowByForumIdPublic(ctx, forumId)
	return result, err
}
