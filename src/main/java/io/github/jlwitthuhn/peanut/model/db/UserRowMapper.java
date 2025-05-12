// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.model.db;

import org.springframework.jdbc.core.RowMapper;

import java.sql.ResultSet;
import java.sql.SQLException;
import java.time.OffsetDateTime;

public class UserRowMapper implements RowMapper<UserRow>
{
	@Override
	public UserRow mapRow(ResultSet rs, int rowNum) throws SQLException {
		long id = rs.getLong("id");
		String displayName = rs.getString("display_name");
		String email = rs.getString("email");
		String password = rs.getString("password");
		OffsetDateTime created = rs.getObject("_created", OffsetDateTime.class);
		OffsetDateTime updated = rs.getObject("_updated", OffsetDateTime.class);
		return new UserRow(id, displayName, email, password, created, updated);
	}
}
