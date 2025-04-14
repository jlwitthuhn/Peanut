// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.jinja;

import com.hubspot.jinjava.interpret.JinjavaInterpreter;
import com.hubspot.jinjava.loader.ResourceLocator;
import org.springframework.context.EnvironmentAware;
import org.springframework.core.env.Environment;
import org.springframework.core.io.DefaultResourceLoader;
import org.springframework.core.io.Resource;
import org.springframework.core.io.ResourceLoader;
import org.springframework.lang.NonNull;
import org.springframework.stereotype.Component;

import java.io.File;
import java.io.IOException;
import java.nio.charset.Charset;
import java.nio.file.Files;

@Component
public class JinjaResourceLocator implements EnvironmentAware, ResourceLocator
{
	private final ResourceLoader resourceLoader = new DefaultResourceLoader();
	private String pathPrefix = "classpath:/templates/";
	private String pathSuffix = ".html";

	@Override
	public String getString(String fullName, Charset encoding, JinjavaInterpreter interpreter) throws IOException
	{
		Resource resource = resourceLoader.getResource(pathPrefix + fullName + pathSuffix);
		if (!resource.exists())
		{
			throw new IOException(fullName + " does not exist");
		}
		File file = resource.getFile();
		return Files.readString(file.toPath(), encoding);
	}

	@Override
	public void setEnvironment(@NonNull Environment environment)
	{
		String maybePrefix = environment.getProperty("jinja.path.prefix");
		if (maybePrefix != null)
		{
			pathPrefix = maybePrefix;
		}
		String maybeSuffix = environment.getProperty("jinja.path.suffix");
		if (maybeSuffix != null) {
			pathSuffix = maybeSuffix;
		}
	}
}
