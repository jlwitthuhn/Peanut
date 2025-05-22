// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.db;

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
}
