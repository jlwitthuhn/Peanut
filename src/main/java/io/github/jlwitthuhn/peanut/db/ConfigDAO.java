// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.db;

import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.stereotype.Component;

@Component
public class ConfigDAO
{
	private final JdbcTemplate jdbcTemplate;

	public ConfigDAO(JdbcTemplate jdbcTemplate)
	{
		this.jdbcTemplate = jdbcTemplate;
	}

	public Long getLong(String name)
	{
		return jdbcTemplate.queryForObject("SELECT value FROM config_int WHERE name = ?", Long.class, name);
	}
}
