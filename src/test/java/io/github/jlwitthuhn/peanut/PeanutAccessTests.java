// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut;

import io.github.jlwitthuhn.peanut.db.GroupDAO;
import io.github.jlwitthuhn.peanut.db.MultiTableDAO;
import io.github.jlwitthuhn.peanut.db.UserDAO;
import io.github.jlwitthuhn.peanut.interceptor.DatabaseInitInterceptor;
import io.github.jlwitthuhn.peanut.service.AdminService;
import io.github.jlwitthuhn.peanut.service.ConfigService;
import io.github.jlwitthuhn.peanut.service.DatabaseService;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.security.test.context.support.WithMockUser;
import org.springframework.test.context.bean.override.mockito.MockitoBean;
import org.springframework.test.web.servlet.MockMvc;
import org.springframework.test.web.servlet.request.MockMvcRequestBuilders;
import org.springframework.test.web.servlet.result.MockMvcResultMatchers;

@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
@AutoConfigureMockMvc
public class PeanutAccessTests
{
	@Autowired
	private MockMvc mvc;

	@MockitoBean
	private DatabaseInitInterceptor databaseInitInterceptor;

	@MockitoBean
	private JdbcTemplate jdbcTemplate;

	@MockitoBean
	private AdminService adminService;

	@MockitoBean
	private ConfigService configService;

	@MockitoBean
	private DatabaseService databaseService;

	@MockitoBean
	private GroupDAO groupDAO;

	@MockitoBean
	private MultiTableDAO multiTableDAO;

	@MockitoBean
	private UserDAO userDAO;

	@BeforeEach
	void commonSetup()
	{
		Mockito.when(databaseInitInterceptor.preHandle(Mockito.any(),Mockito.any(),Mockito.any())).thenReturn(true);
	}

	@Test
	@WithMockUser(authorities = "ADMIN")
	void adminAccessAllowedTest() throws Exception
	{
		mvc.perform(MockMvcRequestBuilders.get("/admin"))
			.andExpect(MockMvcResultMatchers.status().isOk());
	}

	@Test
	@WithMockUser(authorities = "USER")
	void adminAccessDeniedTest() throws Exception
	{
		mvc.perform(MockMvcRequestBuilders.get("/admin"))
			.andExpect(MockMvcResultMatchers.status().isForbidden());
	}
}
