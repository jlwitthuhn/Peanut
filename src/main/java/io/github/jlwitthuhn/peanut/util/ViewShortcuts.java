// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.util;

import org.springframework.http.HttpStatus;
import org.springframework.web.servlet.ModelAndView;
import org.springframework.web.servlet.view.RedirectView;

import java.util.HashMap;

public class ViewShortcuts
{
	public static ModelAndView simpleMessage(String header, String message, HttpStatus status)
	{
		var model = new HashMap<String, String>();
		model.put("header", header);
		model.put("message", message);
		return new ModelAndView("simple_message.html", model, status);
	}

	public static ModelAndView simpleMessage(String header, String message)
	{
		return simpleMessage(header, message, HttpStatus.OK);
	}

	public static ModelAndView simpleRedirect(String redirectPath)
	{
		RedirectView view = new RedirectView(redirectPath);
		view.setStatusCode(HttpStatus.SEE_OTHER);
		view.setExposeModelAttributes(false);
		return new ModelAndView(view);
	}
}
