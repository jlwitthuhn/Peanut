// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.controller;

import io.github.jlwitthuhn.peanut.cfg.PeanutGlobals;
import io.github.jlwitthuhn.peanut.db.MetaDAO;
import io.github.jlwitthuhn.peanut.util.TimeUtil;
import io.github.jlwitthuhn.peanut.util.Tuple2;
import lombok.RequiredArgsConstructor;
import org.springframework.core.SpringVersion;
import org.springframework.security.access.annotation.Secured;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.servlet.ModelAndView;

import java.lang.management.ManagementFactory;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;

@Controller
@RequestMapping("/admin")
@Secured({"ROLE_ADMIN"})
@RequiredArgsConstructor
public class AdminController
{
	private final MetaDAO metaDAO;

	@GetMapping("")
	ModelAndView adminIndex(Map<String, Object> model)
	{
		long uptimeMs = ManagementFactory.getRuntimeMXBean().getUptime();
		String uptimeStr = TimeUtil.formatMillisecondsAsDDHHMMSS(uptimeMs);

		long maxMemory = Runtime.getRuntime().maxMemory();
		String maxMemoryStr = maxMemory / 1024 / 1024 + " MB";
		long totalMemory = Runtime.getRuntime().totalMemory();
		String totalMemoryStr = totalMemory / 1024 / 1024 + " MB";

		double loadAverage = ManagementFactory.getOperatingSystemMXBean().getSystemLoadAverage();

		List<Tuple2<String, String>> environment = new ArrayList<>();
		environment.add(new Tuple2<>("Peanut version", PeanutGlobals.PEANUT_VERSION));
		environment.add(new Tuple2<>("Host OS", System.getProperty("os.name") + " / " + System.getProperty("os.arch")));
		environment.add(new Tuple2<>("Java version", System.getProperty("java.version")));
		environment.add(new Tuple2<>("Java vendor", System.getProperty("java.vendor")));
		environment.add(new Tuple2<>("Spring version", SpringVersion.getVersion()));
		environment.add(new Tuple2<>("Postgres version", metaDAO.getServerVersion()));
		environment.add(new Tuple2<>("Database size", metaDAO.getDatabaseSize()));
		environment.add(new Tuple2<>("JVM memory available", maxMemoryStr));
		environment.add(new Tuple2<>("JVM memory used", totalMemoryStr));
		if (loadAverage > 0)
		{
			// Value under 0 indicates that the host does not support this feature
			environment.add(new Tuple2<>("Load average", String.valueOf(loadAverage)));
		}
		environment.add(new Tuple2<>("Server uptime", uptimeStr));
		model.put("environment", environment);
		return new ModelAndView("admin/index.html", model);
	}
}
