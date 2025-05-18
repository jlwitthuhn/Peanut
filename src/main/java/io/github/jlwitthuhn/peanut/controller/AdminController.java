// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.controller;

import io.github.jlwitthuhn.peanut.cfg.ConfigKeyNames;
import io.github.jlwitthuhn.peanut.cfg.PeanutGlobals;
import io.github.jlwitthuhn.peanut.db.ConfigDAO;
import io.github.jlwitthuhn.peanut.db.MetaDAO;
import io.github.jlwitthuhn.peanut.model.form.AdminDebugCreateUsersForm;
import io.github.jlwitthuhn.peanut.model.spring.PeanutUserDetails;
import io.github.jlwitthuhn.peanut.security.PeanutUserService;
import io.github.jlwitthuhn.peanut.service.AdminService;
import io.github.jlwitthuhn.peanut.util.TimeUtil;
import io.github.jlwitthuhn.peanut.util.Tuple2;
import io.github.jlwitthuhn.peanut.util.ViewShortcuts;
import lombok.RequiredArgsConstructor;
import org.springframework.core.SpringVersion;
import org.springframework.http.HttpStatus;
import org.springframework.security.access.annotation.Secured;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.ModelAttribute;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.servlet.ModelAndView;

import java.lang.management.ManagementFactory;
import java.text.SimpleDateFormat;
import java.time.Instant;
import java.time.OffsetDateTime;
import java.time.ZoneId;
import java.util.ArrayList;
import java.util.Date;
import java.util.List;
import java.util.Map;

@Controller
@RequestMapping("/admin")
@Secured({"ROLE_ADMIN"})
@RequiredArgsConstructor
public class AdminController
{
	private final AdminService adminService;

	private final ConfigDAO configDAO;
	private final MetaDAO metaDAO;

	@GetMapping("")
	ModelAndView adminIndex(Map<String, Object> model)
	{
		Long initCount = configDAO.getLong(ConfigKeyNames.INITIALIZED_TIME);
		Instant initInstant = Instant.ofEpochSecond(initCount);
		OffsetDateTime initDateTime = OffsetDateTime.ofInstant(initInstant, ZoneId.of("UTC"));
		String initTimeStr = TimeUtil.formatOffsetDateTime(initDateTime);

		long uptimeMs = ManagementFactory.getRuntimeMXBean().getUptime();
		String uptimeStr = TimeUtil.formatMillisecondsAsDDHHMMSS(uptimeMs);

		long maxMemory = Runtime.getRuntime().maxMemory();
		String maxMemoryStr = maxMemory / 1024 / 1024 + " MB";
		long totalMemory = Runtime.getRuntime().totalMemory();
		String totalMemoryStr = totalMemory / 1024 / 1024 + " MB";

		double loadAverage = ManagementFactory.getOperatingSystemMXBean().getSystemLoadAverage();

		List<Tuple2<String, String>> versionDetails = new ArrayList<>();
		versionDetails.add(new Tuple2<>("Peanut version", PeanutGlobals.PEANUT_VERSION));
		versionDetails.add(new Tuple2<>("Host OS", System.getProperty("os.name") + " / " + System.getProperty("os.arch")));
		versionDetails.add(new Tuple2<>("Java version", System.getProperty("java.version")));
		versionDetails.add(new Tuple2<>("Java vendor", System.getProperty("java.vendor")));
		versionDetails.add(new Tuple2<>("Spring version", SpringVersion.getVersion()));
		versionDetails.add(new Tuple2<>("Postgres version", metaDAO.getServerVersion()));

		List<Tuple2<String, String>> runtimeDetails = new ArrayList<>();
		runtimeDetails.add(new Tuple2<>("DB initialized", initTimeStr));
		runtimeDetails.add(new Tuple2<>("Database size", metaDAO.getDatabaseSize()));
		runtimeDetails.add(new Tuple2<>("JVM max memory", maxMemoryStr));
		runtimeDetails.add(new Tuple2<>("JVM total memory", totalMemoryStr));
		if (loadAverage > 0)
		{
			// Value under 0 indicates that the host does not support this feature
			runtimeDetails.add(new Tuple2<>("Load average", String.valueOf(loadAverage)));
		}
		runtimeDetails.add(new Tuple2<>("Server uptime", uptimeStr));

		model.put("version_details", versionDetails);
		model.put("runtime_details", runtimeDetails);
		return new ModelAndView("admin/index.html", model);
	}

	@GetMapping("/debug")
	ModelAndView debugIndex()
	{
		return new ModelAndView("admin/debug.html");
	}

	@PostMapping("/debug/create_users")
	ModelAndView debugCreateUsers(@ModelAttribute AdminDebugCreateUsersForm form)
	{
		if (form == null)
		{
			return ViewShortcuts.simpleMessage("Error", "Form is invalid.", HttpStatus.BAD_REQUEST);
		}

		if (!form.getValidationErrors().isEmpty())
		{
			return ViewShortcuts.simpleMessage("Error", form.getValidationErrors().getFirst(), HttpStatus.BAD_REQUEST);
		}

		Integer count = form.getCountInt();
		adminService.createDebugUsers(count, form.getPrefix(), form.getPassword());

		return ViewShortcuts.simpleMessage("Success", "Successfully created " + count + " test accounts.");
	}

	@GetMapping("/users")
	ModelAndView usersIndex()
	{
		return new ModelAndView("admin/users.html");
	}
}
