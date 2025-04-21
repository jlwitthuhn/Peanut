package io.github.jlwitthuhn.peanut;

import io.github.jlwitthuhn.peanut.controller.RootController;
import io.github.jlwitthuhn.peanut.controller.SetupController;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;

import static org.assertj.core.api.Assertions.assertThat;

@SpringBootTest
class PeanutApplicationTests
{
	@Autowired
	private RootController rootController;

	@Autowired
	private SetupController setupController;

	@Test
	void controllersLoad()
	{
		assertThat(rootController).isNotNull();
		assertThat(setupController).isNotNull();
	}
}
