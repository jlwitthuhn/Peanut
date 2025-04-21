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

	public UserDAO(final JdbcTemplate jdbcTemplate)
	{
		this.jdbcTemplate = jdbcTemplate;
	}

	public UserRow getByName(final String name)
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
