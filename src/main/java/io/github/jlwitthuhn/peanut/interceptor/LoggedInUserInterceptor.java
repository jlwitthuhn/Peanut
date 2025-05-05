// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.interceptor;

import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.springframework.security.authentication.AnonymousAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.stereotype.Component;
import org.springframework.web.servlet.HandlerInterceptor;
import org.springframework.web.servlet.ModelAndView;

import java.util.HashMap;

@Component
public class LoggedInUserInterceptor implements HandlerInterceptor
{
	@Override
	public void postHandle(HttpServletRequest request, HttpServletResponse response, Object handler, ModelAndView modelAndView) throws Exception
	{
		if (modelAndView == null)
		{
			return;
		}

		HashMap<String, Object> userInfo = new HashMap<>();
		userInfo.put("logged_in", Boolean.FALSE);
		userInfo.put("name", "Anonymous");
		modelAndView.getModel().put("user", userInfo);

		Authentication auth = SecurityContextHolder.getContext().getAuthentication();
		if (auth == null || auth instanceof AnonymousAuthenticationToken)
		{
			return;
		}

		userInfo.put("logged_in", Boolean.TRUE);
		userInfo.put("name", auth.getName());
	}
}
