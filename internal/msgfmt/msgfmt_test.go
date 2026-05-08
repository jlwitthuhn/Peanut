// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package msgfmt

import "testing"

func TestFormatMessage_Blockquotes(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "simple quote",
			input: "> hello",
			want:  "<blockquote>hello</blockquote>",
		},
		{
			name:  "grouped multi-line quote",
			input: "> line one\n> line two",
			want:  "<blockquote>line one\nline two</blockquote>",
		},
		{
			name:  "nested quote",
			input: ">> inner",
			want:  "<blockquote><blockquote>inner</blockquote></blockquote>",
		},
		{
			name:  "nested then outer",
			input: ">> inner\n> outer",
			want:  "<blockquote><blockquote>inner</blockquote>\nouter</blockquote>",
		},
		{
			name:  "text before and after quote",
			input: "before\n> quoted\nafter",
			want:  "before\n<blockquote>quoted</blockquote>\nafter",
		},
		{
			name:  "inline formatting inside quote",
			input: "> **bold** and *italic*",
			want:  "<blockquote><b>bold</b> and <i>italic</i></blockquote>",
		},
		{
			name:  "heading inside quote",
			input: "> # Title",
			want:  "<blockquote><h1>Title</h1></blockquote>",
		},
		{
			name:  "link inside quote",
			input: "> [click](https://example.com)",
			want:  `<blockquote><a href="https://example.com" rel="nofollow">click</a></blockquote>`,
		},
		{
			name:  "triple nested",
			input: ">>> deep",
			want:  "<blockquote><blockquote><blockquote>deep</blockquote></blockquote></blockquote>",
		},
		{
			name:  "html escaped in quote",
			input: "> <script>alert('xss')</script>",
			want:  "<blockquote>&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;</blockquote>",
		},
		{
			name:  "no space after arrow",
			input: ">nospace",
			want:  "<blockquote>nospace</blockquote>",
		},
		{
			name:  "empty quote",
			input: ">",
			want:  "<blockquote></blockquote>",
		},
		{
			name:  "mixed depths",
			input: "> a\n>> b\n>> c\n> d",
			want:  "<blockquote>a\n<blockquote>b\nc</blockquote>\nd</blockquote>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Format(tt.input)
			if got != tt.want {
				t.Errorf("\ninput: %q\n  got: %q\n want: %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestFormatMessage_PerLineTags(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "bold",
			input: "**hello**",
			want:  "<b>hello</b>",
		},
		{
			name:  "italic",
			input: "*hello*",
			want:  "<i>hello</i>",
		},
		{
			name:  "strikethrough",
			input: "~~hello~~",
			want:  "<s>hello</s>",
		},
		{
			name:  "heading",
			input: "## Title",
			want:  "<h2>Title</h2>",
		},
		{
			name:  "link",
			input: "[text](https://example.com)",
			want:  `<a href="https://example.com" rel="nofollow">text</a>`,
		},
		{
			name:  "image",
			input: "![alt](https://example.com/img.png)",
			want:  `<img src="https://example.com/img.png" alt="alt">`,
		},
		{
			name:  "multi-line",
			input: "line one\nline two\nline three",
			want:  "line one\nline two\nline three",
		},
		{
			name:  "html escaping",
			input: "<b>not bold</b>",
			want:  "&lt;b&gt;not bold&lt;/b&gt;",
		},
		{
			name:  "script injection",
			input: "<script>alert('xss')</script>",
			want:  "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Format(tt.input)
			if got != tt.want {
				t.Errorf("\ninput: %q\n  got: %q\n want: %q", tt.input, got, tt.want)
			}
		})
	}
}
