// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.cfg;

import lombok.RequiredArgsConstructor;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.boot.context.event.ApplicationReadyEvent;
import org.springframework.context.event.EventListener;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.stereotype.Component;

@Component
@RequiredArgsConstructor
public class JdbcConnectOnStartupHelper
{
	private final JdbcTemplate jdbcTemplate;

	private static final Logger logger = LoggerFactory.getLogger(JdbcConnectOnStartupHelper.class);

	@EventListener(ApplicationReadyEvent.class)
	public void connect()
	{
		logger.info("Connecting to the database...");
		jdbcTemplate.execute("SELECT 1");
	}
}
