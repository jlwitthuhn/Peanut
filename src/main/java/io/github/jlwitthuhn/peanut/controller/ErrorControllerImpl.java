// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.controller;

import io.github.jlwitthuhn.peanut.util.ViewShortcuts;
import jakarta.servlet.RequestDispatcher;
import jakarta.servlet.http.HttpServletRequest;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
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
	private static final Logger logger = LoggerFactory.getLogger(ErrorControllerImpl.class);

	@GetMapping
	public ModelAndView error(HttpServletRequest request)
	{
		Object status = request.getAttribute(RequestDispatcher.ERROR_STATUS_CODE);
		if (status instanceof Integer statusInt)
		{
			final String title = String.format("Error - HTTP %d", statusInt);
			return ViewShortcuts.simpleMessage(title, "An error occurred", HttpStatus.valueOf(statusInt));
		}

		logger.error("Caught non-http error in ErrorControllerImpl");
		Object message = request.getAttribute(RequestDispatcher.ERROR_MESSAGE);
		if (message instanceof String messageStr)
		{
			logger.error(messageStr);
		}
		Object exception = request.getAttribute(RequestDispatcher.ERROR_EXCEPTION);
		if (exception instanceof Exception ex)
		{
			logger.error(ex.getMessage(), ex);
		}

		return ViewShortcuts.simpleMessage("Error", "An error occurred, please try again.", HttpStatus.INTERNAL_SERVER_ERROR);
	}
}
