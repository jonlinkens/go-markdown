package serializer

import (
	"encoding/json"
	"fmt"

	"github.com/jonlinkens/go-markdown/parser"
)

type NodeSerializer struct {
	Type     string           `json:"type"`
	Value    string           `json:"value,omitempty"`
	Meta     interface{}      `json:"meta,omitempty"`
	Children []NodeSerializer `json:"children,omitempty"`
}

func convertNode(node *parser.Node) NodeSerializer {
	serialized := NodeSerializer{
		Type:  node.Type.String(),
		Value: node.Value,
		Meta:  node.Meta,
	}

	if len(node.Children) > 0 {
		serialized.Children = make([]NodeSerializer, len(node.Children))
		for i, child := range node.Children {
			serialized.Children[i] = convertNode(child)
		}
	}

	return serialized
}

func ToJSON(doc *parser.Document) (string, error) {
	if doc == nil {
		return "", fmt.Errorf("document is nil")
	}

	serialized := convertNode(doc.Node)
	jsonData, err := json.MarshalIndent(serialized, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshaling JSON: %w", err)
	}

	return string(jsonData), nil
}

func FromJSON(jsonStr string) (*parser.Document, error) {
	var serialized NodeSerializer
	err := json.Unmarshal([]byte(jsonStr), &serialized)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	doc := parser.NewDocument()
	if serialized.Type != "Document" {
		node := convertSerializedNode(serialized)
		doc.AddChild(node)
	} else {
		for _, child := range serialized.Children {
			node := convertSerializedNode(child)
			doc.AddChild(node)
		}
	}

	return doc, nil
}

func convertSerializedNode(serialized NodeSerializer) *parser.Node {
	var nodeType parser.NodeType
	if !nodeType.FromString(serialized.Type) {
		nodeType = parser.NodeText
	}

	node := parser.NewNode(nodeType, serialized.Value)
	node.Meta = serialized.Meta

	if len(serialized.Children) > 0 {
		node.Children = make([]*parser.Node, len(serialized.Children))
		for i, child := range serialized.Children {
			node.Children[i] = convertSerializedNode(child)
		}
	}

	return node
}
