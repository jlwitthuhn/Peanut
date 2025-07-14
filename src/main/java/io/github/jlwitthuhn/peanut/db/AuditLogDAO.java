// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.db;

import io.github.jlwitthuhn.peanut.err.DBCreationDependencyNotSatisfiedException;
import io.github.jlwitthuhn.peanut.model.db.AuditLogEventType;
import io.github.jlwitthuhn.peanut.model.db.AuditLogTargetType;
import io.github.jlwitthuhn.peanut.service.DatabaseService;
import lombok.RequiredArgsConstructor;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.stereotype.Component;

@Component
@RequiredArgsConstructor
public class AuditLogDAO
{
	public static final String TABLE_NAME = "audit_log";

	private final DatabaseService dbService;

	private final JdbcTemplate jdbcTemplate;

	public void createDatabaseObjects() throws DBCreationDependencyNotSatisfiedException
	{
		if (dbService.doesTableExist(TABLE_NAME))
		{
			throw new DBCreationDependencyNotSatisfiedException("Table '" + TABLE_NAME + "' cannot be created because it already exists");
		}
		if (!dbService.doesTableExist(UserDAO.TABLE_NAME))
		{
			throw new DBCreationDependencyNotSatisfiedException("Table '" + TABLE_NAME + "' requires that table '" + UserDAO.TABLE_NAME + "' exists");
		}
		final String SQL_TABLE = """
			CREATE TABLE audit_log (
			    id BIGSERIAL PRIMARY KEY,
			    source_user_id BIGINT REFERENCES users(id) NOT NULL,
			    target_id BIGINT,
			    target_id_type VARCHAR(256) NOT NULL,
			    event_type VARCHAR(256) NOT NULL,
			    message VARCHAR(1024) NOT NULL,
			    _created TIMESTAMP WITH TIME ZONE NOT NULL,
			    _updated TIMESTAMP WITH TIME ZONE NOT NULL
			);
			""";
		jdbcTemplate.execute(SQL_TABLE);
		final String SQL_TRIGGER_BEFORE_INSERT = """
			CREATE TRIGGER
				audit_log_trigger_created_updated_before_insert
			BEFORE INSERT ON
				audit_log
			FOR EACH ROW EXECUTE FUNCTION
				fn_created_updated_before_insert();
			""";
		jdbcTemplate.execute(SQL_TRIGGER_BEFORE_INSERT);
		final String SQL_TRIGGER_BEFORE_UPDATE = """
			CREATE TRIGGER
				audit_log_trigger_created_updated_before_update
			BEFORE UPDATE ON
				audit_log
			FOR EACH ROW EXECUTE FUNCTION
				fn_created_updated_before_update();
			""";
		jdbcTemplate.execute(SQL_TRIGGER_BEFORE_UPDATE);
	}

	public void insertEvent(Long srcUserId, Long targetId, AuditLogTargetType targetType, AuditLogEventType eventType, String message)
	{
		final String SQL = """
			INSERT INTO
				audit_log (source_user_id, target_id, target_id_type, event_type, message)
			VALUES
				(?, ?, ?, ?, ?)
			""";
		jdbcTemplate.update(SQL, srcUserId, targetId, targetType.name(), eventType.name(), message);
	}
}
