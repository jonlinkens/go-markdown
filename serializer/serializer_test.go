package serializer

import (
	"testing"

	"github.com/jonlinkens/go-markdown/parser"
)

func TestSerializerBasic(t *testing.T) {
	doc := parser.NewDocument()
	para := parser.NewNode(parser.NodeParagraph, "")
	text := parser.NewNode(parser.NodeText, "Hello ")
	bold := parser.NewNode(parser.NodeBold, "world")
	text2 := parser.NewNode(parser.NodeText, "!")

	para.AddChild(text)
	para.AddChild(bold)
	para.AddChild(text2)
	doc.AddChild(para)

	json, err := ToJSON(doc)
	if err != nil {
		t.Errorf("ToJSON failed: %v", err)
	}

	newDoc, err := FromJSON(json)
	if err != nil {
		t.Errorf("FromJSON failed: %v", err)
	}

	if len(newDoc.Children) != 1 {
		t.Errorf("Expected 1 child in document, got %d", len(newDoc.Children))
	}

	para = newDoc.Children[0]
	if para.Type != parser.NodeParagraph {
		t.Errorf("Expected paragraph node, got %s", para.Type)
	}

	if len(para.Children) != 3 {
		t.Errorf("Expected 3 children in paragraph, got %d", len(para.Children))
	}

	if para.Children[0].Type != parser.NodeText || para.Children[0].Value != "Hello " {
		t.Errorf("Expected text node with 'Hello ', got %s with '%s'", para.Children[0].Type, para.Children[0].Value)
	}

	if para.Children[1].Type != parser.NodeBold || para.Children[1].Value != "world" {
		t.Errorf("Expected bold node with 'world', got %s with '%s'", para.Children[1].Type, para.Children[1].Value)
	}

	if para.Children[2].Type != parser.NodeText || para.Children[2].Value != "!" {
		t.Errorf("Expected text node with '!', got %s with '%s'", para.Children[2].Type, para.Children[2].Value)
	}
}

func TestSerializerNilDocument(t *testing.T) {
	_, err := ToJSON(nil)
	if err == nil {
		t.Error("Expected error when serializing nil document")
	}
}

func TestSerializerInvalidJSON(t *testing.T) {
	_, err := FromJSON("invalid json")
	if err == nil {
		t.Error("Expected error when deserializing invalid JSON")
	}
}

func TestSerializerUnknownNodeType(t *testing.T) {
	json := `{"type": "UnknownType", "value": "test"}`
	doc, err := FromJSON(json)
	if err != nil {
		t.Errorf("FromJSON failed: %v", err)
	}

	if len(doc.Children) != 1 {
		t.Errorf("Expected document to have 1 child, got %d", len(doc.Children))
		return
	}

	node := doc.Children[0]
	if node.Type != parser.NodeText {
		t.Errorf("Expected unknown type to default to text node, got %s", node.Type)
	}

	if node.Value != "test" {
		t.Errorf("Expected value 'test', got '%s'", node.Value)
	}
}
