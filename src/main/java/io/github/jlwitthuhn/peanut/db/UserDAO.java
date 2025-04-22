// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.db;

import io.github.jlwitthuhn.peanut.model.db.UserRow;
import io.github.jlwitthuhn.peanut.model.db.UserRowMapper;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.stereotype.Component;

import java.util.List;

@Component
public class UserDAO
{
	private final JdbcTemplate jdbcTemplate;

	public UserDAO(JdbcTemplate jdbcTemplate)
	{
		this.jdbcTemplate = jdbcTemplate;
	}

	public void createRow(String name, String password)
	{
		final String SQL = "INSERT INTO users (name, password) VALUES (?, ?)";
		jdbcTemplate.update(SQL, name, password);
	}

	public UserRow getByName(String name)
	{
		final String SQL = "SELECT id, name, password FROM users WHERE name = ?";
		final List<UserRow> results = jdbcTemplate.query(SQL, new UserRowMapper(), name);
		if (results.size() != 1)
		{
			return null;
		}
		return results.getFirst();
	}
}
