// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.service;

import io.github.jlwitthuhn.peanut.cfg.ConfigKeyNames;
import io.github.jlwitthuhn.peanut.db.AuthorityDAO;
import io.github.jlwitthuhn.peanut.db.ConfigDAO;
import io.github.jlwitthuhn.peanut.db.MetaDAO;
import io.github.jlwitthuhn.peanut.db.UserAuthorityDAO;
import io.github.jlwitthuhn.peanut.db.UserDAO;
import io.github.jlwitthuhn.peanut.err.DBCreationDependencyNotSatisfiedException;
import io.github.jlwitthuhn.peanut.model.spring.PeanutUserDetails;
import io.github.jlwitthuhn.peanut.security.PeanutUserService;
import lombok.RequiredArgsConstructor;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.Instant;
import java.util.ArrayList;

@Service
@RequiredArgsConstructor
public class SetupService
{
	private final PasswordEncoder passwordEncoder;
	private final PeanutUserService peanutUserService;

	private final AuthorityDAO authorityDAO;
	private final ConfigDAO configDAO;
	private final MetaDAO metaDAO;
	private final UserDAO userDAO;
	private final UserAuthorityDAO userAuthorityDAO;

	@Transactional
	public void initializeDatabase(String adminName, String plainAdminPassword, String adminEmail) throws DBCreationDependencyNotSatisfiedException
	{
		if (metaDAO.doesTableExist("config_int"))
		{
			throw new DBCreationDependencyNotSatisfiedException("The database is already set up. No changes have been made.");
		}

		initDatabaseObjects();
		initAuthorities();
		initConfig();
		initAdminUser(adminName, plainAdminPassword, adminEmail);
	}

	private void initAdminUser(String username, String plainPassword, String email)
	{
		String hashedPassword = passwordEncoder.encode(plainPassword);
		ArrayList<GrantedAuthority> authorities = new ArrayList<GrantedAuthority>();
		authorities.add(new SimpleGrantedAuthority("ROLE_TURBO_ADMIN"));
		authorities.add(new SimpleGrantedAuthority("ROLE_ADMIN"));
		authorities.add(new SimpleGrantedAuthority("ROLE_USER"));
		PeanutUserDetails userDetails = new PeanutUserDetails(username, email, hashedPassword, authorities);
		peanutUserService.createUser(userDetails);
	}

	private void initAuthorities()
	{
		authorityDAO.insertRow("ROLE_TURBO_ADMIN", true);
		authorityDAO.insertRow("ROLE_ADMIN", true);
		authorityDAO.insertRow("ROLE_USER", true);
	}

	private void initConfig()
	{
		configDAO.setLong(ConfigKeyNames.INITIALIZED_TIME, Instant.now().getEpochSecond());
		configDAO.setLong(ConfigKeyNames.SCHEMA_VERSION, 1);
	}

	private void initDatabaseObjects() throws DBCreationDependencyNotSatisfiedException
	{
		metaDAO.createDatabaseObjects();
		authorityDAO.createDatabaseObjects();
		configDAO.createDatabaseObjects();
		userDAO.createDatabaseObjects();
		userAuthorityDAO.createDatabaseObjects();
	}
}
