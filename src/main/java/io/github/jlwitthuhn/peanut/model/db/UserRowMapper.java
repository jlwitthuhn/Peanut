package io.github.jlwitthuhn.peanut.model.db;

import org.springframework.jdbc.core.RowMapper;

import java.sql.ResultSet;
import java.sql.SQLException;

public class UserRowMapper implements RowMapper<UserRow>
{
	@Override
	public UserRow mapRow(ResultSet rs, int rowNum) throws SQLException {
		long id = rs.getLong("id");
		String name = rs.getString("name");
		String password = rs.getString("password");
		return new UserRow(id, name, password);
	}
}
