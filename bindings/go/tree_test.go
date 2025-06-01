package tree_sitter_markdown_test

import (
	"strings"

	tree_sitter_markdown "github.com/tree-sitter-grammars/tree-sitter-markdown/bindings/go"
	sitter "github.com/tree-sitter/go-tree-sitter"
)

// TreeNode represents a simplified tree structure for comparison
type TreeNode struct {
	Kind     string
	Text     string
	Children []*TreeNode
}

// convertToTreeNode converts a tree-sitter node to our TreeNode format
func convertToTreeNode(node *sitter.Node, content []byte, parser *tree_sitter_markdown.AdfMarkdownParser) *TreeNode {
	result := &TreeNode{
		Kind: node.Kind(),
	}

	// Get text for leaf nodes or nodes with specific text content
	if node.ChildCount() == 0 || node.Kind() == "atx_h1_marker" || node.Kind() == "inline" {
		result.Text = strings.TrimSpace(string(content[node.StartByte():node.EndByte()]))
	}

	if node.Kind() == "inline" {
		inlineTree := parser.InlineToTree[node.Id()]
		if inlineTree != nil {
			return convertToTreeNode(
				inlineTree.RootNode(),
				content[node.StartByte():node.EndByte()],
				parser,
			)
		}
	}

	// Convert children
	childCount := node.ChildCount()
	for i := uint(0); i < childCount; i++ {
		child := node.Child(i)
		if child != nil {
			result.Children = append(result.Children, convertToTreeNode(child, content, parser))
		}
	}

	return result
}

// compareTreeNodes recursively compares two TreeNode structures
func compareTreeNodes(expected, actual *TreeNode) bool {
	if expected.Kind != actual.Kind {
		return false
	}

	if expected.Text != "" && expected.Text != actual.Text {
		return false
	}

	if len(expected.Children) != len(actual.Children) {
		return false
	}

	for i, expectedChild := range expected.Children {
		if !compareTreeNodes(expectedChild, actual.Children[i]) {
			return false
		}
	}

	return true
}

// printTree creates a readable string representation of the tree
func printTree(node *TreeNode, depth int) string {
	indent := strings.Repeat("  ", depth)
	result := indent + node.Kind
	if node.Text != "" {
		result += ": \"" + node.Text + "\""
	}
	result += "\n"

	for _, child := range node.Children {
		result += printTree(child, depth+1)
	}

	return result
}
