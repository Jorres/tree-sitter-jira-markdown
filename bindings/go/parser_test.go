package tree_sitter_markdown_test

import (
	"fmt"
	"testing"

	tree_sitter_markdown "github.com/tree-sitter-grammars/tree-sitter-markdown/bindings/go"
	sitter "github.com/tree-sitter/go-tree-sitter"
)

func TestParseSimpleHeader(t *testing.T) {
	parser := tree_sitter_markdown.NewAdfMarkdownParser()
	content := []byte("# header\n")

	tree, err := parser.Parse(content)
	if err != nil {
		t.Fatalf("Error parsing markdown: %v", err)
	}

	if tree == nil {
		t.Fatal("Parse returned nil tree")
	}

	root := tree.RootNode()
	if root == nil {
		t.Fatal("Root node is nil")
	}

	// Verify we have a document root
	if root.Kind() != "document" {
		t.Errorf("Expected root node to be 'document', got '%s'", root.Kind())
	}

	// Check that we have children
	if root.ChildCount() == 0 {
		t.Fatal("Document should have children")
	}

	// Find the heading node - may be nested in sections
	headingNode := findNodeByType(root, "atx_heading")
	if headingNode == nil {
		t.Fatal("Could not find atx_heading node")
	}

	// Verify heading content
	headingText := string(content[headingNode.StartByte():headingNode.EndByte()])
	if headingText != "# header\n" {
		t.Errorf("Expected heading text '# header', got '%s'", headingText)
	}

	// Check for inline content within the heading
	inlineNode := findNodeByType(headingNode, "inline")
	if inlineNode != nil {
		// Get the inline tree and check its content
		inlineTree := parser.GetInlineTree(inlineNode, content)
		if inlineTree != nil {
			inlineText := string(content[inlineNode.StartByte():inlineNode.EndByte()])
			expectedInlineText := "header"
			if inlineText != expectedInlineText {
				t.Errorf("Expected inline text '%s', got '%s'", expectedInlineText, inlineText)
			}
		}
	}
}

// Helper function to recursively find a node by type
func findNodeByType(node *sitter.Node, nodeType string) *sitter.Node {
	if node.Kind() == nodeType {
		return node
	}

	childCount := node.ChildCount()
	for i := uint(0); i < childCount; i++ {
		child := node.Child(i)
		if child != nil {
			if found := findNodeByType(child, nodeType); found != nil {
				return found
			}
		}
	}
	return nil
}

