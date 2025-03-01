package lexer

import (
	"strconv"
	"strings"
	"unicode"
)

type HeadingMeta struct {
	Level int `json:"level"`
}
type FencedCodeBlockMeta struct {
	Language string `json:"language"`
}

type LinkMeta struct {
	Src string `json:"src"`
}

type ImageMeta struct {
	Src string `json:"src"`
}

type BlockquoteMeta struct {
	Depth int `json:"depth"`
}

type OrderedListMeta struct {
	Number int `json:"number"`
}

func (l *Lexer) enrichToken(token Token) Token {

	switch token.Type {
	case TokenEOF:
		return Token{Type: token.Type, Value: ""}
	case TokenBreak:
		return Token{Type: token.Type, Value: ""}

	case TokenHeading:
		level := countLeadingChars(token.Value, '#')
		return Token{Type: token.Type, Value: token.Value, CleanValue: trimLeadingChars(token.Value, '#', level), Meta: HeadingMeta{Level: level}}

	case TokenBold:
		boldChar := rune(token.Value[0])
		return Token{Type: token.Type, Value: token.Value, CleanValue: trimSurroundingChars(token.Value, boldChar, 2)}

	case TokenItalic:
		italicChar := rune(token.Value[0])
		return Token{Type: token.Type, Value: token.Value, CleanValue: trimSurroundingChars(token.Value, italicChar, 1)}

	case TokenInlineCode:
		return Token{Type: token.Type, Value: token.Value, CleanValue: trimSurroundingChars(token.Value, '`', 1)}

	case TokenFencedCodeBlock:
		language := parseLanguageFromFencedCodeBlock(token.Value)
		if len(language) > 0 {
			return Token{Type: token.Type, Value: token.Value, CleanValue: trimCodeBlock(token.Value, language), Meta: FencedCodeBlockMeta{Language: language}}
		}
		return Token{Type: token.Type, Value: token.Value, CleanValue: trimCodeBlock(token.Value, language)}

	case TokenUnorderedList:
		listChar := rune(token.Value[0])
		return Token{Type: token.Type, Value: token.Value, CleanValue: trimLeadingChars(token.Value, listChar, 1)}

	case TokenOrderedList:
		number, cleanValue := parseOrderedListParts(token.Value)

		return Token{Type: token.Type, Value: token.Value, CleanValue: cleanValue, Meta: OrderedListMeta{Number: number}}

	case TokenLink:
		title, url := parseLink(token.Value)
		return Token{Type: token.Type, Value: token.Value, CleanValue: title, Meta: LinkMeta{Src: url}}

	case TokenImage:
		alt, src := parseImage(token.Value)
		return Token{Type: token.Type, Value: token.Value, CleanValue: alt, Meta: ImageMeta{Src: src}}

	case TokenBlockquote:
		depth := countLeadingChars(token.Value, '>')
		token.CleanValue = trimLeadingChars(token.Value, '>', depth)
		return Token{Type: token.Type, Value: token.Value, CleanValue: trimLeadingChars(token.Value, '>', depth), Meta: BlockquoteMeta{Depth: depth}}
	}

	return Token{Type: token.Type, Value: token.Value, CleanValue: token.Value}
}

func countLeadingChars(s string, char rune) int {
	count := 0
	for _, c := range s {
		if c == char {
			count++
		} else {
			break
		}
	}
	return count
}

func trimLeadingChars(s string, char rune, count int) string {
	i := 0
	for i < len(s) && i < count && rune(s[i]) == char {
		i++
	}

	return strings.TrimFunc(s[i:], func(r rune) bool {
		return unicode.IsSpace(r) && r != '\n'
	})
}

func trimEndingChars(s string, char rune, count int) string {
	i := len(s) - 1
	for i >= 0 && i >= len(s)-count && rune(s[i]) == char {
		i--
	}

	return s[:i+1]

}

func trimSurroundingChars(s string, char rune, count int) string {
	s = trimLeadingChars(s, char, count)
	s = trimEndingChars(s, char, count)
	return s
}

func trimCodeBlock(s string, language string) string {
	s = trimSurroundingChars(s, '`', 3)

	if len(language) <= 0 {
		return s
	}

	return s[len(language):]
}

func parseLanguageFromFencedCodeBlock(s string) string {
	if s[0:4] == "```\n" {
		return ""
	}

	codeBlock := trimLeadingChars(s, '`', 3)

	words := strings.FieldsFunc(codeBlock, func(r rune) bool {
		return r == ' ' || r == '\n'
	})

	return words[0]
}

func parseOrderedListParts(s string) (int, string) {
	parts := strings.Split(s, ".")
	number, err := strconv.Atoi(parts[0])
	if err != nil {
		panic(err)
	}

	return number, strings.Trim(parts[1], " ")
}

func parseLink(s string) (string, string) {
	startText := strings.Index(s, "[") + 1
	endText := strings.Index(s, "]")
	startURL := strings.Index(s, "(") + 1
	endURL := strings.Index(s, ")")

	text := s[startText:endText]
	url := s[startURL:endURL]
	return text, url
}

func parseImage(s string) (string, string) {
	startAlt := strings.Index(s, "![") + 2
	endAlt := strings.Index(s, "]")
	startURL := strings.Index(s, "(") + 1
	endURL := strings.Index(s, ")")

	alt := s[startAlt:endAlt]
	src := s[startURL:endURL]
	return alt, src

}
