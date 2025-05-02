// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.controller;

import io.github.jlwitthuhn.peanut.util.ViewShortcuts;
import jakarta.servlet.RequestDispatcher;
import jakarta.servlet.http.HttpServletRequest;
import org.springframework.boot.web.servlet.error.ErrorController;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.servlet.ModelAndView;

@Controller
@RequestMapping("/error")
public class ErrorControllerImpl implements ErrorController
{
	@GetMapping
	public ModelAndView error(HttpServletRequest request)
	{
		Object status = request.getAttribute(RequestDispatcher.ERROR_STATUS_CODE);
		if (!(status instanceof Integer statusInt))
		{
			return ViewShortcuts.simpleMessage("Error", "An error occurred", HttpStatus.INTERNAL_SERVER_ERROR);
		}
		final String title = String.format("Error - HTTP %d", statusInt);
		return ViewShortcuts.simpleMessage(title, "An error occurred", HttpStatus.valueOf(statusInt));
	}
}
