// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.service;

import io.github.jlwitthuhn.peanut.model.spring.PeanutUserDetails;
import io.github.jlwitthuhn.peanut.security.PeanutUserService;
import lombok.RequiredArgsConstructor;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

@Service
@RequiredArgsConstructor
public class AdminService
{
	private final PasswordEncoder passwordEncoder;
	private final PeanutUserService userService;

	@Transactional
	public void createDebugUsers(int count, String prefix, String plainPassword)
	{
		for (int i = 0; i < count; i++)
		{
			String suffix = String.format("%04d", i);
			String accountName = prefix + suffix;
			String email = "test"; //accountName + "@peanut";
			String hashedPassword = passwordEncoder.encode(plainPassword);
			PeanutUserDetails thisUser = new PeanutUserDetails(accountName, email, hashedPassword);
			userService.createUser(thisUser);
		}
	}
}
