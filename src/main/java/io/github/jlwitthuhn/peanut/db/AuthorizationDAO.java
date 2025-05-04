// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.db;

import io.github.jlwitthuhn.peanut.err.TableAlreadyExistsException;
import lombok.RequiredArgsConstructor;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.stereotype.Component;

@Component
@RequiredArgsConstructor
public class AuthorizationDAO
{
	public static final String TABLE_NAME = "authorizations";

	private final JdbcTemplate jdbcTemplate;
	private final InformationSchemaDAO informationSchemaDAO;

	public void createDatabaseObjects() throws TableAlreadyExistsException
	{
		if (informationSchemaDAO.doesTableExist(TABLE_NAME))
		{
			throw new TableAlreadyExistsException();
		}
		final String SQL = """
			CREATE TABLE authorizations (
			    id BIGSERIAL PRIMARY KEY,
			    name VARCHAR(127) UNIQUE NOT NULL,
			    system_owned BOOLEAN NOT NULL
			);
			""";
		jdbcTemplate.execute(SQL);
	}

	public void insertRow(String name, boolean systemOwned)
	{
		final String SQL = "INSERT INTO authorizations (name, system_owned) VALUES (?, ?)";
		jdbcTemplate.update(SQL, name, systemOwned);
	}
}
