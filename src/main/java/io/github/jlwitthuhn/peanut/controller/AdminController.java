// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.controller;

import io.github.jlwitthuhn.peanut.cfg.PeanutGlobals;
import io.github.jlwitthuhn.peanut.util.TimeUtil;
import io.github.jlwitthuhn.peanut.util.Tuple2;
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
public class AdminController
{
	@GetMapping("")
	ModelAndView adminIndex(Map<String, Object> model)
	{
		List<Tuple2<String, String>> environment = new ArrayList<>();
		environment.add(new Tuple2<>("Peanut version", PeanutGlobals.PEANUT_VERSION));
		environment.add(new Tuple2<>("Host OS", System.getProperty("os.name") + " / " + System.getProperty("os.arch")));
		environment.add(new Tuple2<>("Java version", System.getProperty("java.version")));
		environment.add(new Tuple2<>("Java vendor", System.getProperty("java.vendor")));
		environment.add(new Tuple2<>("Spring version", SpringVersion.getVersion()));
		long uptimeMs = ManagementFactory.getRuntimeMXBean().getUptime();
		String uptimeStr = TimeUtil.formatMillisecondsAsDDHHMMSS(uptimeMs);
		environment.add(new Tuple2<>("Server uptime", uptimeStr));
		model.put("environment", environment);
		return new ModelAndView("admin/index.html", model);
	}
}
