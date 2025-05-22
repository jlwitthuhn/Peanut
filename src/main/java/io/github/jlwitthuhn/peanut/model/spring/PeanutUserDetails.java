// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.model.spring;

import io.github.jlwitthuhn.peanut.model.db.UserRow;
import lombok.Getter;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.core.userdetails.UserDetails;

import java.time.OffsetDateTime;
import java.util.ArrayList;
import java.util.Collection;
import java.util.List;
import java.util.Optional;
import java.util.stream.Collectors;

public class PeanutUserDetails implements UserDetails
{
	private final Long id;
	private final String displayName;
	private final String hashedPassword;
	private final OffsetDateTime createdTimestamp;
	private final OffsetDateTime updatedTimestamp;

	@Getter
	private final List<String> groups;

	@Getter
	private final String email;

	public PeanutUserDetails(UserRow user, List<String> groups)
	{
		id = user.getId();
		displayName = user.getDisplayName();
		email = user.getEmail();
		hashedPassword = user.getPassword();
		this.createdTimestamp = user.getCreatedTimestamp();
		this.updatedTimestamp = user.getUpdatedTimestamp();
		this.groups = groups;
	}

	public PeanutUserDetails(String displayName, String email, String hashedPassword)
	{
		List<String> groups = new ArrayList<>();
		groups.add("USER");

		this.id = null;
		this.displayName = displayName;
		this.email = email;
		this.hashedPassword = hashedPassword;
		this.createdTimestamp = null;
		this.updatedTimestamp = null;
		this.groups = groups;
	}

	public PeanutUserDetails(String displayName, String email, String hashedPassword, List<String> groups)
	{
		this.id = null;
		this.displayName = displayName;
		this.email = email;
		this.hashedPassword = hashedPassword;
		this.groups = groups;
		this.createdTimestamp = null;
		this.updatedTimestamp = null;
	}

	@Override
	public Collection<? extends GrantedAuthority> getAuthorities()
	{
		return groups.stream().map(SimpleGrantedAuthority::new).collect(Collectors.toList());
	}

	@Override
	public String getPassword()
	{
		return hashedPassword;
	}

	@Override
	public String getUsername()
	{
		return displayName;
	}

	public Optional<Long> getId()
	{
		return Optional.ofNullable(id);
	}

	public Optional<OffsetDateTime> getCreatedTimestamp()
	{
		return Optional.ofNullable(createdTimestamp);
	}

	public Optional<OffsetDateTime> getUpdatedTimestamp()
	{
		return Optional.ofNullable(updatedTimestamp);
	}
}
