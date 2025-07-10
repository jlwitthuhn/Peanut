// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.db;

import io.github.jlwitthuhn.peanut.err.DBCreationDependencyNotSatisfiedException;
import io.github.jlwitthuhn.peanut.service.DatabaseService;
import lombok.RequiredArgsConstructor;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.stereotype.Component;

import java.util.Collection;

@Component
@RequiredArgsConstructor
public class GroupMembershipDAO
{
	public static final String TABLE_NAME = "group_membership";

	private final DatabaseService dbService;

	private final JdbcTemplate jdbcTemplate;

	public void createDatabaseObjects() throws DBCreationDependencyNotSatisfiedException
	{
		if (dbService.doesTableExist(TABLE_NAME))
		{
			throw new DBCreationDependencyNotSatisfiedException("Table '" + TABLE_NAME + "' cannot be created because it already exists");
		}
		if (!dbService.doesTableExist(GroupDAO.TABLE_NAME))
		{
			throw new DBCreationDependencyNotSatisfiedException("Table '" + TABLE_NAME + "' requires that table '" + GroupDAO.TABLE_NAME + "' exists");
		}
		if (!dbService.doesTableExist(UserDAO.TABLE_NAME))
		{
			throw new DBCreationDependencyNotSatisfiedException("Table '" + TABLE_NAME + "' requires that table '" + UserDAO.TABLE_NAME + "' exists");
		}
		final String SQL_TABLE = """
			CREATE TABLE group_membership (
			    user_id BIGINT REFERENCES users(id) NOT NULL,
			    group_id BIGINT REFERENCES groups(id) NOT NULL,
			    _created TIMESTAMP WITH TIME ZONE NOT NULL,
			    _updated TIMESTAMP WITH TIME ZONE NOT NULL,
			    PRIMARY KEY (user_id, group_id)
			);
			""";
		jdbcTemplate.execute(SQL_TABLE);
		final String SQL_TRIGGER_BEFORE_INSERT = """
			CREATE TRIGGER
				group_membership_trigger_created_updated_before_insert
			BEFORE INSERT ON
				group_membership
			FOR EACH ROW EXECUTE FUNCTION
				fn_created_updated_before_insert();
			""";
		jdbcTemplate.execute(SQL_TRIGGER_BEFORE_INSERT);
		final String SQL_TRIGGER_BEFORE_UPDATE = """
			CREATE TRIGGER
				group_membership_trigger_created_updated_before_update
			BEFORE UPDATE ON
				group_membership
			FOR EACH ROW EXECUTE FUNCTION
				fn_created_updated_before_update();
			""";
		jdbcTemplate.execute(SQL_TRIGGER_BEFORE_UPDATE);
	}

	public void insertGroupsForUser(long userId, Collection<Long> groupIds)
	{
		final String SQL = "INSERT INTO group_membership (user_id, group_id) VALUES (?, ?)";
		for (long groupId : groupIds)
		{
			jdbcTemplate.update(SQL, userId, groupId);
		}
	}
}
