// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.controller;

import io.github.jlwitthuhn.peanut.util.ViewShortcuts;
import org.springframework.http.HttpStatus;
import org.springframework.security.authentication.AnonymousAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.servlet.ModelAndView;

@Controller
@RequestMapping("/login")
public class LoginController
{
	@GetMapping("")
	public ModelAndView login()
	{
		Authentication auth = SecurityContextHolder.getContext().getAuthentication();
		if (auth == null || auth instanceof AnonymousAuthenticationToken)
		{
			return new ModelAndView("login.html");
		}

		// Already logged in
		return ViewShortcuts.simpleRedirect("/");
	}

	@GetMapping("/failure")
	public ModelAndView loginFailure()
	{
		return ViewShortcuts.simpleMessage("Login Failure", "Username and password do not match.", HttpStatus.FORBIDDEN);
	}

	@GetMapping("/success")
	public ModelAndView loginSuccess()
	{
		return ViewShortcuts.simpleRedirect("/");
	}
}
