// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.model.form;

import lombok.Getter;
import lombok.Setter;

public class AdminGeneralWelcomeMessageForm
{
	@Getter
	@Setter
	private Boolean confirm;

	@Getter
	@Setter
	private String message;
}
