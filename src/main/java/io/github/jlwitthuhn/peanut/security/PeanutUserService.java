// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.security;

import io.github.jlwitthuhn.peanut.db.UserDAO;
import io.github.jlwitthuhn.peanut.model.db.UserRow;
import io.github.jlwitthuhn.peanut.model.spring.PeanutUserDetails;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import org.springframework.security.provisioning.UserDetailsManager;
import org.springframework.stereotype.Component;

@Component
public class PeanutUserService implements UserDetailsManager
{
	private final UserDAO userDAO;

	public PeanutUserService(UserDAO userDAO)
	{
		this.userDAO = userDAO;
	}

	@Override
	public void createUser(UserDetails user)
	{
		userDAO.createRow(user.getUsername(), user.getPassword());
	}

	@Override
	public void updateUser(UserDetails user)
	{
		throw new UnsupportedOperationException("TODO: updateUser");
	}

	@Override
	public void deleteUser(String username)
	{
		throw new UnsupportedOperationException("TODO: deleteUser");
	}

	@Override
	public void changePassword(String oldPassword, String newPassword)
	{
		throw new UnsupportedOperationException("TODO: changePassword");
	}

	@Override
	public boolean userExists(String username)
	{
		return false;
	}

	@Override
	public UserDetails loadUserByUsername(String username) throws UsernameNotFoundException
	{
		UserRow row = userDAO.getByName(username);
		if (row == null)
		{
			throw new UsernameNotFoundException(username);
		}
		return new PeanutUserDetails(row);
	}
}
