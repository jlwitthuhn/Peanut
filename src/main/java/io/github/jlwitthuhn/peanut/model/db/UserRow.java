package io.github.jlwitthuhn.peanut.model.db;

import lombok.Data;

@Data
public class UserRow
{
	private final long id;
	private final String name;
	private final String password;
}
