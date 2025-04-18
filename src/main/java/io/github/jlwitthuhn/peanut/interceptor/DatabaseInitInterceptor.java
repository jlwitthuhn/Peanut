// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.interceptor;

import io.github.jlwitthuhn.peanut.err.DatabaseNotInitializedException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.jdbc.BadSqlGrammarException;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.lang.NonNull;
import org.springframework.stereotype.Component;
import org.springframework.web.servlet.HandlerInterceptor;

@Component
public class DatabaseInitInterceptor implements HandlerInterceptor
{
	JdbcTemplate jdbcTemplate;

	@Autowired
	public DatabaseInitInterceptor(JdbcTemplate jdbcTemplate)
	{
		this.jdbcTemplate = jdbcTemplate;
	}

	@Override
	public boolean preHandle(@NonNull HttpServletRequest request, @NonNull HttpServletResponse response, @NonNull Object handler)
	{
		// These pages do not require the database to be initialized
		if (request.getRequestURI().equals("/design") || request.getRequestURI().equals("/setup"))
		{
			return true;
		}

		try
		{
			String result = jdbcTemplate.queryForObject("SELECT value FROM config WHERE value = 'schema_version'", String.class);
		}
		catch (BadSqlGrammarException ex)
		{
			throw new DatabaseNotInitializedException();
		}

		return true;
	}
}
