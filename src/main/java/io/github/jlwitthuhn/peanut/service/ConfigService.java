// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.service;

import io.github.jlwitthuhn.peanut.db.ConfigDAO;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;

@Service
@RequiredArgsConstructor
public class ConfigService
{
	private final ConfigDAO configDAO;

	public Long getLong(String name)
	{
		return configDAO.selectLongByName(name);
	}

	public String getString(String name)
	{
		return configDAO.selectStringByName(name);
	}

	public void setString(String name, String value)
	{
		// TODO: There should be some access control based on which value is being set here
		configDAO.upsertStringByName(name, value);
	}
}
