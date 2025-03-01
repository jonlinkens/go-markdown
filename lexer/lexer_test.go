package lexer

import (
	"encoding/json"
	"testing"
)

type testCase struct {
	name     string
	input    string
	expected []Token
}

func runTests(t *testing.T, tests []testCase) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			go l.Run()

			for _, expectedToken := range tt.expected {
				actualToken := <-l.GetTokens()
				if actualToken != expectedToken {
					actualJson, _ := json.MarshalIndent(actualToken, "", "    ")
					expectedJson, _ := json.MarshalIndent(expectedToken, "", "    ")
					t.Errorf("%#v : %#v", actualToken, expectedToken)
					t.Errorf("tokens do not match:\n got %s\n expected %s", actualJson, expectedJson)
				}
			}
		})
	}
}

func TestLexerPlainText(t *testing.T) {
	tests := []testCase{
		{
			name:  "simple text",
			input: "This is plain text",
			expected: []Token{
				{Type: TokenText, Value: "This is plain text", CleanValue: "This is plain text"},
				{Type: TokenEOF, Value: ""},
			},
		},

		{
			name:  "simple text with no break",
			input: "This is plain text\n",
			expected: []Token{
				{Type: TokenText, Value: "This is plain text", CleanValue: "This is plain text"},
				{Type: TokenEOF, Value: ""},
			},
		},

		{
			name:  "simple text with break",
			input: "This is plain text\n\n",
			expected: []Token{
				{Type: TokenText, Value: "This is plain text", CleanValue: "This is plain text"},
				{Type: TokenBreak, Value: ""},
				{Type: TokenEOF, Value: ""},
			},
		},
	}
	runTests(t, tests)
}

func TestLexerHeadings(t *testing.T) {
	tests := []testCase{
		{
			name:  "single level heading",
			input: "# Heading",
			expected: []Token{
				{Type: TokenHeading, Value: "# Heading", CleanValue: "Heading", Meta: HeadingMeta{Level: 1}},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "multi level heading",
			input: "### Deep Heading",
			expected: []Token{
				{Type: TokenHeading, Value: "### Deep Heading", CleanValue: "Deep Heading", Meta: HeadingMeta{Level: 3}},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "heading with text after",
			input: "# Heading\nNormal text",
			expected: []Token{
				{Type: TokenHeading, Value: "# Heading", CleanValue: "Heading", Meta: HeadingMeta{Level: 1}},
				{Type: TokenText, Value: "Normal text", CleanValue: "Normal text"},
				{Type: TokenEOF, Value: ""},
			},
		},
	}
	runTests(t, tests)
}

func TestLexerBold(t *testing.T) {
	tests := []testCase{
		{
			name:  "asterisk bold",
			input: "**bold text**",
			expected: []Token{
				{Type: TokenBold, Value: "**bold text**", CleanValue: "bold text"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "underscore bold",
			input: "__bold text__",
			expected: []Token{
				{Type: TokenBold, Value: "__bold text__", CleanValue: "bold text"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "bold with surrounding text",
			input: "before **bold** after",
			expected: []Token{
				{Type: TokenText, Value: "before ", CleanValue: "before "},
				{Type: TokenBold, Value: "**bold**", CleanValue: "bold"},
				{Type: TokenText, Value: " after", CleanValue: " after"},
				{Type: TokenEOF, Value: ""},
			},
		},
	}
	runTests(t, tests)
}

func TestLexerItalic(t *testing.T) {
	tests := []testCase{
		{
			name:  "asterisk italic",
			input: "*italic text*",
			expected: []Token{
				{Type: TokenItalic, Value: "*italic text*", CleanValue: "italic text"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "underscore italic",
			input: "_italic text_",
			expected: []Token{
				{Type: TokenItalic, Value: "_italic text_", CleanValue: "italic text"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "italic with surrounding text",
			input: "before *italic* after",
			expected: []Token{
				{Type: TokenText, Value: "before ", CleanValue: "before "},
				{Type: TokenItalic, Value: "*italic*", CleanValue: "italic"},
				{Type: TokenText, Value: " after", CleanValue: " after"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "unclosed italic becomes text",
			input: "*unclosed italic",
			expected: []Token{
				{Type: TokenText, Value: "*unclosed italic", CleanValue: "*unclosed italic"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "multiple italics in text",
			input: "*first* normal *second*",
			expected: []Token{
				{Type: TokenItalic, Value: "*first*", CleanValue: "first"},
				{Type: TokenText, Value: " normal ", CleanValue: " normal "},
				{Type: TokenItalic, Value: "*second*", CleanValue: "second"},
				{Type: TokenEOF, Value: ""},
			},
		},
	}
	runTests(t, tests)
}

func TestLexerInlineCode(t *testing.T) {
	tests := []testCase{
		{
			name:  "basic inline code",
			input: "`code`",
			expected: []Token{
				{Type: TokenInlineCode, Value: "`code`", CleanValue: "code"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "inline code with spaces",
			input: "`code with spaces`",
			expected: []Token{
				{Type: TokenInlineCode, Value: "`code with spaces`", CleanValue: "code with spaces"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "inline code with surrounding text",
			input: "before `code` after",
			expected: []Token{
				{Type: TokenText, Value: "before ", CleanValue: "before "},
				{Type: TokenInlineCode, Value: "`code`", CleanValue: "code"},
				{Type: TokenText, Value: " after", CleanValue: " after"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "unclosed inline code becomes text",
			input: "`unclosed code",
			expected: []Token{
				{Type: TokenText, Value: "`unclosed code", CleanValue: "`unclosed code"},
				{Type: TokenEOF, Value: ""},
			},
		},
	}
	runTests(t, tests)
}

func TestLexerFencedCodeBlock(t *testing.T) {
	tests := []testCase{
		{
			name:  "basic code block",
			input: "```\ncode block\n```",
			expected: []Token{
				{Type: TokenFencedCodeBlock, Value: "```\ncode block\n```", CleanValue: "\ncode block\n"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "code block with language",
			input: "```python\ndef hello():\n    print('Hello')\n```",
			expected: []Token{
				{Type: TokenFencedCodeBlock, Value: "```python\ndef hello():\n    print('Hello')\n```", CleanValue: "\ndef hello():\n    print('Hello')\n", Meta: FencedCodeBlockMeta{Language: "python"}},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "code block with surrounding text",
			input: "before\n```\ncode\n```\nafter",
			expected: []Token{
				{Type: TokenText, Value: "before", CleanValue: "before"},
				{Type: TokenFencedCodeBlock, Value: "```\ncode\n```", CleanValue: "\ncode\n"},
				{Type: TokenText, Value: "after", CleanValue: "after"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "unclosed code block becomes text",
			input: "```\nunclosed code block",
			expected: []Token{
				{Type: TokenText, Value: "```\nunclosed code block", CleanValue: "```\nunclosed code block"},
				{Type: TokenEOF, Value: ""},
			},
		},
	}
	runTests(t, tests)
}

func TestLexerUnorderedList(t *testing.T) {
	tests := []testCase{
		{
			name:  "dash list item",
			input: "- list item",
			expected: []Token{
				{Type: TokenUnorderedList, Value: "- list item", CleanValue: "list item"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "plus list item",
			input: "+ list item",
			expected: []Token{
				{Type: TokenUnorderedList, Value: "+ list item", CleanValue: "list item"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "asterisk list item",
			input: "* list item",
			expected: []Token{
				{Type: TokenUnorderedList, Value: "* list item", CleanValue: "list item"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "multiple list items",
			input: "- first\n+ second\n* third",
			expected: []Token{
				{Type: TokenUnorderedList, Value: "- first", CleanValue: "first"},

				{Type: TokenUnorderedList, Value: "+ second", CleanValue: "second"},
				{Type: TokenUnorderedList, Value: "* third", CleanValue: "third"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "not a list item",
			input: "-not a list",
			expected: []Token{
				{Type: TokenText, Value: "-not a list", CleanValue: "-not a list"},
				{Type: TokenEOF, Value: ""},
			},
		},
	}
	runTests(t, tests)
}

func TestLexerBlockquote(t *testing.T) {
	tests := []testCase{
		{
			name:  "simple blockquote",
			input: "> quoted text",
			expected: []Token{
				{Type: TokenBlockquote, Value: "> quoted text", CleanValue: "quoted text", Meta: BlockquoteMeta{Depth: 1}},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "blockquote with surrounding text",
			input: "before\n> quoted text\nafter",
			expected: []Token{
				{Type: TokenText, Value: "before", CleanValue: "before"},
				{Type: TokenBlockquote, Value: "> quoted text", CleanValue: "quoted text", Meta: BlockquoteMeta{Depth: 1}},
				{Type: TokenText, Value: "after", CleanValue: "after"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "multiple blockquotes",
			input: "> first quote\n> second quote",
			expected: []Token{
				{Type: TokenBlockquote, Value: "> first quote", CleanValue: "first quote", Meta: BlockquoteMeta{Depth: 1}},

				{Type: TokenBlockquote, Value: "> second quote", CleanValue: "second quote", Meta: BlockquoteMeta{Depth: 1}},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "not a blockquote",
			input: "a>not a quote",
			expected: []Token{
				{Type: TokenText, Value: "a>not a quote", CleanValue: "a>not a quote"},
				{Type: TokenEOF, Value: ""},
			},
		},
	}
	runTests(t, tests)
}

func TestLexerHorizontalRule(t *testing.T) {
	tests := []testCase{
		{
			name:  "hyphen rule",
			input: "---",
			expected: []Token{
				{Type: TokenHorizontalRule, Value: "---", CleanValue: "---"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "asterisk rule",
			input: "***",
			expected: []Token{
				{Type: TokenHorizontalRule, Value: "***", CleanValue: "***"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "underscore rule",
			input: "___",
			expected: []Token{
				{Type: TokenHorizontalRule, Value: "___", CleanValue: "___"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "longer rule",
			input: "-----",
			expected: []Token{
				{Type: TokenHorizontalRule, Value: "-----", CleanValue: "-----"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "rule with surrounding text",
			input: "before\n---\nafter",
			expected: []Token{
				{Type: TokenText, Value: "before", CleanValue: "before"},
				{Type: TokenHorizontalRule, Value: "---", CleanValue: "---"},
				{Type: TokenText, Value: "after", CleanValue: "after"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "not a rule",
			input: "--",
			expected: []Token{
				{Type: TokenText, Value: "--", CleanValue: "--"},
				{Type: TokenEOF, Value: ""},
			},
		},
	}
	runTests(t, tests)
}

func TestLexerOrderedList(t *testing.T) {
	tests := []testCase{
		{
			name:  "simple ordered list",
			input: "1. list item",
			expected: []Token{
				{Type: TokenOrderedList, Value: "1. list item", CleanValue: "list item", Meta: OrderedListMeta{Number: 1}},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "multiple ordered items",
			input: "1. first\n2. second\n3. third",
			expected: []Token{
				{Type: TokenOrderedList, Value: "1. first", CleanValue: "first", Meta: OrderedListMeta{Number: 1}},
				{Type: TokenOrderedList, Value: "2. second", CleanValue: "second", Meta: OrderedListMeta{Number: 2}},
				{Type: TokenOrderedList, Value: "3. third", CleanValue: "third", Meta: OrderedListMeta{Number: 3}},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "not an ordered list",
			input: "1 not a list",
			expected: []Token{
				{Type: TokenText, Value: "1 not a list", CleanValue: "1 not a list"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "ordered list with surrounding text",
			input: "before\n1. list item\nafter",
			expected: []Token{
				{Type: TokenText, Value: "before", CleanValue: "before"},
				{Type: TokenOrderedList, Value: "1. list item", CleanValue: "list item", Meta: OrderedListMeta{Number: 1}},
				{Type: TokenText, Value: "after", CleanValue: "after"},
				{Type: TokenEOF, Value: ""},
			},
		},
	}
	runTests(t, tests)
}

func TestLexerLink(t *testing.T) {
	tests := []testCase{
		{
			name:  "simple link",
			input: "[text](url)",
			expected: []Token{
				{Type: TokenLink, Value: "[text](url)", CleanValue: "text", Meta: LinkMeta{Src: "url"}},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "link with surrounding text",
			input: "before [text](url) after",
			expected: []Token{
				{Type: TokenText, Value: "before ", CleanValue: "before "},
				{Type: TokenLink, Value: "[text](url)", CleanValue: "text", Meta: LinkMeta{Src: "url"}},
				{Type: TokenText, Value: " after", CleanValue: " after"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "multiple links",
			input: "[first](url1) [second](url2)",
			expected: []Token{
				{Type: TokenLink, Value: "[first](url1)", CleanValue: "first", Meta: LinkMeta{Src: "url1"}},
				{Type: TokenText, Value: " ", CleanValue: " "},
				{Type: TokenLink, Value: "[second](url2)", CleanValue: "second", Meta: LinkMeta{Src: "url2"}},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "unclosed link becomes text",
			input: "[text(url)",
			expected: []Token{
				{Type: TokenText, Value: "[text(url)", CleanValue: "[text(url)"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "link without url becomes text",
			input: "[text]",
			expected: []Token{
				{Type: TokenText, Value: "[text]", CleanValue: "[text]"},
				{Type: TokenEOF, Value: ""},
			},
		},
	}
	runTests(t, tests)
}

func TestLexerImage(t *testing.T) {
	tests := []testCase{
		{
			name:  "simple image",
			input: "![alt](src)",
			expected: []Token{
				{Type: TokenImage, Value: "![alt](src)", CleanValue: "alt", Meta: ImageMeta{Src: "src"}},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "image with surrounding text",
			input: "before ![alt](src) after",
			expected: []Token{
				{Type: TokenText, Value: "before ", CleanValue: "before "},
				{Type: TokenImage, Value: "![alt](src)", CleanValue: "alt", Meta: ImageMeta{Src: "src"}},
				{Type: TokenText, Value: " after", CleanValue: " after"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "multiple images",
			input: "![first](src1) ![second](src2)",
			expected: []Token{
				{Type: TokenImage, Value: "![first](src1)", CleanValue: "first", Meta: ImageMeta{Src: "src1"}},
				{Type: TokenText, Value: " ", CleanValue: " "},
				{Type: TokenImage, Value: "![second](src2)", CleanValue: "second", Meta: ImageMeta{Src: "src2"}},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "unclosed image becomes text",
			input: "![alt(src)",
			expected: []Token{
				{Type: TokenText, Value: "![alt(src)", CleanValue: "![alt(src)"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "image without src becomes text",
			input: "![alt]",
			expected: []Token{
				{Type: TokenText, Value: "![alt]", CleanValue: "![alt]"},
				{Type: TokenEOF, Value: ""},
			},
		},
	}
	runTests(t, tests)
}
