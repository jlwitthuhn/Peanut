// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package io.github.jlwitthuhn.peanut.util;

import java.util.concurrent.TimeUnit;

public class TimeUtil
{
	public static String formatMillisecondsAsDDHHMMSS(long milliseconds)
	{
		long day = TimeUnit.MILLISECONDS.toDays(milliseconds);
		long hour = TimeUnit.MILLISECONDS.toHours(milliseconds) % 24;
		long minute = TimeUnit.MILLISECONDS.toMinutes(milliseconds) % 60;
		long second = TimeUnit.MILLISECONDS.toSeconds(milliseconds) % 60;
		return String.format("%02d:%02d:%02d:%02d", day, hour, minute, second);
	}
}
