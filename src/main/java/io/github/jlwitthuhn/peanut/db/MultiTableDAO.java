// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.db;

import lombok.RequiredArgsConstructor;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.stereotype.Component;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;

@Component
@RequiredArgsConstructor
public class MultiTableDAO
{
	private final JdbcTemplate jdbcTemplate;

	public ArrayList<GrantedAuthority> getAuthoritiesByUserId(long userId)
	{
		final String sql = """
			SELECT
				authorities.name AS name
			FROM
				user_authorities INNER JOIN authorities
				ON user_authorities.authority_id = authorities.id
			WHERE
				user_authorities.user_id = ?;
			""";
		List<Map<String, Object>> result = jdbcTemplate.queryForList(sql, userId);
		ArrayList<GrantedAuthority> authorities = new ArrayList<>();
		for (Map<String, Object> row : result)
		{
			if (row.containsKey("name"))
			{
				if (row.get("name") instanceof String nameString)
				{
					authorities.add(new SimpleGrantedAuthority(nameString));
				}
			}
		}
		return authorities;
	}
}
