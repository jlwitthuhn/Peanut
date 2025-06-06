// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.interceptor;

import io.github.jlwitthuhn.peanut.model.spring.PeanutUserDetails;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.security.authentication.AnonymousAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.stereotype.Component;
import org.springframework.web.servlet.HandlerInterceptor;
import org.springframework.web.servlet.ModelAndView;

import java.util.HashMap;

@Component
public class LoggedInUserInterceptor implements HandlerInterceptor
{
	private static final Logger logger = LoggerFactory.getLogger(LoggedInUserInterceptor.class);

	@Override
	public void postHandle(HttpServletRequest request, HttpServletResponse response, Object handler, ModelAndView modelAndView) throws Exception
	{
		if (modelAndView == null)
		{
			return;
		}

		HashMap<String, Object> userInfo = new HashMap<>();
		userInfo.put("loggedIn", Boolean.FALSE);
		userInfo.put("admin", Boolean.FALSE);
		userInfo.put("name", "Anonymous");
		modelAndView.getModel().put("user", userInfo);

		Authentication auth = SecurityContextHolder.getContext().getAuthentication();
		if (auth == null || auth instanceof AnonymousAuthenticationToken)
		{
			return;
		}

		if (!(auth.getPrincipal() instanceof PeanutUserDetails userDetails))
		{
			return;
		}

		boolean isAdmin = false;
		for (GrantedAuthority authority : auth.getAuthorities())
		{
			if (authority.getAuthority().equals("ADMIN"))
			{
				isAdmin = true;
				break;
			}
		}

		if (userDetails.getId().isEmpty())
		{
			logger.error("Found user with no id, name: {}", userDetails.getUsername());
			return;
		}

		userInfo.put("loggedIn", Boolean.TRUE);
		userInfo.put("admin", isAdmin);
		userInfo.put("name", userDetails.getUsername());
		userInfo.put("id", userDetails.getId().get());
	}
}
