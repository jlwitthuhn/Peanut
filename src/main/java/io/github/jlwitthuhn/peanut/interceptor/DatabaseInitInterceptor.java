// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.interceptor;

import io.github.jlwitthuhn.peanut.cfg.ConfigKeyNames;
import io.github.jlwitthuhn.peanut.db.ConfigDAO;
import io.github.jlwitthuhn.peanut.db.MetaDAO;
import io.github.jlwitthuhn.peanut.err.DBBadSchemaException;
import io.github.jlwitthuhn.peanut.err.DBNotInitializedException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import lombok.RequiredArgsConstructor;
import org.springframework.lang.NonNull;
import org.springframework.stereotype.Component;
import org.springframework.web.servlet.HandlerInterceptor;

@Component
@RequiredArgsConstructor
public class DatabaseInitInterceptor implements HandlerInterceptor
{
	final ConfigDAO configDAO;
	final MetaDAO metaDAO;

	@Override
	public boolean preHandle(@NonNull HttpServletRequest request, @NonNull HttpServletResponse response, @NonNull Object handler)
	{
		boolean configExists = metaDAO.doesTableExist("config_int");
		if (!configExists)
		{
			throw new DBNotInitializedException();
		}
		Long version = configDAO.getLong(ConfigKeyNames.SCHEMA_VERSION_INT);
		if (version == null || version != 1)
		{
			throw new DBBadSchemaException();
		}

		return true;
	}
}
