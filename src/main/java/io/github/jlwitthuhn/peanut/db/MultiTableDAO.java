// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.db;

import io.github.jlwitthuhn.peanut.model.db.UserRow;
import io.github.jlwitthuhn.peanut.model.db.UserRowMapper;
import lombok.RequiredArgsConstructor;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.stereotype.Component;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;

@Component
@RequiredArgsConstructor
public class MultiTableDAO
{
	private final JdbcTemplate jdbcTemplate;

	public List<String> getGroupsByUserId(long userId)
	{
		final String sql = """
			SELECT
				groups.name AS name
			FROM
				group_membership INNER JOIN groups
				ON group_membership.group_id = groups.id
			WHERE
				group_membership.user_id = ?;
			""";
		List<Map<String, Object>> result = jdbcTemplate.queryForList(sql, userId);
		ArrayList<String> groups = new ArrayList<>();
		for (Map<String, Object> row : result)
		{
			if (row.containsKey("name"))
			{
				if (row.get("name") instanceof String nameString)
				{
					groups.add(nameString);
				}
			}
		}
		return groups;
	}

	public List<UserRow> getUsersByGroupName(String groupName)
	{
		final String sql = """
			WITH
			group_id AS (
				SELECT id FROM groups WHERE name = ?
			),
			user_ids_in_group AS (
				SELECT user_id FROM group_membership WHERE group_id IN (SELECT * FROM group_id)
			)
			SELECT
				*
			FROM
				users
			WHERE
				id IN (SELECT * FROM user_ids_in_group)
			ORDER BY
			    id
			""";
		return jdbcTemplate.query(sql, new UserRowMapper(), groupName);
	}
}
