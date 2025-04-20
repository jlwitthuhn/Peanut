// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.controller;

import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;

@Controller
public class RootController
{
	@GetMapping("/")
	public String index()
	{
		return "index.html";
	}

	@GetMapping("/design")
	public String design()
	{
		return "design.html";
	}

	@GetMapping("/login")
	public String login()
	{
		return "login.html";
	}
}
