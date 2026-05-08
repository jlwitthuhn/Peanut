// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package msgfmt

import (
	"fmt"
	"html"
	"regexp"
	"strings"
)

func Format(message string) string {
	escaped := html.EscapeString(message)
	return formatBlocks(strings.Split(escaped, "\n"))
}

// formatBlocks groups consecutive lines by quote depth and renders them.
// Empty normal lines are dropped so the output never contains <p></p>.
func formatBlocks(lines []string) string {
	var blocks []string

	i := 0
	for i < len(lines) {
		depth, _ := quotePrefix(lines[i])

		if depth == 0 {
			if rendered := formatLine(lines[i]); rendered != "" {
				blocks = append(blocks, rendered)
			}
			i++
			continue
		}

		// Collect a run of consecutive lines that all have depth >= 1.
		start := i
		for i < len(lines) {
			d, _ := quotePrefix(lines[i])
			if d == 0 {
				break
			}
			i++
		}

		// Strip one level of '>' from each line in the run and recurse.
		inner := make([]string, i-start)
		for j := start; j < i; j++ {
			_, rest := stripOneLevel(lines[j])
			inner[j-start] = rest
		}

		blocks = append(blocks, "<blockquote>"+formatBlocks(inner)+"</blockquote>")
	}

	return strings.Join(blocks, "\n")
}

// quotePrefix returns how many leading `&gt;` markers the line has and the
// remaining text after stripping them.
func quotePrefix(line string) (depth int, rest string) {
	const marker = "&gt;"
	for strings.HasPrefix(line, marker) {
		depth++
		line = line[len(marker):]
		if len(line) > 0 && line[0] == ' ' {
			line = line[1:]
		}
	}
	return depth, line
}

// stripOneLevel removes exactly one `&gt;` prefix (plus optional space).
func stripOneLevel(line string) (ok bool, rest string) {
	const marker = "&gt;"
	if !strings.HasPrefix(line, marker) {
		return false, line
	}
	line = line[len(marker):]
	if len(line) > 0 && line[0] == ' ' {
		line = line[1:]
	}
	return true, line
}

// formatLine handles a single already-escaped line.
func formatLine(line string) string {
	if level, content, ok := parseHeading(line); ok {
		return fmt.Sprintf("<h%d>%s</h%d>", level, formatInline(content), level)
	}
	if strings.TrimSpace(line) == "" {
		return ""
	}
	return "<p>" + formatInline(line) + "</p>"
}

// parseHeading recognises a markdown `# Title`.
func parseHeading(line string) (level int, content string, ok bool) {
	const maxLevel = 6
	count := 0
	for count < len(line) && count < maxLevel && line[count] == '#' {
		count++
	}
	if count == 0 || count >= len(line) || line[count] != ' ' {
		return 0, "", false
	}
	return count, line[count+1:], true
}

// formatInline applies all inline transforms to a single line of escaped text.
func formatInline(text string) string {
	text = applyImages(text)
	text = applyLinks(text)
	text = applyBold(text)
	text = applyItalic(text)
	text = applyStrikethrough(text)
	return text
}

var imagePattern = regexp.MustCompile(`!\[([^\]]+)\]\(([^)]+)\)`)

// applyImages rewrites ![alt](url) into <img src="url" alt="alt">.
func applyImages(text string) string {
	return imagePattern.ReplaceAllStringFunc(text, func(match string) string {
		parts := imagePattern.FindStringSubmatch(match)
		alt, url := parts[1], parts[2]

		if !isSafeURL(url) {
			return match
		}

		var b strings.Builder
		b.WriteString(`<img src="`)
		b.WriteString(url)
		b.WriteString(`" alt="`)
		b.WriteString(alt)
		b.WriteString(`">`)
		return b.String()
	})
}

// linkPattern matches [visible text](url).
var linkPattern = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)

// applyLinks rewrites [text](url) into <a href="url" rel="nofollow">text</a>..
func applyLinks(text string) string {
	return linkPattern.ReplaceAllStringFunc(text, func(match string) string {
		parts := linkPattern.FindStringSubmatch(match)
		linkText, url := parts[1], parts[2]

		if !isSafeURL(url) {
			return match
		}

		var b strings.Builder
		b.WriteString(`<a href="`)
		b.WriteString(url)
		b.WriteString(`" rel="nofollow">`)
		b.WriteString(linkText)
		b.WriteString(`</a>`)
		return b.String()
	})
}

// isSafeURL determines what URLs are allowed to be used in links
func isSafeURL(url string) bool {
	trimmed := strings.TrimSpace(url)
	lower := strings.ToLower(trimmed)

	switch {
	case strings.HasPrefix(lower, "http://"),
		strings.HasPrefix(lower, "https://"):
		return true
	}

	colon := strings.Index(lower, ":")
	if colon == -1 {
		return true
	}
	slash := strings.Index(lower, "/")
	return slash != -1 && slash < colon
}

var boldPattern = regexp.MustCompile(`\*\*([^*]+)\*\*`)

// applyBold implements **bold text**.
func applyBold(text string) string {
	return boldPattern.ReplaceAllString(text, "<b>$1</b>")
}

var italicPattern = regexp.MustCompile(`\*([^*]+)\*`)

// applyItalic implements *italic text*. Run after applyBold so the double
// asterisks of bold spans have already been consumed.
func applyItalic(text string) string {
	return italicPattern.ReplaceAllString(text, "<i>$1</i>")
}

var strikethroughPattern = regexp.MustCompile(`~~([^~]+)~~`)

// applyStrikethrough implements ~~struck text~~.
func applyStrikethrough(text string) string {
	return strikethroughPattern.ReplaceAllString(text, "<s>$1</s>")
}
