package io.github.jlwitthuhn.peanut.controller;

import io.github.jlwitthuhn.peanut.util.ViewShortcuts;
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
	public ModelAndView error()
	{
		return ViewShortcuts.simpleMessage("Error", "Something happened.", HttpStatus.INTERNAL_SERVER_ERROR);
	}
}
