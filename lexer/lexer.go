package lexer

import (
	"strings"
	"unicode"
)

type Lexer struct {
	input  string
	start  int
	pos    int
	tokens chan Token
}

type stateFn func(*Lexer) stateFn

const eof = -1

func (l *Lexer) next() rune {
	if l.pos >= len(l.input) {
		return eof
	}
	r := rune(l.input[l.pos])
	l.pos++
	return r
}

func (l *Lexer) peek() rune {
	if l.pos >= len(l.input) {
		return eof
	}
	return rune(l.input[l.pos])
}

func (l *Lexer) backup() {
	l.pos--
}

func (l *Lexer) ignore() {
	l.start = l.pos
}

func (l *Lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.peek()) {
		l.next()
	}
	l.ignore()
}

func (l *Lexer) emit(t TokenType) {
	token := Token{Type: t, Value: l.input[l.start:l.pos]}
	enrichedToken := l.enrichToken(token)
	l.tokens <- enrichedToken
	l.start = l.pos
}

func (l *Lexer) Run() {
	for state := lexText; state != nil; {
		state = state(l)
	}
	close(l.tokens)
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		input:  input,
		tokens: make(chan Token),
	}
}

func (l *Lexer) GetTokens() <-chan Token {
	return l.tokens
}

func lexText(l *Lexer) stateFn {
	for {
		if l.peek() == '\n' {
			if l.pos > l.start {
				originalPos := l.pos
				l.emit(TokenText)
				l.start = originalPos
			}
			l.next()
			l.emit(TokenNewLine)
			l.start = l.pos
			continue
		}

		if l.pos == 0 || l.input[l.pos-1] == '\n' {
			if unicode.IsDigit(l.peek()) {
				if l.pos > l.start {
					l.emit(TokenText)
				}
				return lexPotentialOrderedList
			}
		}
		if l.pos < len(l.input)-2 && (strings.HasPrefix(l.input[l.pos:], "---") ||
			strings.HasPrefix(l.input[l.pos:], "***") ||
			strings.HasPrefix(l.input[l.pos:], "___")) {
			if l.pos > l.start {
				l.emit(TokenText)
			}
			return lexHorizontalRule
		}
		if strings.HasPrefix(l.input[l.pos:], "- ") ||
			strings.HasPrefix(l.input[l.pos:], "+ ") ||
			strings.HasPrefix(l.input[l.pos:], "* ") {
			if l.pos > l.start {
				l.emit(TokenText)
			}
			return lexUnorderedList
		}
		if strings.HasPrefix(l.input[l.pos:], "#") {
			if l.pos > l.start {
				l.emit(TokenText)
			}
			return lexHeading
		}
		if strings.HasPrefix(l.input[l.pos:], "**") || strings.HasPrefix(l.input[l.pos:], "__") {
			if l.pos > l.start {
				l.emit(TokenText)
			}
			return lexBold
		}
		if strings.HasPrefix(l.input[l.pos:], "*") || strings.HasPrefix(l.input[l.pos:], "_") {
			if l.pos > l.start {
				l.emit(TokenText)
			}
			return lexItalic
		}
		if strings.HasPrefix(l.input[l.pos:], "`") {
			if l.pos > l.start {
				l.emit(TokenText)
			}
			if strings.HasPrefix(l.input[l.pos:], "```") {
				return lexFencedCodeBlock
			}
			return lexInlineCode
		}
		if strings.HasPrefix(l.input[l.pos:], "[") {
			if l.pos > l.start {
				l.emit(TokenText)
			}
			return lexLink
		}
		if strings.HasPrefix(l.input[l.pos:], "![") {
			if l.pos > l.start {
				l.emit(TokenText)
			}
			return lexImage
		}
		if strings.HasPrefix(l.input[l.pos:], ">") {
			if l.pos == 0 || l.input[l.pos-1] == '\n' {
				if l.pos > l.start {
					l.emit(TokenText)
				}
				return lexBlockquote
			}
		}
		if l.next() == eof {
			break
		}
	}
	if l.pos > l.start {
		l.emit(TokenText)
	}
	l.emit(TokenEOF)
	return nil
}

func lexHeading(l *Lexer) stateFn {
	start := l.pos

	l.pos++
	for l.peek() == '#' {
		l.pos++
	}
	l.skipWhitespace()
	for l.peek() != '\n' && l.peek() != eof {
		l.pos++
	}

	l.start = start
	l.emit(TokenHeading)
	return lexText
}

func lexBold(l *Lexer) stateFn {
	marker := l.input[l.pos : l.pos+2]
	l.pos += 2
	for !strings.HasPrefix(l.input[l.pos:], marker) {
		if l.next() == eof {
			l.pos = l.start
			return lexText
		}
	}
	l.pos += 2
	l.emit(TokenBold)
	return lexText
}

func lexItalic(l *Lexer) stateFn {
	marker := l.input[l.pos : l.pos+1]
	start := l.pos
	l.pos++

	for !strings.HasPrefix(l.input[l.pos:], marker) {
		if l.next() == eof {
			l.pos = start
			l.next()
			l.consumeUntilEOF()
			l.emit(TokenText)
			return lexText
		}
	}
	l.pos++
	l.emit(TokenItalic)
	return lexText
}

func lexInlineCode(l *Lexer) stateFn {
	start := l.pos
	l.pos++

	for !strings.HasPrefix(l.input[l.pos:], "`") {
		if l.next() == eof {
			l.pos = start
			l.next()
			l.consumeUntilEOF()
			l.emit(TokenText)
			return lexText
		}
	}
	l.pos++
	l.emit(TokenInlineCode)
	return lexText
}

func lexFencedCodeBlock(l *Lexer) stateFn {
	start := l.pos
	l.pos += 3

	for !strings.HasPrefix(l.input[l.pos:], "```") {
		if l.next() == eof {
			l.pos = start
			l.next()
			l.consumeUntilEOF()
			l.emit(TokenText)
			return lexText
		}
	}
	l.pos += 3
	l.emit(TokenFencedCodeBlock)
	return lexText
}

func lexUnorderedList(l *Lexer) stateFn {
	l.pos += 2
	for l.peek() != '\n' && l.peek() != eof {
		l.pos++
	}
	l.emit(TokenUnorderedList)
	return lexText
}

func lexPotentialOrderedList(l *Lexer) stateFn {
	start := l.pos

	l.consumeWhile("0123456789")

	if l.accept(".") && l.accept(" ") {

		for l.peek() != '\n' && l.peek() != eof {
			l.pos++
		}
		l.start = start
		l.emit(TokenOrderedList)
	} else {

		l.pos = start
		l.next()
		for l.peek() != '\n' && l.peek() != eof {
			l.next()
		}
		l.emit(TokenText)
	}
	return lexText
}

func lexLink(l *Lexer) stateFn {
	start := l.pos
	l.pos++

	for l.peek() != ']' && l.peek() != eof {
		l.pos++
	}

	if l.peek() == ']' && strings.HasPrefix(l.input[l.pos+1:], "(") {
		l.pos++
		l.pos++
		for l.peek() != ')' && l.peek() != eof {
			l.pos++
		}
		if l.peek() == ')' {
			l.pos++
			l.start = start
			l.emit(TokenLink)
			return lexText
		}
	}

	l.pos = start
	l.next()
	l.consumeUntilEOF()
	l.emit(TokenText)
	return lexText
}

func lexImage(l *Lexer) stateFn {
	start := l.pos
	l.pos += 2

	for l.peek() != ']' && l.peek() != eof {
		l.pos++
	}

	if l.peek() == ']' && strings.HasPrefix(l.input[l.pos+1:], "(") {
		l.pos++
		l.pos++
		for l.peek() != ')' && l.peek() != eof {
			l.pos++
		}
		if l.peek() == ')' {
			l.pos++
			l.start = start
			l.emit(TokenImage)
			return lexText
		}
	}

	l.pos = start
	l.next()
	l.consumeUntilEOF()
	l.emit(TokenText)
	return lexText
}

func lexBlockquote(l *Lexer) stateFn {
	start := l.pos
	l.pos++
	l.skipWhitespace()
	for l.peek() != '\n' && l.peek() != eof {
		l.pos++
	}

	l.start = start
	l.emit(TokenBlockquote)
	return lexText
}

func lexHorizontalRule(l *Lexer) stateFn {
	start := l.pos
	marker := rune(l.input[l.pos])

	l.pos += 3

	for l.peek() == marker {
		l.pos++
	}

	l.start = start
	l.emit(TokenHorizontalRule)
	return lexText
}

func (l *Lexer) consumeWhile(valid string) {
	for {
		nextChar := l.next()
		if !strings.ContainsRune(valid, nextChar) {
			break
		}
		// Character was valid, continue to next
	}
	l.backup()
}

func (l *Lexer) consumeUntilEOF() {
	for {
		if l.next() == eof {
			break
		}
		// Continue consuming until EOF
	}
}
