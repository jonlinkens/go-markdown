package parser

import "github.com/jonlinkens/go-markdown/lexer"

type NodeType int

const (
	NodeDocument NodeType = iota
	NodeParagraph
	NodeHeading
	NodeText
	NodeBold
	NodeItalic
	NodeInlineCode
	NodeCodeBlock
	NodeUnorderedList
	NodeOrderedList
	NodeListItem
	NodeLink
	NodeImage
	NodeBlockquote
	NodeHorizontalRule
)

var NodeTypeStrings = map[NodeType]string{
	NodeDocument:       "Document",
	NodeParagraph:      "Paragraph",
	NodeHeading:        "Heading",
	NodeText:           "Text",
	NodeBold:           "Bold",
	NodeItalic:         "Italic",
	NodeInlineCode:     "InlineCode",
	NodeCodeBlock:      "CodeBlock",
	NodeUnorderedList:  "UnorderedList",
	NodeOrderedList:    "OrderedList",
	NodeListItem:       "ListItem",
	NodeLink:           "Link",
	NodeImage:          "Image",
	NodeBlockquote:     "Blockquote",
	NodeHorizontalRule: "HorizontalRule",
}

var StringToNodeType = func() map[string]NodeType {
	m := make(map[string]NodeType, len(NodeTypeStrings))
	for nodeType, str := range NodeTypeStrings {
		m[str] = nodeType
	}
	return m
}()

func (n *NodeType) FromString(s string) bool {
	if nodeType, ok := StringToNodeType[s]; ok {
		*n = nodeType
		return true
	}
	return false
}

type Node struct {
	Type     NodeType     `json:"type"`
	Value    string       `json:"value"`
	Meta     any          `json:"meta,omitempty"`
	Children []*Node      `json:"children,omitempty"`
	Token    *lexer.Token `json:"-"`
}

type Document struct {
	*Node
}

func NewDocument() *Document {
	return &Document{
		Node: &Node{
			Type:     NodeDocument,
			Children: make([]*Node, 0),
		},
	}
}

func NewNode(nodeType NodeType, value string) *Node {
	return &Node{
		Type:     nodeType,
		Value:    value,
		Children: make([]*Node, 0),
	}
}

func (n *Node) AddChild(child *Node) {
	n.Children = append(n.Children, child)
}

func (n NodeType) String() string {
	if str, ok := NodeTypeStrings[n]; ok {
		return str
	}
	return "Unknown"
}
