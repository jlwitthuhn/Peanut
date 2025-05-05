// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.db;

import io.github.jlwitthuhn.peanut.err.TableAlreadyExistsException;
import io.github.jlwitthuhn.peanut.err.TableCreationDependencyNotSatisfiedException;
import lombok.RequiredArgsConstructor;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.stereotype.Component;

import java.util.Collection;

@Component
@RequiredArgsConstructor
public class  UserAuthorizationDAO
{
	public static final String TABLE_NAME = "user_authorizations";

	private final JdbcTemplate jdbcTemplate;
	private final InformationSchemaDAO informationSchemaDAO;

	public void createDatabaseObjects() throws TableAlreadyExistsException, TableCreationDependencyNotSatisfiedException
	{
		if (informationSchemaDAO.doesTableExist(TABLE_NAME))
		{
			throw new TableAlreadyExistsException();
		}
		if (!informationSchemaDAO.doesTableExist(AuthorizationDAO.TABLE_NAME))
		{
			throw new TableCreationDependencyNotSatisfiedException("Table 'user_authorizations' requires that table 'authorizations' exists");
		}
		if (!informationSchemaDAO.doesTableExist(UserDAO.TABLE_NAME))
		{
			throw new TableCreationDependencyNotSatisfiedException("Table 'user_authorizations' requires that table 'users' exists");
		}
		final String SQL = """
			CREATE TABLE user_authorizations (
			    user_id BIGINT REFERENCES users(id),
			    authorization_id BIGINT REFERENCES authorizations(id),
			    PRIMARY KEY (user_id, authorization_id)
			);
			""";
		jdbcTemplate.execute(SQL);
	}

	public void insertAuthoritiesForUser(long userId, Collection<Long> authorityIds)
	{
		final String SQL = "INSERT INTO user_authorizations (user_id, authorization_id) VALUES (?, ?)";
		for (long authorityId : authorityIds)
		{
			jdbcTemplate.update(SQL, userId, authorityId);
		}
	}
}
