# go-markdown

> [!NOTE]
> This is **not standard-compliant** according to any specification such as [CommonMark](https://commonmark.org/) or [GitHub Markdown](https://docs.github.com/en/get-started/writing-on-github/getting-started-with-writing-and-formatting-on-github/basic-writing-and-formatting-syntax). If you're looking for something like that, please use [goldmark](https://github.com/yuin/goldmark)

A simple, zero-dependency library that tokenizes raw markdown and offers serialization to JSON. Its approach is very much inspired by Rob Pike's talk on Lexical Scanning in Go found [here](https://www.youtube.com/watch?v=HxaD_trXwRE). Expect bugs here and there, as well as some very messy code - this is a quick project I did over a few weekends!

https://github.com/user-attachments/assets/5a197516-1804-40f3-9fd0-87df05495162

> The example above involves naive file watching, and so is much slower than in practice.

`go-markdown` can process markdown (2000+ lines) files in <3.5ms:

<p align="center"><img width="100%" alt="Screenshot 2025-02-04 at 19 29 58" src="https://github.com/user-attachments/assets/7e80af67-8f6d-4b4f-81d5-c7c3e7f27cf2" /></p>

---

The library is split into several packages:

- `lexer`: Handles tokenization of raw markdown text
- `parser`: Converts flat token stream into a hierarchical AST
- `serializer`: Provides JSON serialization for both tokens and AST nodes
- `common`: Shared functionality like type mapping between packages

Each package is designed to be used independently, so you can use just the lexer if you only need tokens, or the full pipeline for AST generation.

## Example usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/jonlinkens/go-markdown/lexer"
	"github.com/jonlinkens/go-markdown/parser"
	"github.com/jonlinkens/go-markdown/serializer"
)

func main() {
	input := []byte("# Hello\nThis is **bold** text")

	l := lexer.NewLexer(string(input))
	go l.Run()

	var tokens []lexer.Token
	for token := range l.GetTokens() {
		tokens = append(tokens, token)
	}

	p := parser.NewParser(tokens)
	doc := p.Parse()

	json, err := serializer.ToJSON(doc)
	if err != nil {
		log.Fatalf("Error serializing AST: %v", err)
	}

	fmt.Println(json)
}
```

## Supported Elements

All markdown elements are represented as both tokens and AST nodes. Supported elements can be found in:

- Tokens: [`lexer/token.go`](https://github.com/jonlinkens/go-markdown/blob/main/lexer/token.go)
- AST Nodes: [`parser/node.go`](https://github.com/jonlinkens/go-markdown/blob/main/parser/node.go)
