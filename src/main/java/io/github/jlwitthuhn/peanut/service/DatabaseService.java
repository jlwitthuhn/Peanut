// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.service;

import io.github.jlwitthuhn.peanut.db.MetaDAO;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;

@Service
@RequiredArgsConstructor
public class DatabaseService
{
	private final MetaDAO metaDAO;

	public boolean doesTableExist(String tableName)
	{
		return metaDAO.doesTableExist(tableName);
	}

	public String getDatabaseSize()
	{
		return metaDAO.getDatabaseSize();
	}

	public String getServerVersion()
	{
		return metaDAO.getServerVersion();
	}
}
