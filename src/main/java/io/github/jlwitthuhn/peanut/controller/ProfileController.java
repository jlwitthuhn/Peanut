// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.controller;

import io.github.jlwitthuhn.peanut.db.UserDAO;
import io.github.jlwitthuhn.peanut.model.db.UserRow;
import io.github.jlwitthuhn.peanut.util.TimeUtil;
import io.github.jlwitthuhn.peanut.util.ViewShortcuts;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.servlet.ModelAndView;

import java.util.HashMap;
import java.util.Map;

@Controller
@RequestMapping("/profile")
@RequiredArgsConstructor
public class ProfileController
{
	private final UserDAO userDAO;

	@GetMapping("/view/{id}")
	public ModelAndView viewById(@PathVariable long id, Map<String, Object> model)
	{
		UserRow userRow = userDAO.selectRowById(id);

		if(userRow == null)
		{
			return ViewShortcuts.simpleMessage("404 - User not found", "No user exists with id " + id, HttpStatus.NOT_FOUND);
		}

		Map<String, String> profileData = new HashMap<>();
		profileData.put("displayName", userRow.getDisplayName());
		profileData.put("id", Long.toString(userRow.getId()));
		profileData.put("created", TimeUtil.formatOffsetDateTime(userRow.getCreatedTimestamp()));

		model.put("profile", profileData);
		return new ModelAndView("profile.html", model);
	}
}
