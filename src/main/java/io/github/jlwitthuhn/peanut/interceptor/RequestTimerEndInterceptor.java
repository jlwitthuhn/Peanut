// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.interceptor;

import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Component;
import org.springframework.web.servlet.HandlerInterceptor;
import org.springframework.web.servlet.ModelAndView;

@Component
public class RequestTimerEndInterceptor implements HandlerInterceptor
{
	private static final Logger logger = LoggerFactory.getLogger(RequestTimerEndInterceptor.class);

	@Override
	public void postHandle(HttpServletRequest request, HttpServletResponse response, Object handler, ModelAndView modelAndView) throws Exception
	{
		Object requestStartMillisObj = request.getAttribute("requestStartMillis");
		if (!(requestStartMillisObj instanceof Long requestStartMillis))
		{
			logger.error("Attribute 'requestStartMillis' does not exist or is not a Long");
			return;
		}
		Long requestDuration = System.currentTimeMillis() - requestStartMillis;
		modelAndView.getModel().put("requestDurationMs", requestDuration);
	}
}
