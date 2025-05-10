// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.db;

import io.github.jlwitthuhn.peanut.err.DBObjectAlreadyExistsException;
import lombok.RequiredArgsConstructor;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.stereotype.Component;

@Component
@RequiredArgsConstructor
public class ConfigDAO
{
	public final String TABLE_NAME_INT = "config_int";

	private final JdbcTemplate jdbcTemplate;

	private final MetaDAO metaDAO;

	public void createDatabaseObjects() throws DBObjectAlreadyExistsException
	{
		if (metaDAO.doesTableExist(TABLE_NAME_INT))
		{
			throw new DBObjectAlreadyExistsException();
		}
		final String SQL = """
			CREATE TABLE config_int (
				name VARCHAR(255) PRIMARY KEY,
				value BIGINT NOT NULL
			);
			""";
		jdbcTemplate.execute(SQL);
	}

	public Long getLong(String name)
	{
		final String SQL = "SELECT value FROM config_int WHERE name = ?";
		return jdbcTemplate.queryForObject(SQL, Long.class, name);
	}

	public void setLong(String name, long value)
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
}
