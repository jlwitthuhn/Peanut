// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.controller;

import io.github.jlwitthuhn.peanut.err.DatabaseNotInitializedException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.HttpStatus;
import org.springframework.jdbc.BadSqlGrammarException;
import org.springframework.jdbc.CannotGetJdbcConnectionException;
import org.springframework.stereotype.Component;
import org.springframework.web.servlet.ModelAndView;
import org.springframework.web.servlet.handler.AbstractHandlerExceptionResolver;

import java.util.HashMap;

@Component
public class ExceptionResolver extends AbstractHandlerExceptionResolver
{
	private static final Logger logger = LoggerFactory.getLogger(ExceptionResolver.class);

	@Override
	protected ModelAndView doResolveException(HttpServletRequest request, HttpServletResponse response, Object handler, Exception ex)
	{
		switch (ex) {
			case BadSqlGrammarException badSqlGrammarException ->
			{
				logger.error("Caught BadSqlGrammarException in ExceptionResolver", ex);
				HashMap<String, String> model = new HashMap<>();
				model.put("message", "SQL error");
				return new ModelAndView("error.html", model, HttpStatus.INTERNAL_SERVER_ERROR);
			}
			case CannotGetJdbcConnectionException cannotGetJdbcConnectionException ->
			{
				HashMap<String, String> model = new HashMap<>();
				model.put("message", "Unable to connect to the database");
				return new ModelAndView("error.html", model, HttpStatus.INTERNAL_SERVER_ERROR);
			}
			case DatabaseNotInitializedException databaseNotInitializedException ->
			{
				HashMap<String, String> model = new HashMap<>();
				model.put("message", "Database has not been initialized");
				return new ModelAndView("error.html", model, HttpStatus.INTERNAL_SERVER_ERROR);
			}
			default ->
			{
				logger.error("Caught unknown exception in ExceptionResolver", ex);
				HashMap<String, String> model = new HashMap<>();
				model.put("message", "Unknown error occurred");
				return new ModelAndView("error.html", model, HttpStatus.INTERNAL_SERVER_ERROR);
			}
		}
	}
}
