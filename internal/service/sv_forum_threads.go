// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"context"
	"errors"
	"peanut/internal/data"
	"peanut/internal/keynames/contextkeys"
	"peanut/internal/middleutil"
	"peanut/internal/security/perms"
)

const ForumThreadsPerPage = 20

type ForumThreadListing struct {
	Threads     []data.ForumThreadSummaryViewRow
	CurrentPage int
	TotalPages  int
}

type ForumThreadService interface {
	AddThreadPost(ctx context.Context, threadId string, message string) (string, error)
	CanPostReplyInThread(ctx context.Context, forumId string) (bool, error)
	CanPostThreadInForum(ctx context.Context, forumId string) (bool, error)
	CreateThread(ctx context.Context, forumId string, title string, message string) (string, error)
	GetForumThreadListingPublic(ctx context.Context, forumId string, pageNum int) (*ForumThreadListing, error)
	GetThreadRowById(ctx context.Context, id string) (*data.ForumThreadRow, error)
	GetPostsByThreadId(ctx context.Context, threadId string) ([]data.ForumPostListingViewRow, error)
}

func NewForumThreadService(forumService ForumService, forumPostsDao data.ForumPostsDao, forumThreadsDao data.ForumThreadsDao, forumThreadSummaryDvao data.ForumThreadSummaryDvao, forumPostListingDvao data.ForumPostListingDvao) ForumThreadService {
	return &forumThreadServiceImpl{forumService: forumService, forumPostsDao: forumPostsDao, forumThreadsDao: forumThreadsDao, forumThreadSummaryDvao: forumThreadSummaryDvao, forumPostListingDvao: forumPostListingDvao}
}

type forumThreadServiceImpl struct {
	forumService           ForumService
	forumPostsDao          data.ForumPostsDao
	forumThreadsDao        data.ForumThreadsDao
	forumThreadSummaryDvao data.ForumThreadSummaryDvao
	forumPostListingDvao   data.ForumPostListingDvao
}

func (this *forumThreadServiceImpl) AddThreadPost(ctx context.Context, threadId string, message string) (string, error) {
	userId, ok := ctx.Value(contextkeys.UserId).(string)
	if !ok {
		return "", errors.New("cannot add post: no user id in context")
	}

	postId, err := this.forumPostsDao.InsertRow(ctx, threadId, userId, message)
	if err != nil {
		return "", err
	}

	return postId, nil
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

func (this *forumThreadServiceImpl) CanPostThreadInForum(ctx context.Context, forumId string) (bool, error) {
	if !middleutil.ContextHasPermission(ctx, perms.Forum_Thread_Post) {
		return false, nil
	}
	return this.forumService.IsForumReadable(ctx, forumId)
}

func (this *forumThreadServiceImpl) CanPostReplyInThread(ctx context.Context, forumId string) (bool, error) {
	if !middleutil.ContextHasPermission(ctx, perms.Forum_Thread_Reply) {
		return false, nil
	}
	return this.forumService.IsForumReadable(ctx, forumId)
}

func (this *forumThreadServiceImpl) GetForumThreadListingPublic(ctx context.Context, forumId string, pageNum int) (*ForumThreadListing, error) {
	totalThreads, err := this.forumThreadSummaryDvao.CountForumThreadSummaryViewRowByForumIdPublic(ctx, forumId)
	if err != nil {
		return nil, err
	}

	totalPages := int((totalThreads + ForumThreadsPerPage - 1) / ForumThreadsPerPage)
	if totalPages < 1 {
		totalPages = 1
	}

	if pageNum < 1 || pageNum > totalPages {
		return nil, nil
	}

	beginIndex := (pageNum - 1) * ForumThreadsPerPage
	threads, err := this.forumThreadSummaryDvao.SelectPageForumThreadSummaryViewRowByForumIdPublic(ctx, forumId, beginIndex, ForumThreadsPerPage)
	if err != nil {
		return nil, err
	}

	return &ForumThreadListing{
		Threads:     threads,
		CurrentPage: pageNum,
		TotalPages:  totalPages,
	}, nil
}

func (this *forumThreadServiceImpl) GetThreadRowById(ctx context.Context, id string) (*data.ForumThreadRow, error) {
	return this.forumThreadsDao.SelectRowById(ctx, id)
}

func (this *forumThreadServiceImpl) GetPostsByThreadId(ctx context.Context, threadId string) ([]data.ForumPostListingViewRow, error) {
	return this.forumPostListingDvao.SelectByThreadId(ctx, threadId)
}
