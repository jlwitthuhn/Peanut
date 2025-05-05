// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.controller;

import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.servlet.view.RedirectView;

@Controller
@RequestMapping("/logout")
public class LogoutController
{
	@GetMapping("/success")
	public RedirectView loginSuccess()
	{
		RedirectView view = new RedirectView("/");
		view.setStatusCode(HttpStatus.SEE_OTHER);
		return view;
	}
}
