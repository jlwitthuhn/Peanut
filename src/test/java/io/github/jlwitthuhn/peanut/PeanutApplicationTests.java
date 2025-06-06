// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut;

import io.github.jlwitthuhn.peanut.controller.AdminController;
import io.github.jlwitthuhn.peanut.controller.LoginController;
import io.github.jlwitthuhn.peanut.controller.LogoutController;
import io.github.jlwitthuhn.peanut.controller.RegisterController;
import io.github.jlwitthuhn.peanut.controller.RootController;
import io.github.jlwitthuhn.peanut.controller.SetupController;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;

import static org.assertj.core.api.Assertions.assertThat;

@SpringBootTest
class PeanutApplicationTests
{
	@Autowired
	private AdminController adminController;

	@Autowired
	private LoginController loginController;

	@Autowired
	private LogoutController logoutController;

	@Autowired
	private RegisterController registerController;

	@Autowired
	private RootController rootController;

	@Autowired
	private SetupController setupController;

	@Test
	void controllersLoad()
	{
		assertThat(adminController).isNotNull();
		assertThat(loginController).isNotNull();
		assertThat(logoutController).isNotNull();
		assertThat(registerController).isNotNull();
		assertThat(rootController).isNotNull();
		assertThat(setupController).isNotNull();
	}
}
