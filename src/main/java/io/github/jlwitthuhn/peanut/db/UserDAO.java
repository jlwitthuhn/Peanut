// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.db;

import io.github.jlwitthuhn.peanut.err.TableAlreadyExistsException;
import io.github.jlwitthuhn.peanut.model.db.UserRow;
import io.github.jlwitthuhn.peanut.model.db.UserRowMapper;
import lombok.RequiredArgsConstructor;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.stereotype.Component;

import java.util.List;

@Component
@RequiredArgsConstructor
public class UserDAO
{
	public final String TABLE_NAME = "users";

	private final JdbcTemplate jdbcTemplate;

	private final InformationSchemaDAO informationSchemaDAO;

	public void createDatabaseObjects() throws TableAlreadyExistsException
	{
		if (informationSchemaDAO.doesTableExist(TABLE_NAME))
		{
			throw new TableAlreadyExistsException();
		}
		final String SQL = """
			CREATE TABLE users (
			    id BIGSERIAL PRIMARY KEY,
			    display_name VARCHAR(255) UNIQUE NOT NULL,
			    email VARCHAR(255) UNIQUE NOT NULL,
			    password VARCHAR(255) NOT NULL
			);
			""";
		jdbcTemplate.execute(SQL);
	}

	public void insertRow(String displayName, String email, String password)
	{
		final String SQL = "INSERT INTO users (display_name, email, password) VALUES (?, ?, ?)";
		jdbcTemplate.update(SQL, displayName, email, password);
	}

	public UserRow selectRowByDisplayName(String displayName)
	{
		final String SQL = "SELECT id, display_name, email, password FROM users WHERE display_name = ?";
		final List<UserRow> results = jdbcTemplate.query(SQL, new UserRowMapper(), displayName);
		if (results.size() != 1)
		{
			return null;
		}
		return results.getFirst();
	}
}
