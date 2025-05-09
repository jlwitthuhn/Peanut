// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.controller;

import io.github.jlwitthuhn.peanut.db.MetaDAO;
import io.github.jlwitthuhn.peanut.db.setup.DatabaseInitializer;
import io.github.jlwitthuhn.peanut.model.form.SetupForm;
import io.github.jlwitthuhn.peanut.model.spring.PeanutUserDetails;
import io.github.jlwitthuhn.peanut.security.PeanutUserService;
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

import java.util.ArrayList;
import java.util.Map;

@Controller
@RequestMapping("/setup")
@RequiredArgsConstructor
public class SetupController
{
	private final DatabaseInitializer databaseInitializer;
	private final MetaDAO metaDAO;
	private final PasswordEncoder passwordEncoder;
	private final PeanutUserService peanutUserService;

	@GetMapping("")
	public ModelAndView index(Map<String, Object> model)
	{
		if (metaDAO.doesTableExist("config_int"))
		{
			return ViewShortcuts.simpleMessage("Setup complete", "The database has already been initialized. No further setup action is needed.");
		}
		return new ModelAndView("setup.html", model);
	}

	@PostMapping("")
	public ModelAndView post(@ModelAttribute SetupForm form)
	{
		if (metaDAO.doesTableExist("config_int"))
		{
			return ViewShortcuts.simpleMessage("Error", "The database was already set up. Its existing state has not been modified.", HttpStatus.CONFLICT);
		}

		if (form == null || !form.isValid())
		{
			return ViewShortcuts.simpleMessage("Error", "Form is invalid.", HttpStatus.BAD_REQUEST);
		}

		if (!databaseInitializer.doInit())
		{
			return ViewShortcuts.simpleMessage("Error", "Encountered error while initializing database.", HttpStatus.INTERNAL_SERVER_ERROR);
		}

		// Create admin user
		String hashedPassword = passwordEncoder.encode(form.getAdminPass());
		ArrayList<GrantedAuthority> authorities = new ArrayList<GrantedAuthority>();
		authorities.add(new SimpleGrantedAuthority("ROLE_TURBO_ADMIN"));
		authorities.add(new SimpleGrantedAuthority("ROLE_ADMIN"));
		authorities.add(new SimpleGrantedAuthority("ROLE_USER"));
		PeanutUserDetails userDetails = new PeanutUserDetails(form.getAdminName(), form.getEmail(), hashedPassword, authorities);
		peanutUserService.createUser(userDetails);

		return ViewShortcuts.simpleMessage("Success", "Database initialized successfully.");
	}
}
