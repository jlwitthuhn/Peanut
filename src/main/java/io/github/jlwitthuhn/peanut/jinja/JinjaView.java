// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.jinja;

import com.google.common.base.Charsets;
import com.hubspot.jinjava.Jinjava;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import lombok.Setter;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.lang.NonNull;
import org.springframework.security.web.csrf.CsrfToken;
import org.springframework.web.servlet.view.AbstractTemplateView;

import java.io.PrintWriter;
import java.util.Map;

public class JinjaView extends AbstractTemplateView
{
	@Setter
	private Jinjava jinjava;

	private final Logger logger = LoggerFactory.getLogger(JinjaView.class);

	@Override
	protected void renderMergedTemplateModel(@NonNull Map<String, Object> model, @NonNull HttpServletRequest request, @NonNull HttpServletResponse response) throws Exception
	{
		response.setContentType("text/html");
		response.setCharacterEncoding("UTF-8");

		if (request.getAttribute("_csrf") instanceof CsrfToken csrfToken)
		{
			model.put("_csrfParam", csrfToken.getParameterName());
			model.put("_csrfToken", csrfToken.getToken());
		}

		PrintWriter responseWriter = response.getWriter();
		try
		{
			String template = jinjava.getResourceLocator().getString(getUrl(), Charsets.UTF_8, null);
			String rendered = jinjava.render(template, model);
			responseWriter.write(rendered);
		}
		catch (Exception e)
		{
			logger.error("JinjaResourceLocator::renderMergedTemplateModel caught exception: {}", String.valueOf(e));
		}
	}
}
