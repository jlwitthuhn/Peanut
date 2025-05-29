// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.cfg;

import io.github.jlwitthuhn.peanut.interceptor.DatabaseInitInterceptor;
import io.github.jlwitthuhn.peanut.interceptor.LoggedInUserInterceptor;
import io.github.jlwitthuhn.peanut.interceptor.RequestTimerEndInterceptor;
import io.github.jlwitthuhn.peanut.interceptor.RequestTimerStartInterceptor;
import lombok.RequiredArgsConstructor;
import org.springframework.context.annotation.Configuration;
import org.springframework.lang.NonNull;
import org.springframework.web.servlet.config.annotation.InterceptorRegistry;
import org.springframework.web.servlet.config.annotation.WebMvcConfigurer;

@Configuration
@RequiredArgsConstructor
public class PeanutConfiguration implements WebMvcConfigurer
{
	private final DatabaseInitInterceptor databaseInitInterceptor;
	private final LoggedInUserInterceptor loggedInUserInterceptor;
	private final RequestTimerStartInterceptor requestTimerStartInterceptor;
	private final RequestTimerEndInterceptor requestTimerEndInterceptor;

	@Override
	public void addInterceptors(@NonNull InterceptorRegistry registry)
	{
		// Must be first
		registry.addInterceptor(requestTimerStartInterceptor);

		// Order here matters and might carry hidden dependencies, do not change unless there is a reason
		registry.addInterceptor(databaseInitInterceptor).excludePathPatterns("/design*", "/setup*");
		registry.addInterceptor(loggedInUserInterceptor);

		// Must be last
		registry.addInterceptor(requestTimerEndInterceptor);
	}
}
