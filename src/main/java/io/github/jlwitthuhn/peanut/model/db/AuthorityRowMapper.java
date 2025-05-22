// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.model.db;

import org.springframework.jdbc.core.RowMapper;

import java.sql.ResultSet;
import java.sql.SQLException;
import java.time.OffsetDateTime;

public class AuthorityRowMapper implements RowMapper<AuthorityRow>
{
	@Override
	public AuthorityRow mapRow(ResultSet rs, int rowNum) throws SQLException
	{
		long id = rs.getLong("id");
		String displayName = rs.getString("name");
		String description = rs.getString("description");
		boolean systemOwned = rs.getBoolean("system_owned");
		OffsetDateTime created = rs.getObject("_created", OffsetDateTime.class);
		OffsetDateTime updated = rs.getObject("_updated", OffsetDateTime.class);
		return new AuthorityRow(id, displayName, description, systemOwned, created, updated);
	}
}
