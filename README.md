# go-markdown

> [!NOTE]
> This is **not standard-compliant** according to any specification such as [CommonMark](https://commonmark.org/) or [GitHub Markdown](https://docs.github.com/en/get-started/writing-on-github/getting-started-with-writing-and-formatting-on-github/basic-writing-and-formatting-syntax). If you're looking for something like that, please use [goldmark](https://github.com/yuin/goldmark)

A simple, zero-dependency library that tokenizes raw markdown and offers serialization to JSON. Its approach is very much inspired by Rob Pike's talk on Lexical Scanning in Go found [here](https://www.youtube.com/watch?v=HxaD_trXwRE). Expect bugs here and there, as well as some very messy code - this is a quick project I did over a few weekends!

https://github.com/user-attachments/assets/5a197516-1804-40f3-9fd0-87df05495162

> The example above involves naive file watching, and so is much slower than in practice. 


`go-markdown` can process markdown (2000+ lines) files in >3.5ms:

<p align="center"><img width="100%" alt="Screenshot 2025-02-04 at 19 29 58" src="https://github.com/user-attachments/assets/7e80af67-8f6d-4b4f-81d5-c7c3e7f27cf2" /></p>

## Example usage

```go
func main() {
    input, err := os.ReadFile("myFile.md")
    if err != nil {
        log.Fatalf("Error reading file: %v", err)
    }

    l := lexer.NewLexer(string(input))

    go l.Run()

    var tokens []lexer.Token

    for token := range l.GetTokens() {
        tokens = append(tokens, token)
    }

    tokenSlice := serializer.TokenSlice(tokens)
    tokenJson, err := tokenSlice.ToJson()

    fmt.Println(tokenJson)
}
```

Supported elements can be found in `lexer/token.go`:
https://github.com/jonlinkens/go-markdown/blob/b7eb4f01d5cdf83cd7e1300801894c2d27564c29/lexer/token.go#L10-L25
