// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.model.db;

import lombok.Data;

import java.time.OffsetDateTime;

@Data
public class UserRow
{
	private final long id;
	private final String displayName;
	private final String email;
	private final String password;
	private final OffsetDateTime createdTimestamp;
	private final OffsetDateTime updatedTimestamp;
}
