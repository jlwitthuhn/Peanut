// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.controller;

import io.github.jlwitthuhn.peanut.err.DatabaseNotInitializedException;
import io.github.jlwitthuhn.peanut.util.ViewShortcuts;
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
				return ViewShortcuts.simpleMessage("SQL Error", "Failed to execute SQL query.", HttpStatus.INTERNAL_SERVER_ERROR);
			}
			case CannotGetJdbcConnectionException cannotGetJdbcConnectionException ->
			{
				return ViewShortcuts.simpleMessage("DB Error", "Unable to connect to the database.", HttpStatus.INTERNAL_SERVER_ERROR);
			}
			case DatabaseNotInitializedException databaseNotInitializedException ->
			{
				return ViewShortcuts.simpleMessage("DB Error", "Database has not been initialized.", HttpStatus.INTERNAL_SERVER_ERROR);
			}
			default ->
			{
				logger.error("Caught unknown exception in ExceptionResolver", ex);
				return ViewShortcuts.simpleMessage("Error", "Unknown error occurred.", HttpStatus.INTERNAL_SERVER_ERROR);
			}
		}
	}
}
