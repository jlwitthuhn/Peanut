// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.cfg;

import io.github.jlwitthuhn.peanut.interceptor.DatabaseInitInterceptor;
import org.springframework.context.annotation.Configuration;
import org.springframework.lang.NonNull;
import org.springframework.web.servlet.config.annotation.InterceptorRegistry;
import org.springframework.web.servlet.config.annotation.WebMvcConfigurer;

@Configuration
public class PeanutConfigurer implements WebMvcConfigurer
{
	private final DatabaseInitInterceptor databaseInitInterceptor;

	public PeanutConfigurer(DatabaseInitInterceptor databaseInitInterceptor)
	{
		this.databaseInitInterceptor = databaseInitInterceptor;
	}

	@Override
	public void addInterceptors(@NonNull InterceptorRegistry registry)
	{
		registry.addInterceptor(databaseInitInterceptor).excludePathPatterns("/design*", "/setup*");
	}
}
