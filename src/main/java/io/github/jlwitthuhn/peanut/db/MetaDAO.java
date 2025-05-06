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

	public boolean doesTableExist(String table)
	{
		return doesTableExist("public", table);
	}

	public boolean doesTableExist(String schema, String table)
	{
		Long count = jdbcTemplate.queryForObject("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = ? AND table_name = ?", Long.class, schema, table);
		return count != null && count > 0;
	}
}
