// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.db;

import io.github.jlwitthuhn.peanut.err.AuthorityNotFoundException;
import io.github.jlwitthuhn.peanut.err.DBCreationDependencyNotSatisfiedException;
import io.github.jlwitthuhn.peanut.model.db.AuthorityRow;
import io.github.jlwitthuhn.peanut.model.db.AuthorityRowMapper;
import lombok.RequiredArgsConstructor;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.stereotype.Component;

import java.util.ArrayList;
import java.util.Collection;
import java.util.Collections;
import java.util.HashSet;
import java.util.List;
import java.util.Map;

@Component
@RequiredArgsConstructor
public class AuthorityDAO
{
	public static final String TABLE_NAME = "authorities";

	private final JdbcTemplate jdbcTemplate;
	private final MetaDAO metaDAO;

	public void createDatabaseObjects() throws DBCreationDependencyNotSatisfiedException
	{
		if (metaDAO.doesTableExist(TABLE_NAME))
		{
			throw new DBCreationDependencyNotSatisfiedException("Table '" + TABLE_NAME + "' cannot be created because it already exists");
		}
		final String SQL_TABLE = """
			CREATE TABLE authorities (
			    id BIGSERIAL PRIMARY KEY,
			    name VARCHAR(127) UNIQUE NOT NULL,
			    description VARCHAR(255) NOT NULL,
			    system_owned BOOLEAN NOT NULL,
			    _created TIMESTAMP WITH TIME ZONE NOT NULL,
			    _updated TIMESTAMP WITH TIME ZONE NOT NULL
			);
			""";
		jdbcTemplate.execute(SQL_TABLE);
		final String SQL_TRIGGER_BEFORE_INSERT = """
			CREATE TRIGGER
				authorities_trigger_created_updated_before_insert
			BEFORE INSERT ON
				authorities
			FOR EACH ROW EXECUTE FUNCTION
				fn_created_updated_before_insert();
			""";
		jdbcTemplate.execute(SQL_TRIGGER_BEFORE_INSERT);
		final String SQL_TRIGGER_BEFORE_UPDATE = """
			CREATE TRIGGER
				authorities_trigger_created_updated_before_update
			BEFORE UPDATE ON
				authorities
			FOR EACH ROW EXECUTE FUNCTION
				fn_created_updated_before_update();
			""";
		jdbcTemplate.execute(SQL_TRIGGER_BEFORE_UPDATE);
	}

	public Collection<Long> getIdsFromNames(Collection<String> names) throws AuthorityNotFoundException
	{
		final String namesQuestions = String.join(",", Collections.nCopies(names.size(), "?"));
		final String SQL = "SELECT id, name FROM authorities WHERE name IN (" + namesQuestions + ")";
		List<Map<String, Object>> result = jdbcTemplate.queryForList(SQL, names.toArray());
		HashSet<String> remainingNames = new HashSet<>(names);
		ArrayList<Long> ids = new ArrayList<Long>();
		for (Map<String, Object> row : result)
		{
			boolean nameValid = row.containsKey("name") && row.get("name") instanceof String;
			boolean idValid = row.containsKey("id") && row.get("id") instanceof Long;
			if (nameValid && idValid)
			{
				String rowName = (String) row.get("name");
				if (!remainingNames.contains(rowName))
				{
					throw new RuntimeException("Found an authority name that was not requested");
				}
				remainingNames.remove(rowName);

				Long rowId = (Long) row.get("id");
				ids.add(rowId);
			}
		}
		if (!remainingNames.isEmpty())
		{
			throw new AuthorityNotFoundException("Could not find authorities: " + remainingNames.toString());
		}
		return ids;
	}

	public List<AuthorityRow> selectAll()
	{
		final String SQL = "SELECT id, name, description, system_owned, _created, _updated FROM authorities ORDER BY id";
		return jdbcTemplate.query(SQL, new AuthorityRowMapper());
	}

	public void insertRow(String name, String description, boolean systemOwned)
	{
		final String SQL = "INSERT INTO authorities (name, description, system_owned) VALUES (?, ?, ?)";
		jdbcTemplate.update(SQL, name, description, systemOwned);
	}
}
