package io.github.jlwitthuhn.peanut.db;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.stereotype.Component;

@Component
public class ConfigDAO
{
	private final JdbcTemplate jdbcTemplate;

	@Autowired
	public ConfigDAO(JdbcTemplate jdbcTemplate)
	{
		this.jdbcTemplate = jdbcTemplate;
	}

	public Long getLong(String name)
	{
		return jdbcTemplate.queryForObject("SELECT value FROM config_int WHERE name = ?", Long.class, name);
	}
}
