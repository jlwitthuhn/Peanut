// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.db.setup;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.stereotype.Component;

@Component
public class DatabaseInitializer
{
	private final JdbcTemplate jdbcTemplate;

	private final Logger logger = LoggerFactory.getLogger(DatabaseInitializer.class);

	public DatabaseInitializer(JdbcTemplate jdbcTemplate)
	{
		this.jdbcTemplate = jdbcTemplate;
	}

	public boolean doInit()
	{
		if (!createTables())
		{
			return false;
		}
		return insertConfig();
	}

	private boolean createTables()
	{
		final String CREATE_TABLE_CONFIG = """
			CREATE TABLE config_int (
				name VARCHAR(255) PRIMARY KEY,
				value BIGINT NOT NULL
			);
			""";
		final String CREATE_TABLE_USERS = """
			CREATE TABLE users (
			    id BIGSERIAL PRIMARY KEY,
			    name VARCHAR(255) UNIQUE NOT NULL,
			    password VARCHAR(255) NOT NULL
			);
			""";
		try
		{
			jdbcTemplate.execute(CREATE_TABLE_CONFIG);
			jdbcTemplate.execute(CREATE_TABLE_USERS);
		}
		catch (Exception e)
		{
			logger.error("Caught exception while creating tables", e);
			return false;
		}
		return true;
	}

	private boolean insertConfig()
	{
		final String INSERT_SCHEMA_VERSION = """
			INSERT INTO config_int (name, value) VALUES (?, ?);
			""";
		try
		{
			jdbcTemplate.update(INSERT_SCHEMA_VERSION, "schemaVersion", 1);
		}
		catch (Exception e)
		{
			logger.error("Caught exception while inserting config", e);
			return false;
		}
		return true;
	}
}
