// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.model.form;

import lombok.Getter;
import lombok.Setter;

public class SetupForm
{
	@Getter
	@Setter
	private String adminName;

	@Getter
	@Setter
	private String email;

	@Getter
	@Setter
	private String adminPass;

	@Getter
	@Setter
	private String adminPass2;

	public boolean isValid()
	{
		final boolean notNull = adminName != null && email != null && adminPass != null && adminPass2 != null;
		final boolean goodLength = notNull && adminName.length() > 1 && email.length() > 1 && adminPass.length() > 1 && adminPass2.length() > 1;
		final boolean passwordsMatch = adminPass != null && adminPass.equals(adminPass2);
		return notNull && goodLength && passwordsMatch;
	}
}
