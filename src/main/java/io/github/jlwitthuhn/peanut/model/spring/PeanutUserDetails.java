// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.model.spring;

import io.github.jlwitthuhn.peanut.model.db.UserRow;
import lombok.Getter;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.userdetails.UserDetails;

import java.util.Collection;
import java.util.List;

public class PeanutUserDetails implements UserDetails
{
	private final String displayName;
	private final String password;

	@Getter
	private final String email;

	public PeanutUserDetails(UserRow user)
	{
		displayName = user.getDisplayName();
		email = user.getEmail();
		password = user.getPassword();
	}

	public PeanutUserDetails(String displayName, String email, String password)
	{
		this.displayName = displayName;
		this.email = email;
		this.password = password;
	}

	@Override
	public Collection<? extends GrantedAuthority> getAuthorities() {
		return List.of();
	}

	@Override
	public String getPassword() {
		return password;
	}

	@Override
	public String getUsername() {
		return displayName;
	}
}
