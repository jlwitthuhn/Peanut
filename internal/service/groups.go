// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"database/sql"
	"peanut/internal/data"
)

type GroupService interface {
	CreateGroup(tx *sql.Tx, name string, desc string, systemOwned bool) error
}

func NewGroupService() GroupService {
	return &groupServiceImpl{}
}

type groupServiceImpl struct{}

func (*groupServiceImpl) CreateGroup(tx *sql.Tx, name string, desc string, systemOwned bool) error {
	err := data.GroupDaoInst().InsertRow(tx, name, desc, systemOwned)
	return err
}
