// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.cfg;

import org.springframework.boot.web.embedded.tomcat.TomcatServletWebServerFactory;
import org.springframework.boot.web.server.WebServerFactoryCustomizer;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class TomcatConfiguration
{

	@Bean
	public WebServerFactoryCustomizer<TomcatServletWebServerFactory> customTomcat() {
		return factory ->
			factory.addConnectorCustomizers(
				connector ->
					connector.setEncodedSolidusHandling("decode")
			);
	}
}
