// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.controller;

import io.github.jlwitthuhn.peanut.err.DBNotInitializedException;
import io.github.jlwitthuhn.peanut.interceptor.LoggedInUserInterceptor;
import io.github.jlwitthuhn.peanut.util.ViewShortcuts;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import lombok.RequiredArgsConstructor;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.HttpStatus;
import org.springframework.jdbc.BadSqlGrammarException;
import org.springframework.jdbc.CannotGetJdbcConnectionException;
import org.springframework.security.authorization.AuthorizationDeniedException;
import org.springframework.stereotype.Component;
import org.springframework.web.servlet.ModelAndView;
import org.springframework.web.servlet.handler.AbstractHandlerExceptionResolver;

@Component
@RequiredArgsConstructor
public class ExceptionResolver extends AbstractHandlerExceptionResolver
{
	private final LoggedInUserInterceptor loggedInUserInterceptor;

	private static final Logger logger = LoggerFactory.getLogger(ExceptionResolver.class);

	@Override
	protected ModelAndView doResolveException(HttpServletRequest request, HttpServletResponse response, Object handler, Exception ex)
	{
		ModelAndView result = null;
		switch (ex) {
			case BadSqlGrammarException badSqlGrammarException ->
			{
				logger.error("Caught BadSqlGrammarException in ExceptionResolver", ex);
				result = ViewShortcuts.simpleMessage("SQL Error", "Failed to execute SQL query.", HttpStatus.INTERNAL_SERVER_ERROR);
			}
			case CannotGetJdbcConnectionException cannotGetJdbcConnectionException ->
			{
				result = ViewShortcuts.simpleMessage("DB Error", "Unable to connect to the database.", HttpStatus.INTERNAL_SERVER_ERROR);
			}
			case DBNotInitializedException dbNotInitializedException ->
			{
				result = ViewShortcuts.simpleMessage("DB Error", "Database has not been initialized.", HttpStatus.INTERNAL_SERVER_ERROR);
			}
			case AuthorizationDeniedException authorizationDeniedException ->
			{
				result = ViewShortcuts.simpleMessage("Error - HTTP 403", "Forbidden", HttpStatus.FORBIDDEN);
			}
			default ->
			{
				logger.error("Caught unknown exception in ExceptionResolver", ex);
				result = ViewShortcuts.simpleMessage("Error", "Unknown error occurred.", HttpStatus.INTERNAL_SERVER_ERROR);
			}
		}

		try
		{
			loggedInUserInterceptor.postHandle(request, response, handler, result);
		}
		catch (Exception exInner)
		{
			logger.error("Caught unknown exception while fetching login info", exInner);
		}
		return result;
	}
}
