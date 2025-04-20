// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.interceptor;

import io.github.jlwitthuhn.peanut.db.ConfigDAO;
import io.github.jlwitthuhn.peanut.err.BadDatabaseSchemaException;
import io.github.jlwitthuhn.peanut.err.DatabaseNotInitializedException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.springframework.jdbc.BadSqlGrammarException;
import org.springframework.lang.NonNull;
import org.springframework.stereotype.Component;
import org.springframework.web.servlet.HandlerInterceptor;

@Component
public class DatabaseInitInterceptor implements HandlerInterceptor
{
	ConfigDAO configDAO;

	public DatabaseInitInterceptor(ConfigDAO configDAO)
	{
		this.configDAO = configDAO;
	}

	@Override
	public boolean preHandle(@NonNull HttpServletRequest request, @NonNull HttpServletResponse response, @NonNull Object handler)
	{
		try
		{
			Long version = configDAO.getLong("schemaVersion");
			if (version == null || version != 1)
			{
				throw new BadDatabaseSchemaException();
			}
		}
		catch (BadSqlGrammarException ex)
		{
			throw new DatabaseNotInitializedException();
		}

		return true;
	}
}
