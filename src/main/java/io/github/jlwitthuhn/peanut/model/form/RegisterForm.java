// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.model.form;

import lombok.Getter;
import lombok.Setter;

public class RegisterForm
{
	@Getter
	@Setter
	private String username;

	@Getter
	@Setter
	private String email;

	@Getter
	@Setter
	private String password;

	@Getter
	@Setter
	private String password2;

	public boolean isValid()
	{
		final boolean notNull = username != null && email != null && password != null && password2 != null;
		final boolean goodLength = notNull && username.length() > 1 && email.length() > 1 && password.length() > 1 && password2.length() > 1;
		final boolean passwordsMatch = password != null && password.equals(password2);
		return notNull && goodLength && passwordsMatch;
	}
}
