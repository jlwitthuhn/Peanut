// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.controller;

import io.github.jlwitthuhn.peanut.err.UserDetailsConflictException;
import io.github.jlwitthuhn.peanut.model.form.RegisterForm;
import io.github.jlwitthuhn.peanut.model.spring.PeanutUserDetails;
import io.github.jlwitthuhn.peanut.security.PeanutUserService;
import io.github.jlwitthuhn.peanut.util.AuthorizationUtil;
import io.github.jlwitthuhn.peanut.util.ViewShortcuts;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.ModelAttribute;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.servlet.ModelAndView;
import org.springframework.web.servlet.view.RedirectView;

import java.util.ArrayList;

@Controller
@RequestMapping("/register")
@RequiredArgsConstructor
public class RegisterController
{
	private final PasswordEncoder passwordEncoder;
	private final PeanutUserService userService;

	@GetMapping("")
	public ModelAndView getIndex()
	{
		if (AuthorizationUtil.currentUserIsLoggedIn())
		{
			RedirectView view = new RedirectView("/");
			view.setStatusCode(HttpStatus.SEE_OTHER);
			return new ModelAndView(view);
		}
		return new ModelAndView("register.html");
	}

	@PostMapping("")
	public ModelAndView postIndex(@ModelAttribute RegisterForm form)
	{
		if (AuthorizationUtil.currentUserIsLoggedIn())
		{
			RedirectView view = new RedirectView("/");
			view.setStatusCode(HttpStatus.SEE_OTHER);
			return new ModelAndView(view);
		}

		if (form == null || !form.isValid())
		{
			return ViewShortcuts.simpleMessage("Error", "Form is invalid.", HttpStatus.BAD_REQUEST);
		}

		ArrayList<GrantedAuthority> authorities = new ArrayList<>();
		authorities.add(new SimpleGrantedAuthority("USER"));
		PeanutUserDetails newUser = new PeanutUserDetails(form.getUsername(), form.getEmail(), passwordEncoder.encode(form.getPassword()), authorities);
		try
		{
			userService.createUser(newUser);
		}
		catch (UserDetailsConflictException ex)
		{
			return ViewShortcuts.simpleMessage("Failed to create user", ex.getMessage(), HttpStatus.CONFLICT);
		}

		return ViewShortcuts.simpleMessage("Success", "Your account has been created. You may now log in.");
	}
}
