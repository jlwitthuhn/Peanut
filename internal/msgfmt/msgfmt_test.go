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
			want:  "<blockquote><p>hello</p></blockquote>",
		},
		{
			name:  "grouped multi-line quote",
			input: "> line one\n> line two",
			want:  "<blockquote><p>line one</p>\n<p>line two</p></blockquote>",
		},
		{
			name:  "nested quote",
			input: ">> inner",
			want:  "<blockquote><blockquote><p>inner</p></blockquote></blockquote>",
		},
		{
			name:  "nested then outer",
			input: ">> inner\n> outer",
			want:  "<blockquote><blockquote><p>inner</p></blockquote>\n<p>outer</p></blockquote>",
		},
		{
			name:  "text before and after quote",
			input: "before\n> quoted\nafter",
			want:  "<p>before</p>\n<blockquote><p>quoted</p></blockquote>\n<p>after</p>",
		},
		{
			name:  "inline formatting inside quote",
			input: "> **bold** and *italic*",
			want:  "<blockquote><p><b>bold</b> and <i>italic</i></p></blockquote>",
		},
		{
			name:  "heading inside quote",
			input: "> # Title",
			want:  "<blockquote><h1>Title</h1></blockquote>",
		},
		{
			name:  "link inside quote",
			input: "> [click](https://example.com)",
			want:  `<blockquote><p><a href="https://example.com" rel="nofollow">click</a></p></blockquote>`,
		},
		{
			name:  "triple nested",
			input: ">>> deep",
			want:  "<blockquote><blockquote><blockquote><p>deep</p></blockquote></blockquote></blockquote>",
		},
		{
			name:  "html escaped in quote",
			input: "> <script>alert('xss')</script>",
			want:  "<blockquote><p>&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;</p></blockquote>",
		},
		{
			name:  "no space after arrow",
			input: ">nospace",
			want:  "<blockquote><p>nospace</p></blockquote>",
		},
		{
			name:  "empty quote",
			input: ">",
			want:  "<blockquote></blockquote>",
		},
		{
			name:  "mixed depths",
			input: "> a\n>> b\n>> c\n> d",
			want:  "<blockquote><p>a</p>\n<blockquote><p>b</p>\n<p>c</p></blockquote>\n<p>d</p></blockquote>",
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
			want:  "<p><b>hello</b></p>",
		},
		{
			name:  "italic",
			input: "*hello*",
			want:  "<p><i>hello</i></p>",
		},
		{
			name:  "strikethrough",
			input: "~~hello~~",
			want:  "<p><s>hello</s></p>",
		},
		{
			name:  "heading",
			input: "## Title",
			want:  "<h2>Title</h2>",
		},
		{
			name:  "link",
			input: "[text](https://example.com)",
			want:  `<p><a href="https://example.com" rel="nofollow">text</a></p>`,
		},
		{
			name:  "image",
			input: "![alt](https://example.com/img.png)",
			want:  `<p><img src="https://example.com/img.png" alt="alt"></p>`,
		},
		{
			name:  "multi-line",
			input: "line one\nline two\nline three",
			want:  "<p>line one</p>\n<p>line two</p>\n<p>line three</p>",
		},
		{
			name:  "blank lines dropped",
			input: "line one\n\nline two",
			want:  "<p>line one</p>\n<p>line two</p>",
		},
		{
			name:  "only blank lines",
			input: "\n\n",
			want:  "",
		},
		{
			name:  "heading then paragraph",
			input: "# Title\nbody",
			want:  "<h1>Title</h1>\n<p>body</p>",
		},
		{
			name:  "html escaping",
			input: "<b>not bold</b>",
			want:  "<p>&lt;b&gt;not bold&lt;/b&gt;</p>",
		},
		{
			name:  "script injection",
			input: "<script>alert('xss')</script>",
			want:  "<p>&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;</p>",
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
