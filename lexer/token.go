package lexer

import (
	"encoding/json"
	"fmt"
)

type TokenType int

const (
	TokenEOF TokenType = iota
	TokenBreak
	TokenText
	TokenHeading
	TokenBold
	TokenItalic
	TokenInlineCode
	TokenFencedCodeBlock
	TokenUnorderedList
	TokenOrderedList
	TokenLink
	TokenImage
	TokenBlockquote
	TokenHorizontalRule
)

var TokenTypeStrings = map[TokenType]string{
	TokenEOF:             "EOF",
	TokenBreak:           "Break",
	TokenText:            "Text",
	TokenHeading:         "Heading",
	TokenBold:            "Bold",
	TokenItalic:          "Italic",
	TokenInlineCode:      "InlineCode",
	TokenFencedCodeBlock: "FencedCodeBlock",
	TokenUnorderedList:   "UnorderedList",
	TokenOrderedList:     "OrderedList",
	TokenLink:            "Link",
	TokenImage:           "Image",
	TokenBlockquote:      "Blockquote",
	TokenHorizontalRule:  "HorizontalRule",
}

type TokenMeta any

type Token struct {
	Type       TokenType `json:"type"`
	Value      string    `json:"rawValue"`
	CleanValue string    `json:"cleanValue"`
	Meta       any       `json:"meta"`
}

func (t Token) String() string {
	return fmt.Sprintf("{ Type: %s, Value: %q }\n", TokenTypeStrings[t.Type], t.Value)
}

func (t TokenType) MarshalJSON() ([]byte, error) {
	str, exists := TokenTypeStrings[t]
	if !exists {
		return nil, fmt.Errorf("invalid TokenType: %d", t)
	}
	return json.Marshal(str)
}
