package io.github.jlwitthuhn.peanut.cfg;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.http.HttpMethod;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.config.annotation.web.configurers.AbstractHttpConfigurer;
import org.springframework.security.web.SecurityFilterChain;

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
		http.csrf(AbstractHttpConfigurer::disable);

		return http.build();
	}
}
