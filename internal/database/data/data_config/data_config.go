// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package data_config

import "database/sql"

var sqlCreateTableInt = `
	CREATE TABLE config_int (
		name VARCHAR(255) PRIMARY KEY,
		value BIGINT NOT NULL,
		_created TIMESTAMP WITH TIME ZONE NOT NULL,
		_updated TIMESTAMP WITH TIME ZONE NOT NULL
	);

	CREATE TRIGGER
		config_int_trigger_created_updated_before_insert
	BEFORE INSERT ON
		config_int
	FOR EACH ROW EXECUTE FUNCTION
		fn_created_updated_before_insert();

	CREATE TRIGGER
		config_int_trigger_created_updated_before_update
	BEFORE UPDATE ON
		config_int
	FOR EACH ROW EXECUTE FUNCTION
		fn_created_updated_before_update();
`

func CreateDBObjects(tx *sql.Tx) error {
	_, intErr := tx.Exec(sqlCreateTableInt)
	if intErr != nil {
		return intErr
	}
	return nil
}
