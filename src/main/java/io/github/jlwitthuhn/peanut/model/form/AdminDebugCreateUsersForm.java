// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.model.form;

import lombok.Getter;
import lombok.Setter;

import java.util.ArrayList;
import java.util.List;

public class AdminDebugCreateUsersForm
{
	@Getter
	@Setter
	private String count;

	@Getter
	@Setter
	private String prefix;

	@Getter
	@Setter
	private String password;

	public Integer getCountInt()
	{
		return Integer.parseInt(count);
	}

	public List<String> getValidationErrors()
	{
		List<String> result = new ArrayList<String>();

		try
		{
			Integer.parseInt(count);
		}
		catch (NumberFormatException ex)
		{
			result.add("Count must be an integer");
		}

		return result;
	}
}
