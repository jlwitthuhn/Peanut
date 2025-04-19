// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.controller;

import io.github.jlwitthuhn.peanut.model.form.SetupForm;
import io.github.jlwitthuhn.peanut.util.ViewShortcuts;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.ModelAttribute;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.servlet.ModelAndView;

import java.util.HashMap;

@Controller
@RequestMapping("/setup")
public class SetupController
{
	@GetMapping("")
	public String index()
	{
		return "setup.html";
	}

	@PostMapping("")
	public ModelAndView post(@ModelAttribute SetupForm form)
	{
		if (form == null || !form.isValid())
		{
			return ViewShortcuts.simpleMessage("Error", "Form is invalid.", HttpStatus.BAD_REQUEST);
		}

		// TODO: Create tables

		return ViewShortcuts.simpleMessage("Success", "Database initialized successfully.");
	}
}
