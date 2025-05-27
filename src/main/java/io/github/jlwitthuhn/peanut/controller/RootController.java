// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.controller;

import io.github.jlwitthuhn.peanut.cfg.ConfigKeyNames;
import io.github.jlwitthuhn.peanut.db.ConfigDAO;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.servlet.ModelAndView;

import java.util.Map;

@Controller
@RequiredArgsConstructor
public class RootController
{
	private final ConfigDAO configDAO;

	@GetMapping("/")
	public ModelAndView index(Map<String, Object> model)
	{
		String welcomeMessage = configDAO.getString(ConfigKeyNames.WELCOME_MESSAGE_STR);
		model.put("welcomeMessage", welcomeMessage);
		return new ModelAndView("index.html", model);
	}

	@GetMapping("/design")
	public String design()
	{
		return "design.html";
	}
}
