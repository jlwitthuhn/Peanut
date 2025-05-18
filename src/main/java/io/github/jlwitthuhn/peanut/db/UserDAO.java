// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.db;

import io.github.jlwitthuhn.peanut.err.DBCreationDependencyNotSatisfiedException;
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
	public static final String TABLE_NAME = "users";

	private final JdbcTemplate jdbcTemplate;

	private final MetaDAO metaDAO;

	public void createDatabaseObjects() throws DBCreationDependencyNotSatisfiedException
	{
		if (metaDAO.doesTableExist(TABLE_NAME))
		{
			throw new DBCreationDependencyNotSatisfiedException("Table '" + TABLE_NAME + "' cannot be created because it already exists");
		}
		final String SQL_TABLE = """
			CREATE TABLE users (
			    id BIGSERIAL PRIMARY KEY,
			    display_name VARCHAR(255) UNIQUE NOT NULL,
			    email VARCHAR(255) UNIQUE NOT NULL,
			    password VARCHAR(255) NOT NULL,
			    _created TIMESTAMP WITH TIME ZONE NOT NULL,
			    _updated TIMESTAMP WITH TIME ZONE NOT NULL
			);
			""";
		jdbcTemplate.execute(SQL_TABLE);
		final String SQL_TRIGGER_BEFORE_INSERT = """
			CREATE TRIGGER
				users_trigger_created_updated_before_insert
			BEFORE INSERT ON
				users
			FOR EACH ROW EXECUTE FUNCTION
				fn_created_updated_before_insert();
			""";
		jdbcTemplate.execute(SQL_TRIGGER_BEFORE_INSERT);
		final String SQL_TRIGGER_BEFORE_UPDATE = """
			CREATE TRIGGER
				users_trigger_created_updated_before_update
			BEFORE UPDATE ON
				users
			FOR EACH ROW EXECUTE FUNCTION
				fn_created_updated_before_update();
			""";
		jdbcTemplate.execute(SQL_TRIGGER_BEFORE_UPDATE);
	}

	public void insertRow(String displayName, String email, String password)
	{
		final String SQL = "INSERT INTO users (display_name, email, password) VALUES (?, ?, ?)";
		jdbcTemplate.update(SQL, displayName, email, password);
	}

	public UserRow selectRowByDisplayName(String displayName)
	{
		final String SQL = "SELECT id, display_name, email, password, _created, _updated FROM users WHERE display_name = ?";
		final List<UserRow> results = jdbcTemplate.query(SQL, new UserRowMapper(), displayName);
		if (results.size() != 1)
		{
			return null;
		}
		return results.getFirst();
	}

	public UserRow selectRowByEmail(String email)
	{
		final String SQL = "SELECT id, display_name, email, password, _created, _updated FROM users WHERE email = ?";
		final List<UserRow> results = jdbcTemplate.query(SQL, new UserRowMapper(), email);
		if (results.size() != 1)
		{
			return null;
		}
		return results.getFirst();
	}

	public UserRow selectRowById(long id)
	{
		final String SQL = "SELECT id, display_name, email, password, _created, _updated FROM users WHERE id = ?";
		final List<UserRow> results = jdbcTemplate.query(SQL, new UserRowMapper(), id);
		if (results.size() != 1)
		{
			return null;
		}
		return results.getFirst();
	}
}
