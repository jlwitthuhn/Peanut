// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.jinja;

import com.hubspot.jinjava.Jinjava;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class JinjaConfiguration
{
	@Bean
	Jinjava jinjava(JinjaResourceLocator jinjaResourceLocator)
	{
		Jinjava jinjava = new Jinjava();
		jinjava.setResourceLocator(jinjaResourceLocator);
		return jinjava;
	}
}
