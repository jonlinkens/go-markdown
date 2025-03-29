package parser

import (
	"testing"

	"github.com/jonlinkens/go-markdown/lexer"
)

func TestParserBasic(t *testing.T) {
	tests := []struct {
		name     string
		tokens   []lexer.Token
		expected *Document
	}{
		{
			name: "simple text",
			tokens: []lexer.Token{
				{Type: lexer.TokenText, Value: "Hello world", CleanValue: "Hello world"},
				{Type: lexer.TokenEOF},
			},
			expected: &Document{
				Node: &Node{
					Type: NodeDocument,
					Children: []*Node{
						{
							Type: NodeParagraph,
							Children: []*Node{
								{Type: NodeText, Value: "Hello world"},
							},
						},
					},
				},
			},
		},
		{
			name: "heading with text",
			tokens: []lexer.Token{
				{Type: lexer.TokenHeading, Value: "# Title", CleanValue: "Title", Meta: lexer.HeadingMeta{Level: 1}},
				{Type: lexer.TokenText, Value: "Hello world", CleanValue: "Hello world"},
				{Type: lexer.TokenEOF},
			},
			expected: &Document{
				Node: &Node{
					Type: NodeDocument,
					Children: []*Node{
						{
							Type:  NodeHeading,
							Value: "Title",
							Meta:  lexer.HeadingMeta{Level: 1},
						},
						{
							Type: NodeParagraph,
							Children: []*Node{
								{Type: NodeText, Value: "Hello world"},
							},
						},
					},
				},
			},
		},
		{
			name: "text with inline formatting",
			tokens: []lexer.Token{
				{Type: lexer.TokenText, Value: "Hello ", CleanValue: "Hello "},
				{Type: lexer.TokenBold, Value: "**bold**", CleanValue: "bold"},
				{Type: lexer.TokenText, Value: " and ", CleanValue: " and "},
				{Type: lexer.TokenItalic, Value: "*italic*", CleanValue: "italic"},
				{Type: lexer.TokenEOF},
			},
			expected: &Document{
				Node: &Node{
					Type: NodeDocument,
					Children: []*Node{
						{
							Type: NodeParagraph,
							Children: []*Node{
								{Type: NodeText, Value: "Hello "},
								{Type: NodeBold, Value: "bold"},
								{Type: NodeText, Value: " and "},
								{Type: NodeItalic, Value: "italic"},
							},
						},
					},
				},
			},
		},
		{
			name: "unordered list",
			tokens: []lexer.Token{
				{Type: lexer.TokenUnorderedList, Value: "- First", CleanValue: "First"},
				{Type: lexer.TokenUnorderedList, Value: "- Second", CleanValue: "Second"},
				{Type: lexer.TokenEOF},
			},
			expected: &Document{
				Node: &Node{
					Type: NodeDocument,
					Children: []*Node{
						{
							Type: NodeUnorderedList,
							Children: []*Node{
								{Type: NodeListItem, Value: "First"},
								{Type: NodeListItem, Value: "Second"},
							},
						},
					},
				},
			},
		},
		{
			name: "ordered list",
			tokens: []lexer.Token{
				{Type: lexer.TokenOrderedList, Value: "1. First", CleanValue: "First", Meta: lexer.OrderedListMeta{Number: 1}},
				{Type: lexer.TokenOrderedList, Value: "2. Second", CleanValue: "Second", Meta: lexer.OrderedListMeta{Number: 2}},
				{Type: lexer.TokenEOF},
			},
			expected: &Document{
				Node: &Node{
					Type: NodeDocument,
					Children: []*Node{
						{
							Type: NodeOrderedList,
							Children: []*Node{
								{Type: NodeListItem, Value: "First", Meta: lexer.OrderedListMeta{Number: 1}},
								{Type: NodeListItem, Value: "Second", Meta: lexer.OrderedListMeta{Number: 2}},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.tokens)
			result := parser.Parse()
			compareNodes(t, result.Node, tt.expected.Node)
		})
	}
}

func TestParserInline(t *testing.T) {
	tests := []struct {
		name     string
		tokens   []lexer.Token
		expected *Document
	}{
		{
			name: "mixed inline formatting",
			tokens: []lexer.Token{
				{Type: lexer.TokenText, Value: "This is ", CleanValue: "This is "},
				{Type: lexer.TokenBold, Value: "**bold**", CleanValue: "bold"},
				{Type: lexer.TokenText, Value: " with ", CleanValue: " with "},
				{Type: lexer.TokenItalic, Value: "*italic*", CleanValue: "italic"},
				{Type: lexer.TokenText, Value: " and ", CleanValue: " and "},
				{Type: lexer.TokenInlineCode, Value: "`code`", CleanValue: "code"},
				{Type: lexer.TokenEOF},
			},
			expected: &Document{
				Node: &Node{
					Type: NodeDocument,
					Children: []*Node{
						{
							Type: NodeParagraph,
							Children: []*Node{
								{Type: NodeText, Value: "This is "},
								{Type: NodeBold, Value: "bold"},
								{Type: NodeText, Value: " with "},
								{Type: NodeItalic, Value: "italic"},
								{Type: NodeText, Value: " and "},
								{Type: NodeInlineCode, Value: "code"},
							},
						},
					},
				},
			},
		},
		{
			name: "inline elements with link",
			tokens: []lexer.Token{
				{Type: lexer.TokenText, Value: "Check ", CleanValue: "Check "},
				{Type: lexer.TokenBold, Value: "**this**", CleanValue: "this"},
				{Type: lexer.TokenText, Value: " ", CleanValue: " "},
				{Type: lexer.TokenLink, Value: "[link](https://example.com)", CleanValue: "link", Meta: lexer.LinkMeta{Src: "https://example.com"}},
				{Type: lexer.TokenText, Value: " out", CleanValue: " out"},
				{Type: lexer.TokenEOF},
			},
			expected: &Document{
				Node: &Node{
					Type: NodeDocument,
					Children: []*Node{
						{
							Type: NodeParagraph,
							Children: []*Node{
								{Type: NodeText, Value: "Check "},
								{Type: NodeBold, Value: "this"},
								{Type: NodeText, Value: " "},
								{Type: NodeLink, Value: "link", Meta: lexer.LinkMeta{Src: "https://example.com"}},
								{Type: NodeText, Value: " out"},
							},
						},
					},
				},
			},
		},
		{
			name: "multiple paragraphs with inline elements",
			tokens: []lexer.Token{
				{Type: lexer.TokenText, Value: "First ", CleanValue: "First "},
				{Type: lexer.TokenBold, Value: "**paragraph**", CleanValue: "paragraph"},
				{Type: lexer.TokenBreak},
				{Type: lexer.TokenText, Value: "Second with ", CleanValue: "Second with "},
				{Type: lexer.TokenItalic, Value: "*emphasis*", CleanValue: "emphasis"},
				{Type: lexer.TokenEOF},
			},
			expected: &Document{
				Node: &Node{
					Type: NodeDocument,
					Children: []*Node{
						{
							Type: NodeParagraph,
							Children: []*Node{
								{Type: NodeText, Value: "First "},
								{Type: NodeBold, Value: "paragraph"},
							},
						},
						{
							Type: NodeParagraph,
							Children: []*Node{
								{Type: NodeText, Value: "Second with "},
								{Type: NodeItalic, Value: "emphasis"},
							},
						},
					},
				},
			},
		},
		{
			name: "inline code with special characters",
			tokens: []lexer.Token{
				{Type: lexer.TokenText, Value: "Use ", CleanValue: "Use "},
				{Type: lexer.TokenInlineCode, Value: "`const x = 42;`", CleanValue: "const x = 42;"},
				{Type: lexer.TokenText, Value: " in your code", CleanValue: " in your code"},
				{Type: lexer.TokenEOF},
			},
			expected: &Document{
				Node: &Node{
					Type: NodeDocument,
					Children: []*Node{
						{
							Type: NodeParagraph,
							Children: []*Node{
								{Type: NodeText, Value: "Use "},
								{Type: NodeInlineCode, Value: "const x = 42;"},
								{Type: NodeText, Value: " in your code"},
							},
						},
					},
				},
			},
		},
		{
			name: "inline elements in list items",
			tokens: []lexer.Token{
				{Type: lexer.TokenUnorderedList, Value: "- Item with ", CleanValue: "Item with "},
				{Type: lexer.TokenBold, Value: "**bold**", CleanValue: "bold"},
				{Type: lexer.TokenText, Value: " text", CleanValue: " text"},
				{Type: lexer.TokenUnorderedList, Value: "- Another with ", CleanValue: "Another with "},
				{Type: lexer.TokenItalic, Value: "*italic*", CleanValue: "italic"},
				{Type: lexer.TokenEOF},
			},
			expected: &Document{
				Node: &Node{
					Type: NodeDocument,
					Children: []*Node{
						{
							Type: NodeUnorderedList,
							Children: []*Node{
								{
									Type:  NodeListItem,
									Value: "Item with bold text",
									Children: []*Node{
										{Type: NodeText, Value: "Item with "},
										{Type: NodeBold, Value: "bold"},
										{Type: NodeText, Value: " text"},
									},
								},
								{
									Type:  NodeListItem,
									Value: "Another with italic",
									Children: []*Node{
										{Type: NodeText, Value: "Another with "},
										{Type: NodeItalic, Value: "italic"},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.tokens)
			result := parser.Parse()
			compareNodes(t, result.Node, tt.expected.Node)
		})
	}
}

func compareNodes(t *testing.T, got, want *Node) {
	if got.Type != want.Type {
		t.Errorf("node type mismatch: got %v, want %v", got.Type, want.Type)
		return
	}

	if got.Value != want.Value {
		t.Errorf("node value mismatch: got %q, want %q", got.Value, want.Value)
		return
	}

	if len(got.Children) != len(want.Children) {
		t.Errorf("children length mismatch: got %d, want %d", len(got.Children), len(want.Children))
		return
	}

	for i := range got.Children {
		compareNodes(t, got.Children[i], want.Children[i])
	}
}
