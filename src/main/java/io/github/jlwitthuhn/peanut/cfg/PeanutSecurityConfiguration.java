// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.cfg;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.http.HttpMethod;
import org.springframework.security.config.annotation.method.configuration.EnableMethodSecurity;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.config.core.GrantedAuthorityDefaults;
import org.springframework.security.crypto.argon2.Argon2PasswordEncoder;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.DelegatingPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.security.web.SecurityFilterChain;
import org.springframework.security.web.csrf.HttpSessionCsrfTokenRepository;
import org.springframework.security.web.firewall.HttpFirewall;
import org.springframework.security.web.firewall.StrictHttpFirewall;
import org.springframework.security.web.savedrequest.HttpSessionRequestCache;

import java.util.HashMap;
import java.util.Map;

@Configuration
@EnableMethodSecurity(securedEnabled = true)
@EnableWebSecurity
public class PeanutSecurityConfiguration
{
	@Bean
	public SecurityFilterChain securityFilterChain(HttpSecurity http) throws Exception
	{
		http.authorizeHttpRequests(
			(authorizeHttpRequests) ->
				authorizeHttpRequests
					.requestMatchers(
						HttpMethod.GET,
						"/error",
						"/favicon*",
						"/login",
						"/logout/success",
						"/register",
						"/setup",
						"/design"
					).permitAll()
					.requestMatchers(HttpMethod.POST, "/login", "/register", "/setup").permitAll()
					.anyRequest().authenticated()
		);
		http.formLogin(
			(formLogin) ->
				formLogin
					.defaultSuccessUrl("/login/success")
					.failureUrl("/login/failure")
					.loginPage("/login")
		);
		http.logout(
			(logout) ->
				logout
					.logoutUrl("/logout")
					.logoutSuccessUrl("/logout/success")
		);
		http.csrf(
			(csrf) ->
				csrf.csrfTokenRepository(new HttpSessionCsrfTokenRepository())
		);

		// Disable the 'continue' parameter being added after login
		HttpSessionRequestCache requestCache = new HttpSessionRequestCache();
		requestCache.setMatchingRequestParameterName(null);
		http.requestCache(
			(cache) ->
				cache.requestCache(requestCache)
		);

		return http.build();
	}

	@Bean
	public GrantedAuthorityDefaults grantedAuthorityDefaults() {
		// This is the role prefix, default is 'ROLE_'
		return new GrantedAuthorityDefaults("");
	}

	@Bean
	public HttpFirewall httpFirewall()
	{
		StrictHttpFirewall httpFirewall = new StrictHttpFirewall();
		httpFirewall.setAllowUrlEncodedPercent(true);
		httpFirewall.setAllowUrlEncodedSlash(true);
		return httpFirewall;
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
