package parser

import (
	"strings"

	"github.com/jonlinkens/go-markdown/lexer"
)

type Parser struct {
	tokens []lexer.Token
	pos    int
}

func NewParser(tokens []lexer.Token) *Parser {
	return &Parser{
		tokens: tokens,
		pos:    0,
	}
}

func (p *Parser) Parse() *Document {
	doc := NewDocument()
	currentParagraph := NewNode(NodeParagraph, "")
	var currentListItem *Node
	var listItemContent []string
	var hasInlineContent bool

	for !p.isAtEnd() {
		token := p.current()

		switch token.Type {
		case lexer.TokenEOF:
			if len(currentParagraph.Children) > 0 {
				doc.AddChild(currentParagraph)
			}
			return doc

		case lexer.TokenBreak:
			if len(currentParagraph.Children) > 0 {
				doc.AddChild(currentParagraph)
				currentParagraph = NewNode(NodeParagraph, "")
			}
			currentListItem = nil
			listItemContent = nil
			hasInlineContent = false

		case lexer.TokenHeading:
			if len(currentParagraph.Children) > 0 {
				doc.AddChild(currentParagraph)
				currentParagraph = NewNode(NodeParagraph, "")
			}
			node := NewNode(NodeHeading, token.CleanValue)
			node.Meta = token.Meta
			doc.AddChild(node)
			currentListItem = nil
			listItemContent = nil
			hasInlineContent = false

		case lexer.TokenBold, lexer.TokenItalic, lexer.TokenInlineCode, lexer.TokenLink:
			node := NewNode(nodeTypeFromToken(token.Type), token.CleanValue)
			if token.Meta != nil {
				node.Meta = token.Meta
			}
			if currentListItem != nil {
				if !hasInlineContent {
					currentListItem.Children = make([]*Node, 0)
					if len(listItemContent) > 0 {
						textNode := NewNode(NodeText, listItemContent[0])
						currentListItem.AddChild(textNode)
					}
					hasInlineContent = true
				}
				currentListItem.AddChild(node)
				listItemContent = append(listItemContent, token.CleanValue)
			} else {
				currentParagraph.AddChild(node)
			}

		case lexer.TokenText:
			if currentListItem != nil {
				if hasInlineContent {
					node := NewNode(NodeText, token.CleanValue)
					currentListItem.AddChild(node)
				}
				listItemContent = append(listItemContent, token.CleanValue)
			} else {
				node := NewNode(NodeText, token.CleanValue)
				currentParagraph.AddChild(node)
			}

		case lexer.TokenFencedCodeBlock:
			if len(currentParagraph.Children) > 0 {
				doc.AddChild(currentParagraph)
				currentParagraph = NewNode(NodeParagraph, "")
			}
			node := NewNode(NodeCodeBlock, token.CleanValue)
			node.Meta = token.Meta
			doc.AddChild(node)
			currentListItem = nil
			listItemContent = nil
			hasInlineContent = false

		case lexer.TokenUnorderedList:
			if len(currentParagraph.Children) > 0 {
				doc.AddChild(currentParagraph)
				currentParagraph = NewNode(NodeParagraph, "")
			}

			var listNode *Node
			if len(doc.Children) > 0 && doc.Children[len(doc.Children)-1].Type == NodeUnorderedList {
				listNode = doc.Children[len(doc.Children)-1]
			} else {
				listNode = NewNode(NodeUnorderedList, "")
				doc.AddChild(listNode)
			}

			currentListItem = NewNode(NodeListItem, "")
			listNode.AddChild(currentListItem)
			listItemContent = nil
			hasInlineContent = false

			if token.CleanValue != "" {
				listItemContent = []string{token.CleanValue}
			}

		case lexer.TokenOrderedList:
			if len(currentParagraph.Children) > 0 {
				doc.AddChild(currentParagraph)
				currentParagraph = NewNode(NodeParagraph, "")
			}

			var listNode *Node
			if len(doc.Children) > 0 && doc.Children[len(doc.Children)-1].Type == NodeOrderedList {
				listNode = doc.Children[len(doc.Children)-1]
			} else {
				listNode = NewNode(NodeOrderedList, "")
				doc.AddChild(listNode)
			}

			currentListItem = NewNode(NodeListItem, "")
			currentListItem.Meta = token.Meta
			listNode.AddChild(currentListItem)
			listItemContent = nil
			hasInlineContent = false

			if token.CleanValue != "" {
				listItemContent = []string{token.CleanValue}
			}

		case lexer.TokenImage:
			node := NewNode(NodeImage, token.CleanValue)
			node.Meta = token.Meta
			if currentListItem != nil {
				if !hasInlineContent {
					currentListItem.Children = make([]*Node, 0)
					if len(listItemContent) > 0 {
						textNode := NewNode(NodeText, listItemContent[0])
						currentListItem.AddChild(textNode)
					}
					hasInlineContent = true
				}
				currentListItem.AddChild(node)
				listItemContent = append(listItemContent, token.CleanValue)
			} else {
				currentParagraph.AddChild(node)
			}

		case lexer.TokenBlockquote:
			if len(currentParagraph.Children) > 0 {
				doc.AddChild(currentParagraph)
				currentParagraph = NewNode(NodeParagraph, "")
			}
			node := NewNode(NodeBlockquote, token.CleanValue)
			node.Meta = token.Meta
			doc.AddChild(node)
			currentListItem = nil
			listItemContent = nil
			hasInlineContent = false

		case lexer.TokenHorizontalRule:
			if len(currentParagraph.Children) > 0 {
				doc.AddChild(currentParagraph)
				currentParagraph = NewNode(NodeParagraph, "")
			}
			node := NewNode(NodeHorizontalRule, token.CleanValue)
			doc.AddChild(node)
			currentListItem = nil
			listItemContent = nil
			hasInlineContent = false
		}

		if currentListItem != nil {
			currentListItem.Value = strings.Join(listItemContent, "")
		}

		p.advance()
	}

	if len(currentParagraph.Children) > 0 {
		doc.AddChild(currentParagraph)
	}

	return doc
}

func (p *Parser) current() lexer.Token {
	return p.tokens[p.pos]
}

func (p *Parser) advance() {
	p.pos++
}

func (p *Parser) isAtEnd() bool {
	return p.pos >= len(p.tokens)
}

func nodeTypeFromToken(tokenType lexer.TokenType) NodeType {
	switch tokenType {
	case lexer.TokenText:
		return NodeText
	case lexer.TokenBold:
		return NodeBold
	case lexer.TokenItalic:
		return NodeItalic
	case lexer.TokenInlineCode:
		return NodeInlineCode
	case lexer.TokenLink:
		return NodeLink
	default:
		return NodeText
	}
}
