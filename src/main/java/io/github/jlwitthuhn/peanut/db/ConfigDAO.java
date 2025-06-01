// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.db;

import io.github.jlwitthuhn.peanut.err.DBCreationDependencyNotSatisfiedException;
import lombok.RequiredArgsConstructor;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.stereotype.Component;

@Component
@RequiredArgsConstructor
public class ConfigDAO
{
	public final String TABLE_NAME_INT = "config_int";
	public final String TABLE_NAME_STRING = "config_string";

	private final JdbcTemplate jdbcTemplate;

	private final MetaDAO metaDAO;

	public void createDatabaseObjects() throws DBCreationDependencyNotSatisfiedException
	{
		if (metaDAO.doesTableExist(TABLE_NAME_INT))
		{
			throw new DBCreationDependencyNotSatisfiedException("Table '" + TABLE_NAME_INT + "' cannot be created because it already exists");
		}
		if (metaDAO.doesTableExist(TABLE_NAME_STRING))
		{
			throw new DBCreationDependencyNotSatisfiedException("Table '" + TABLE_NAME_STRING + "' cannot be created because it already exists");
		}
		{
			final String SQL_TABLE = """
				CREATE TABLE config_int (
					name VARCHAR(255) PRIMARY KEY,
					value BIGINT NOT NULL,
					_created TIMESTAMP WITH TIME ZONE NOT NULL,
					_updated TIMESTAMP WITH TIME ZONE NOT NULL
				);
				""";
			jdbcTemplate.execute(SQL_TABLE);
			final String SQL_TRIGGER_BEFORE_INSERT = """
				CREATE TRIGGER
					config_int_trigger_created_updated_before_insert
				BEFORE INSERT ON
					config_int
				FOR EACH ROW EXECUTE FUNCTION
					fn_created_updated_before_insert();
				""";
			jdbcTemplate.execute(SQL_TRIGGER_BEFORE_INSERT);
			final String SQL_TRIGGER_BEFORE_UPDATE = """
				CREATE TRIGGER
					config_int_trigger_created_updated_before_update
				BEFORE UPDATE ON
					config_int
				FOR EACH ROW EXECUTE FUNCTION
					fn_created_updated_before_update();
				""";
			jdbcTemplate.execute(SQL_TRIGGER_BEFORE_UPDATE);
		}
		{
			final String SQL_TABLE = """
				CREATE TABLE config_string (
					name VARCHAR(255) PRIMARY KEY,
					value TEXT NOT NULL,
					_created TIMESTAMP WITH TIME ZONE NOT NULL,
					_updated TIMESTAMP WITH TIME ZONE NOT NULL
				);
				""";
			jdbcTemplate.execute(SQL_TABLE);
			final String SQL_TRIGGER_BEFORE_INSERT = """
				CREATE TRIGGER
					config_string_trigger_created_updated_before_insert
				BEFORE INSERT ON
					config_string
				FOR EACH ROW EXECUTE FUNCTION
					fn_created_updated_before_insert();
				""";
			jdbcTemplate.execute(SQL_TRIGGER_BEFORE_INSERT);
			final String SQL_TRIGGER_BEFORE_UPDATE = """
				CREATE TRIGGER
					config_string_trigger_created_updated_before_update
				BEFORE UPDATE ON
					config_string
				FOR EACH ROW EXECUTE FUNCTION
					fn_created_updated_before_update();
				""";
			jdbcTemplate.execute(SQL_TRIGGER_BEFORE_UPDATE);
		}
	}

	public Long selectLongByName(String name)
	{
		final String SQL = "SELECT value FROM config_int WHERE name = ?";
		return jdbcTemplate.queryForObject(SQL, Long.class, name);
	}

	public String selectStringByName(String name)
	{
		final String SQL = "SELECT value FROM config_string WHERE name = ?";
		return jdbcTemplate.queryForObject(SQL, String.class, name);
	}

	public void upsertLongByName(String name, long value)
	{
		final String SQL = """
			INSERT INTO
				config_int (name, value)
			VALUES
				(?, ?)
			ON CONFLICT (name)
				DO UPDATE SET value = EXCLUDED.value;
			""";
		jdbcTemplate.update(SQL, name, value);
	}

	public void upsertStringByName(String name, String value)
	{
		final String SQL = """
			INSERT INTO
				config_string (name, value)
			VALUES
				(?, ?)
			ON CONFLICT (name)
				DO UPDATE SET value = EXCLUDED.value;
			""";
		jdbcTemplate.update(SQL, name, value);
	}
}
