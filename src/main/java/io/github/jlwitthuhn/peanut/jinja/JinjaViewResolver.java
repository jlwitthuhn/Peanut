// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.jinja;

import com.hubspot.jinjava.Jinjava;
import org.springframework.lang.NonNull;
import org.springframework.stereotype.Component;
import org.springframework.web.servlet.view.AbstractTemplateViewResolver;
import org.springframework.web.servlet.view.AbstractUrlBasedView;

@Component
public class JinjaViewResolver extends AbstractTemplateViewResolver
{
	private final Jinjava jinjava;

	public JinjaViewResolver(Jinjava jinjava)
	{
		this.jinjava = jinjava;
		setViewClass(JinjaView.class);
	}

	@Override
	@NonNull
	protected Class<?> requiredViewClass()
	{
		return JinjaView.class;
	}

	@Override
	protected AbstractUrlBasedView buildView(String viewName) throws Exception
	{
		JinjaView view = (JinjaView) super.buildView(viewName);
		view.setJinjava(jinjava);
		return view;
	}
}
