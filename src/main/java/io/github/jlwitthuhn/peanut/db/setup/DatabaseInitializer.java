// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.db.setup;

import io.github.jlwitthuhn.peanut.cfg.ConfigKeyNames;
import io.github.jlwitthuhn.peanut.db.ConfigDAO;
import io.github.jlwitthuhn.peanut.db.UserDAO;
import lombok.RequiredArgsConstructor;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Component;

@Component
@RequiredArgsConstructor
public class DatabaseInitializer
{
	private final ConfigDAO configDAO;
	private final UserDAO userDAO;

	private final Logger logger = LoggerFactory.getLogger(DatabaseInitializer.class);

	public boolean doInit()
	{
		if (!createDatabaseObjects())
		{
			return false;
		}
		return insertDefaultConfig();
	}

	private boolean createDatabaseObjects()
	{
		try
		{
			configDAO.createDatabaseObjects();
			userDAO.createDatabaseObjects();
		}
		catch (Exception e)
		{
			logger.error("Caught exception while creating tables", e);
			return false;
		}
		return true;
	}

	private boolean insertDefaultConfig()
	{
		try
		{
			configDAO.setLong(ConfigKeyNames.SCHEMA_VERSION, 1);
		}
		catch (Exception e)
		{
			logger.error("Caught exception while inserting config", e);
			return false;
		}
		return true;
	}
}
