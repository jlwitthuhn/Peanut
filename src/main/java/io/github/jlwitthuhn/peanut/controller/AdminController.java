// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.controller;

import io.github.jlwitthuhn.peanut.cfg.ConfigKeyNames;
import io.github.jlwitthuhn.peanut.cfg.PeanutGlobals;
import io.github.jlwitthuhn.peanut.db.GroupDAO;
import io.github.jlwitthuhn.peanut.db.MultiTableDAO;
import io.github.jlwitthuhn.peanut.db.UserDAO;
import io.github.jlwitthuhn.peanut.model.db.GroupRow;
import io.github.jlwitthuhn.peanut.model.db.UserRow;
import io.github.jlwitthuhn.peanut.model.form.AdminDebugCreateUsersForm;
import io.github.jlwitthuhn.peanut.model.form.AdminGeneralWelcomeMessageForm;
import io.github.jlwitthuhn.peanut.model.form.AdminUsersListByGroupForm;
import io.github.jlwitthuhn.peanut.model.form.AdminUsersSearchByNamePatternForm;
import io.github.jlwitthuhn.peanut.service.AdminService;
import io.github.jlwitthuhn.peanut.service.ConfigService;
import io.github.jlwitthuhn.peanut.service.DatabaseService;
import io.github.jlwitthuhn.peanut.util.TimeUtil;
import io.github.jlwitthuhn.peanut.util.Tuple2;
import io.github.jlwitthuhn.peanut.util.ViewShortcuts;
import lombok.RequiredArgsConstructor;
import org.springframework.boot.SpringBootVersion;
import org.springframework.core.SpringVersion;
import org.springframework.http.HttpStatus;
import org.springframework.security.access.annotation.Secured;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.ModelAttribute;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.servlet.ModelAndView;
import org.springframework.web.servlet.view.RedirectView;

import java.lang.management.ManagementFactory;
import java.net.URLEncoder;
import java.nio.charset.StandardCharsets;
import java.time.Instant;
import java.time.OffsetDateTime;
import java.time.ZoneId;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

@Controller
@RequestMapping("/admin")
@Secured({"ADMIN"})
@RequiredArgsConstructor
public class AdminController
{
	private final AdminService adminService;
	private final ConfigService configService;
	private final DatabaseService dbService;

	private final GroupDAO groupDAO;
	private final MultiTableDAO multiTableDAO;
	private final UserDAO userDAO;

	@GetMapping("")
	ModelAndView adminIndex(Map<String, Object> model)
	{
		Long initCount = configService.getLong(ConfigKeyNames.INITIALIZED_TIME_INT);
		Instant initInstant = Instant.ofEpochSecond(initCount);
		OffsetDateTime initDateTime = OffsetDateTime.ofInstant(initInstant, ZoneId.of("UTC"));
		String initTimeStr = TimeUtil.formatOffsetDateTime(initDateTime);

		long uptimeMs = ManagementFactory.getRuntimeMXBean().getUptime();
		String uptimeStr = TimeUtil.formatMillisecondsAsDDHHMMSS(uptimeMs);

		long maxMemory = Runtime.getRuntime().maxMemory();
		String maxMemoryStr = maxMemory / 1024 / 1024 + " MB";
		long totalMemory = Runtime.getRuntime().totalMemory();
		String totalMemoryStr = totalMemory / 1024 / 1024 + " MB";

		int cpuCount = ManagementFactory.getOperatingSystemMXBean().getAvailableProcessors();
		double systemLoad = ManagementFactory.getOperatingSystemMXBean().getSystemLoadAverage();

		List<Tuple2<String, String>> versionDetails = new ArrayList<>();
		versionDetails.add(new Tuple2<>("Peanut version", PeanutGlobals.PEANUT_VERSION));
		versionDetails.add(new Tuple2<>("Java version", System.getProperty("java.version")));
		versionDetails.add(new Tuple2<>("Java vendor", System.getProperty("java.vendor")));
		versionDetails.add(new Tuple2<>("Host OS", System.getProperty("os.name") + " / " + System.getProperty("os.arch")));
		versionDetails.add(new Tuple2<>("Spring Boot version", SpringBootVersion.getVersion()));
		versionDetails.add(new Tuple2<>("Spring version", SpringVersion.getVersion()));
		versionDetails.add(new Tuple2<>("Postgres version", dbService.getServerVersion()));

		List<Tuple2<String, String>> runtimeDetails = new ArrayList<>();
		runtimeDetails.add(new Tuple2<>("DB initialized", initTimeStr));
		runtimeDetails.add(new Tuple2<>("Database size", dbService.getDatabaseSize()));
		runtimeDetails.add(new Tuple2<>("JVM max memory", maxMemoryStr));
		runtimeDetails.add(new Tuple2<>("JVM total memory", totalMemoryStr));
		runtimeDetails.add(new Tuple2<>("CPU cores", String.valueOf(cpuCount)));
		if (systemLoad > 0)
		{
			// Value under 0 indicates that the host does not support this feature
			runtimeDetails.add(new Tuple2<>("System load", String.valueOf(systemLoad)));
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

	@GetMapping("/front_page")
	ModelAndView generalFrontPage(Map<String, Object> model)
	{
		String welcomeMessage = configService.getString(ConfigKeyNames.WELCOME_MESSAGE_STR);
		model.put("welcomeMessage", welcomeMessage);
		return new ModelAndView("admin/front_page.html");
	}

	@PostMapping("/front_page/welcome_message")
	ModelAndView generalWelcomeMessagePost(@ModelAttribute AdminGeneralWelcomeMessageForm form, Map<String, Object> model)
	{
		if (form.getConfirm() == null || !form.getConfirm())
		{
			return renderSimpleMessage("Error - Message Not Set", "You must check the 'Confirm' box to set the welcome message.", HttpStatus.BAD_REQUEST, model);
		}
		configService.setString(ConfigKeyNames.WELCOME_MESSAGE_STR, form.getMessage());
		return ViewShortcuts.simpleRedirect("/admin/front_page");
	}

	@GetMapping("/groups")
	ModelAndView groupsIndex()
	{
		return new ModelAndView("admin/groups.html");
	}

	@GetMapping("/groups/list")
	ModelAndView groupsListGet()
	{
		RedirectView view = new RedirectView("/admin/groups");
		view.setStatusCode(HttpStatus.SEE_OTHER);
		return new ModelAndView(view);
	}

	@GetMapping("/groups/list/all")
	ModelAndView groupsListPost(Map<String, Object> model)
	{
		List<GroupRow> groupRows = groupDAO.selectAll();
		List<Map<String, String>> groups = new ArrayList<>();
		for (GroupRow groupRow : groupRows)
		{
			String createdTimestamp = TimeUtil.formatOffsetDateTime(groupRow.getCreated());
			Map<String, String> thisGroup = new HashMap<>();
			thisGroup.put("id", String.valueOf(groupRow.getId()));
			thisGroup.put("name", groupRow.getName());
			thisGroup.put("description", groupRow.getDescription());
			thisGroup.put("system_owned", String.valueOf(groupRow.getSystemOwned()));
			thisGroup.put("created", createdTimestamp);
			groups.add(thisGroup);
		}
		model.put("groups", groups);
		return new ModelAndView("admin/groups_list.html", model);
	}

	@GetMapping("/users")
	ModelAndView usersIndex(Map<String, Object> model)
	{
		List<GroupRow> allGroupRows = groupDAO.selectAll();
		List<String> groups = allGroupRows.stream().map(GroupRow::getName).toList();
		model.put("groups", groups);
		return new ModelAndView("admin/users.html", model);
	}

	@GetMapping("/users/list")
	ModelAndView usersListGet()
	{
		RedirectView view = new RedirectView("/admin/users");
		view.setStatusCode(HttpStatus.SEE_OTHER);
		return new ModelAndView(view);
	}

	@PostMapping("/users/list")
	ModelAndView usersListPost(
		@ModelAttribute AdminUsersListByGroupForm listByGroupForm,
		@ModelAttribute AdminUsersSearchByNamePatternForm listByNameForm
	)
	{
		boolean listByGroupValid = listByGroupForm != null && listByGroupForm.getGroupName() != null && !listByGroupForm.getGroupName().isEmpty();
		if (listByGroupValid)
		{
			String encodedPattern = URLEncoder.encode(listByGroupForm.getGroupName(), StandardCharsets.UTF_8);
			RedirectView view = new RedirectView("/admin/users/list/by_group/" + encodedPattern);
			view.setStatusCode(HttpStatus.SEE_OTHER);
			return new ModelAndView(view);
		}
		boolean listByNameValid = listByNameForm != null && listByNameForm.getPattern() != null && !listByNameForm.getPattern().isEmpty();
		if (listByNameValid)
		{
			String encodedPattern = URLEncoder.encode(listByNameForm.getPattern(), StandardCharsets.UTF_8);
			RedirectView view = new RedirectView("/admin/users/list/by_name/" + encodedPattern);
			view.setStatusCode(HttpStatus.SEE_OTHER);
			return new ModelAndView(view);
		}

		return ViewShortcuts.simpleMessage("Error", "Invalid search parameters.", HttpStatus.BAD_REQUEST);
	}

	@GetMapping("/users/list/all")
	ModelAndView usersListAll(Map<String, Object> model)
	{
		List<UserRow> userRows = userDAO.selectAll();
		return renderUserList(userRows, model);
	}

	@GetMapping("/users/list/by_group/{groupName}")
	ModelAndView usersListByGroup(@PathVariable String groupName, Map<String, Object> model)
	{
		if (!groupDAO.doesGroupExist(groupName))
		{
			return ViewShortcuts.simpleMessage("Error", "Selected group does not exist.", HttpStatus.BAD_REQUEST);
		}
		List<UserRow> userRows = multiTableDAO.getUsersByGroupName(groupName);
		return renderUserList(userRows, model);
	}

	@GetMapping("/users/list/by_name/{pattern}")
	ModelAndView usersListByName(@PathVariable String pattern, Map<String, Object> model)
	{
		List<UserRow> userRows = userDAO.selectRowsByDisplayNameLike(pattern);
		return renderUserList(userRows, model);
	}

	private ModelAndView renderUserList(List<UserRow> userRows, Map<String, Object> model)
	{
		List<Map<String, String>> users = new ArrayList<>();
		for (UserRow userRow : userRows)
		{
			String createdTimestamp = TimeUtil.formatOffsetDateTime(userRow.getCreatedTimestamp());
			Map<String, String> thisUser = new HashMap<>();
			thisUser.put("id", String.valueOf(userRow.getId()));
			thisUser.put("name", userRow.getDisplayName());
			thisUser.put("created", createdTimestamp);
			users.add(thisUser);
		}
		model.put("users", users);
		return new ModelAndView("admin/users_list.html", model);
	}

	private ModelAndView renderSimpleMessage(String header, String message, HttpStatus status, Map<String, Object> model)
	{
		model.put("header", header);
		model.put("message", message);
		return new ModelAndView("admin/simple_message.html", model, status);
	}
}
