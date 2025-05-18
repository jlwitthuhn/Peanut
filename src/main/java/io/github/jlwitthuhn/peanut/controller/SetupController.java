// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.controller;

import io.github.jlwitthuhn.peanut.db.MetaDAO;
import io.github.jlwitthuhn.peanut.err.DBCreationDependencyNotSatisfiedException;
import io.github.jlwitthuhn.peanut.model.form.SetupForm;
import io.github.jlwitthuhn.peanut.service.SetupService;
import io.github.jlwitthuhn.peanut.util.ViewShortcuts;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.ModelAttribute;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.servlet.ModelAndView;

import java.util.Map;

@Controller
@RequestMapping("/setup")
@RequiredArgsConstructor
public class SetupController
{
	private final SetupService setupService;

	private final MetaDAO metaDAO;

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
		if (form == null || !form.isValid())
		{
			return ViewShortcuts.simpleMessage("Error", "Form is invalid.", HttpStatus.BAD_REQUEST);
		}

		try
		{
			setupService.initializeDatabase(form.getAdminName(), form.getAdminPass(), form.getEmail());
		}
		catch (DBCreationDependencyNotSatisfiedException ex)
		{
			return ViewShortcuts.simpleMessage("Init Error", "Failed to initialize database: " + ex.getMessage(), HttpStatus.INTERNAL_SERVER_ERROR);
		}

		return ViewShortcuts.simpleMessage("Success", "Database initialized successfully.");
	}
}
