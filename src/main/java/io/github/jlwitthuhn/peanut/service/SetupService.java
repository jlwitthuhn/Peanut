// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.service;

import io.github.jlwitthuhn.peanut.cfg.ConfigKeyNames;
import io.github.jlwitthuhn.peanut.db.ConfigDAO;
import io.github.jlwitthuhn.peanut.db.EventLogDAO;
import io.github.jlwitthuhn.peanut.db.GroupDAO;
import io.github.jlwitthuhn.peanut.db.GroupMembershipDAO;
import io.github.jlwitthuhn.peanut.db.MetaDAO;
import io.github.jlwitthuhn.peanut.db.UserDAO;
import io.github.jlwitthuhn.peanut.err.DBCreationDependencyNotSatisfiedException;
import io.github.jlwitthuhn.peanut.model.spring.PeanutUserDetails;
import io.github.jlwitthuhn.peanut.security.PeanutUserService;
import lombok.RequiredArgsConstructor;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.Instant;
import java.util.ArrayList;
import java.util.List;

@Service
@RequiredArgsConstructor
public class SetupService
{
	private final DatabaseService dbService;
	private final PeanutUserService peanutUserService;
	private final PasswordEncoder passwordEncoder;

	private final ConfigDAO configDAO;
	private final EventLogDAO eventLogDAO;
	private final GroupDAO groupDAO;
	private final MetaDAO metaDAO;
	private final UserDAO userDAO;
	private final GroupMembershipDAO groupMembershipDAO;

	private static final Logger logger = LoggerFactory.getLogger(SetupService.class);

	@Transactional
	public void initializeDatabase(String adminName, String plainAdminPassword, String adminEmail) throws DBCreationDependencyNotSatisfiedException
	{
		if (dbService.doesTableExist("config_int"))
		{
			throw new DBCreationDependencyNotSatisfiedException("The database is already set up. No changes have been made.");
		}

		logger.info("Initializing database...");

		initDatabaseObjects();
		initGroups();
		initConfig();
		initAdminUser(adminName, plainAdminPassword, adminEmail);
	}

	private void initAdminUser(String username, String plainPassword, String email)
	{
		logger.info("Initializing admin user...");

		String hashedPassword = passwordEncoder.encode(plainPassword);
		List<String> groups = new ArrayList<>();
		groups.add("TURBO_ADMIN");
		groups.add("ADMIN");
		groups.add("USER");
		PeanutUserDetails userDetails = new PeanutUserDetails(username, email, hashedPassword, groups);
		peanutUserService.createUser(userDetails);
	}

	private void initGroups()
	{
		logger.info("Initializing groups...");

		groupDAO.insertRow("TURBO_ADMIN", "Full control over everything, access cannot be limited.", true);
		groupDAO.insertRow("ADMIN", "Full control over everything by default, access can be limited with permissions.", true);
		groupDAO.insertRow("USER", "Standard role given to all users. This carries no special permissions.", true);
	}

	private void initConfig()
	{
		logger.info("Initializing config...");
		configDAO.upsertLongByName(ConfigKeyNames.INITIALIZED_TIME_INT, Instant.now().getEpochSecond());
		configDAO.upsertLongByName(ConfigKeyNames.SCHEMA_VERSION_INT, 1);
		configDAO.upsertStringByName(ConfigKeyNames.WELCOME_MESSAGE_STR, "Welcome to my underground lair.");
	}

	private void initDatabaseObjects() throws DBCreationDependencyNotSatisfiedException
	{
		logger.info("Initializing database objects...");
		metaDAO.createDatabaseObjects();
		groupDAO.createDatabaseObjects();
		configDAO.createDatabaseObjects();
		userDAO.createDatabaseObjects();
		groupMembershipDAO.createDatabaseObjects();
		eventLogDAO.createDatabaseObjects();
	}
}
