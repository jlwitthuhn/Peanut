// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.db;

import io.github.jlwitthuhn.peanut.err.DBCreationDependencyNotSatisfiedException;
import io.github.jlwitthuhn.peanut.err.DBObjectAlreadyExistsException;
import lombok.RequiredArgsConstructor;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.stereotype.Component;

import java.util.Collection;

@Component
@RequiredArgsConstructor
public class UserAuthorityDAO
{
	public static final String TABLE_NAME = "user_authorities";

	private final JdbcTemplate jdbcTemplate;
	private final MetaDAO metaDAO;

	public void createDatabaseObjects() throws DBObjectAlreadyExistsException, DBCreationDependencyNotSatisfiedException
	{
		if (metaDAO.doesTableExist(TABLE_NAME))
		{
			throw new DBObjectAlreadyExistsException();
		}
		if (!metaDAO.doesTableExist(AuthorityDAO.TABLE_NAME))
		{
			throw new DBCreationDependencyNotSatisfiedException("Table 'user_authorities' requires that table 'authorities' exists");
		}
		if (!metaDAO.doesTableExist(UserDAO.TABLE_NAME))
		{
			throw new DBCreationDependencyNotSatisfiedException("Table 'user_authorities' requires that table 'users' exists");
		}
		final String SQL_TABLE = """
			CREATE TABLE user_authorities (
			    user_id BIGINT REFERENCES users(id),
			    authority_id BIGINT REFERENCES authorities(id),
			    _created TIMESTAMP WITH TIME ZONE NOT NULL,
			    _updated TIMESTAMP WITH TIME ZONE NOT NULL,
			    PRIMARY KEY (user_id, authority_id)
			);
			""";
		jdbcTemplate.execute(SQL_TABLE);
		final String SQL_TRIGGER_BEFORE_INSERT = """
			CREATE TRIGGER
				user_authorities_trigger_created_updated_before_insert
			BEFORE INSERT ON
				user_authorities
			FOR EACH ROW EXECUTE FUNCTION
				fn_created_updated_before_insert();
			""";
		jdbcTemplate.execute(SQL_TRIGGER_BEFORE_INSERT);
		final String SQL_TRIGGER_BEFORE_UPDATE = """
			CREATE TRIGGER
				user_authorities_trigger_created_updated_before_update
			BEFORE UPDATE ON
				user_authorities
			FOR EACH ROW EXECUTE FUNCTION
				fn_created_updated_before_update();
			""";
		jdbcTemplate.execute(SQL_TRIGGER_BEFORE_UPDATE);
	}

	public void insertAuthoritiesForUser(long userId, Collection<Long> authorityIds)
	{
		final String SQL = "INSERT INTO user_authorities (user_id, authority_id) VALUES (?, ?)";
		for (long authorityId : authorityIds)
		{
			jdbcTemplate.update(SQL, userId, authorityId);
		}
	}
}
