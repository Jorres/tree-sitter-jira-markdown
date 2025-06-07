package main

import (
	"fmt"
	"os"
	"strings"

	tree_sitter_markdown "github.com/jorres/tree-sitter-jira-markdown/bindings/go"
	sitter "github.com/tree-sitter/go-tree-sitter"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <markdown_file>\n", os.Args[0])
		os.Exit(1)
	}

	filename := os.Args[1]
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file %q: %v\n", filename, err)
		os.Exit(1)
	}

	parser := tree_sitter_markdown.NewAdfMarkdownParser()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating parser: %v\n", err)
		os.Exit(1)
	}

	tree, err := parser.Parse([]byte(content))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing: %v\n", err)
		os.Exit(1)
	}

	dumpTree(parser, tree.RootNode(), []byte(content), 0)
}

// dumpTree dumps a tree-sitter tree in the standard format
func dumpTree(parser *tree_sitter_markdown.AdfMarkdownParser, node *sitter.Node, content []byte, depth int) {
	indent := strings.Repeat("  ", depth)

	// Print the node with position info
	nodeText := string(content[node.StartByte():node.EndByte()])

	if node.Kind() == "inline" {
		if inlineTree, exists := parser.InlineToTree[node.Id()]; exists {
			dumpTree(parser, inlineTree.RootNode(), []byte(nodeText), depth+1)
			return
		}
	}

	// Escape special characters for display
	displayText := strings.ReplaceAll(nodeText, "\n", "\\n")
	displayText = strings.ReplaceAll(displayText, "\t", "\\t")
	if len(displayText) > 50 {
		displayText = displayText[:47] + "..."
	}

	fmt.Printf("%s(%s [%d, %d] - [%d, %d]",
		indent,
		node.Kind(),
		node.StartPosition().Row, node.StartPosition().Column,
		node.EndPosition().Row, node.EndPosition().Column)

	if strings.TrimSpace(nodeText) != "" && node.ChildCount() == 0 {
		fmt.Printf(" \"%s\"", displayText)
	}
	fmt.Printf(")\n")

	// Recursively print children
	childCount := node.ChildCount()
	for i := uint(0); i < childCount; i++ {
		child := node.Child(i)
		if child != nil {
			dumpTree(parser, child, content, depth+1)
		}
	}
}
