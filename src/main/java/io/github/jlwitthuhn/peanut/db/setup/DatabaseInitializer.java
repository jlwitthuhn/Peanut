// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.db.setup;

import io.github.jlwitthuhn.peanut.cfg.ConfigKeyNames;
import io.github.jlwitthuhn.peanut.db.AuthorityDAO;
import io.github.jlwitthuhn.peanut.db.ConfigDAO;
import io.github.jlwitthuhn.peanut.db.UserAuthorityDAO;
import io.github.jlwitthuhn.peanut.db.UserDAO;
import lombok.RequiredArgsConstructor;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.security.config.annotation.authentication.configuration.GlobalAuthenticationConfigurerAdapter;
import org.springframework.stereotype.Component;

@Component
@RequiredArgsConstructor
public class DatabaseInitializer
{
	private final AuthorityDAO authorityDAO;
	private final ConfigDAO configDAO;
	private final UserDAO userDAO;
	private final UserAuthorityDAO userAuthorityDAO;

	private final Logger logger = LoggerFactory.getLogger(DatabaseInitializer.class);
	private final GlobalAuthenticationConfigurerAdapter enableGlobalAuthenticationAutowiredConfigurer;

	public boolean doInit()
	{
		boolean success = createDatabaseObjects();
		success = success && insertDefaultAuthorities();
		success = success && insertDefaultConfig();
		return success;
	}

	private boolean createDatabaseObjects()
	{
		try
		{
			// Order here matters as foreign keys create dependencies between tables
			authorityDAO.createDatabaseObjects();
			configDAO.createDatabaseObjects();
			userDAO.createDatabaseObjects();
			userAuthorityDAO.createDatabaseObjects();
		}
		catch (Exception e)
		{
			logger.error("Caught exception while creating tables", e);
			return false;
		}
		return true;
	}

	private boolean insertDefaultAuthorities()
	{
		try
		{
			authorityDAO.insertRow("ROLE_USER", true);
			authorityDAO.insertRow("ROLE_ADMIN", true);
		}
		catch (Exception e)
		{
			logger.error("Caught exception while inserting authority", e);
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
