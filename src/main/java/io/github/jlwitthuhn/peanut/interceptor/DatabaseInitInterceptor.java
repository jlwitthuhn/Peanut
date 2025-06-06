// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.interceptor;

import io.github.jlwitthuhn.peanut.cfg.ConfigKeyNames;
import io.github.jlwitthuhn.peanut.err.DBBadSchemaException;
import io.github.jlwitthuhn.peanut.err.DBNotInitializedException;
import io.github.jlwitthuhn.peanut.service.ConfigService;
import io.github.jlwitthuhn.peanut.service.DatabaseService;
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
	private final ConfigService configService;
	private final DatabaseService dbService;

	@Override
	public boolean preHandle(@NonNull HttpServletRequest request, @NonNull HttpServletResponse response, @NonNull Object handler)
	{
		boolean configExists = dbService.doesTableExist("config_int");
		if (!configExists)
		{
			throw new DBNotInitializedException();
		}
		Long version = configService.getLong(ConfigKeyNames.SCHEMA_VERSION_INT);
		if (version == null || version != 1)
		{
			throw new DBBadSchemaException();
		}

		return true;
	}
}
