// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.db;

import lombok.RequiredArgsConstructor;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.stereotype.Component;

@Component
@RequiredArgsConstructor
public class MetaDAO
{
	private final JdbcTemplate jdbcTemplate;

	public void createDatabaseObjects()
	{
		final String SQL_FN_CREATED_UPDATED_BEFORE_INSERT = """
			CREATE FUNCTION fn_created_updated_before_insert()
			RETURNS TRIGGER AS $$
			BEGIN
				NEW._created := now();
				NEW._updated := NEW._created;
				RETURN NEW;
			END;
			$$ LANGUAGE plpgsql;
			""";
		jdbcTemplate.execute(SQL_FN_CREATED_UPDATED_BEFORE_INSERT);
		final String SQL_FN_CREATED_UPDATED_BEFORE_UPDATE = """
			CREATE FUNCTION fn_created_updated_before_update()
			RETURNS TRIGGER AS $$
			BEGIN
				NEW._created := OLD._created;
				NEW._updated := now();
				RETURN NEW;
			END;
			$$ LANGUAGE plpgsql;
			""";
		jdbcTemplate.execute(SQL_FN_CREATED_UPDATED_BEFORE_UPDATE);
	}

	public boolean doesTableExist(String table)
	{
		return doesTableExist("public", table);
	}

	public boolean doesTableExist(String schema, String table)
	{
		Long count = jdbcTemplate.queryForObject("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = ? AND table_name = ?", Long.class, schema, table);
		return count != null && count > 0;
	}

	public String getDatabaseSize()
	{
		return jdbcTemplate.queryForObject("SELECT pg_size_pretty(pg_database_size(CURRENT_CATALOG))", String.class);
	}

	public String getServerVersion()
	{
		return jdbcTemplate.queryForObject("SHOW server_version", String.class);
	}
}
