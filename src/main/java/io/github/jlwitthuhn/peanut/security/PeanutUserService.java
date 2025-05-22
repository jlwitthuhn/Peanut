// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.security;

import io.github.jlwitthuhn.peanut.db.GroupDAO;
import io.github.jlwitthuhn.peanut.db.GroupMembershipDAO;
import io.github.jlwitthuhn.peanut.db.MultiTableDAO;
import io.github.jlwitthuhn.peanut.db.UserDAO;
import io.github.jlwitthuhn.peanut.err.GroupNotFoundException;
import io.github.jlwitthuhn.peanut.err.UserDetailsConflictException;
import io.github.jlwitthuhn.peanut.model.db.UserRow;
import io.github.jlwitthuhn.peanut.model.spring.PeanutUserDetails;
import lombok.RequiredArgsConstructor;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import org.springframework.security.provisioning.UserDetailsManager;
import org.springframework.stereotype.Component;
import org.springframework.transaction.annotation.Transactional;

import java.util.ArrayList;
import java.util.Collection;

@Component
@RequiredArgsConstructor
public class PeanutUserService implements UserDetailsManager
{
	private final GroupDAO groupDAO;
	private final MultiTableDAO multiTableDAO;
	private final UserDAO userDAO;
	private final GroupMembershipDAO groupMembershipDAO;

	@Override
	@Transactional
	public void createUser(UserDetails user)
	{
		if (!(user instanceof PeanutUserDetails peanutUserDetails))
		{
			throw new IllegalArgumentException("User is not PeanutUserDetails");
		}
		if (peanutUserDetails.getId().isPresent())
		{
			throw new IllegalArgumentException("User ID must not be set when creating a new user");
		}
		if (userDAO.selectRowByDisplayName(peanutUserDetails.getUsername()) != null)
		{
			throw new UserDetailsConflictException("Username is already in use");
		}
		if (userDAO.selectRowByEmail(peanutUserDetails.getEmail()) != null)
		{
			throw new UserDetailsConflictException("Email address is already in use");
		}
		userDAO.insertRow(
			peanutUserDetails.getUsername(),
			peanutUserDetails.getEmail(),
			peanutUserDetails.getPassword()
		);
		UserRow newRow = userDAO.selectRowByDisplayName(peanutUserDetails.getUsername());
		try
		{
			addRolesToUser(newRow.getId(), user.getAuthorities());
		}
		catch (GroupNotFoundException e)
		{
			throw new RuntimeException(e);
		}
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
		UserRow row = userDAO.selectRowByDisplayName(username);
		return row != null;
	}

	@Override
	public UserDetails loadUserByUsername(String username) throws UsernameNotFoundException
	{
		UserRow row = userDAO.selectRowByDisplayName(username);
		if (row == null)
		{
			throw new UsernameNotFoundException(username);
		}
		ArrayList<GrantedAuthority> groups = multiTableDAO.getGroupsByUserId(row.getId());
		return new PeanutUserDetails(row, groups);
	}

	private void addRolesToUser(long userId, Collection<? extends GrantedAuthority> groups) throws GroupNotFoundException
	{
		ArrayList<String> authorityStrings = new ArrayList<>();
		for (GrantedAuthority authority : groups)
		{
			authorityStrings.add(authority.getAuthority());
		}
		Collection<Long> authorityIds = groupDAO.getIdsFromNames(authorityStrings);
		groupMembershipDAO.insertGroupsForUser(userId, authorityIds);
	}
}
