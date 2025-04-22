// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.cfg;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.http.HttpMethod;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.crypto.argon2.Argon2PasswordEncoder;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.DelegatingPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.security.web.SecurityFilterChain;
import org.springframework.security.web.csrf.HttpSessionCsrfTokenRepository;

import java.util.HashMap;
import java.util.Map;

@Configuration
@EnableWebSecurity
public class PeanutSecurityConfiguration
{
	@Bean
	public SecurityFilterChain securityFilterChain(HttpSecurity http) throws Exception
	{
		http.authorizeHttpRequests(
			(authorizeHttpRequests) ->
				authorizeHttpRequests
					.requestMatchers(HttpMethod.GET, "/login", "/setup", "/design").permitAll()
					.requestMatchers(HttpMethod.POST, "/login", "/setup").permitAll()
					.anyRequest().authenticated()
		);
		http.formLogin(
			(formLogin) ->
				formLogin
					.defaultSuccessUrl("/")
					.loginPage("/login")
		);
		http.csrf(
			(csrf) ->
				csrf.csrfTokenRepository(new HttpSessionCsrfTokenRepository())
		);

		return http.build();
	}

	@Bean
	public PasswordEncoder passwordEncoder()
	{
		Map<String, PasswordEncoder> encoders = new HashMap<>();
		encoders.put("argon2", Argon2PasswordEncoder.defaultsForSpringSecurity_v5_8());
		encoders.put("bcrypt", new BCryptPasswordEncoder());
		return new DelegatingPasswordEncoder("bcrypt", encoders);
	}
}
