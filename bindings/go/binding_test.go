package tree_sitter_markdown_test

import (
	"testing"

	tree_sitter_markdown "github.com/jorres/tree-sitter-jira-markdown/bindings/go"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

func TestCanLoadBlockGrammar(t *testing.T) {
	language := tree_sitter.NewLanguage(tree_sitter_markdown.Language())
	if language == nil {
		t.Errorf("Error loading Markdown block grammar")
	}
}

func TestCanLoadInlineGrammar(t *testing.T) {
	language := tree_sitter.NewLanguage(tree_sitter_markdown.InlineLanguage())
	if language == nil {
		t.Errorf("Error loading Markdown inline grammar")
	}
}

func TestParseUnderlineInBold(t *testing.T) {
	parser := tree_sitter_markdown.NewAdfMarkdownParser()
	content := []byte("**<u>text</u>**\n")

	tree, err := parser.Parse(content)
	if err != nil {
		t.Fatalf("Error parsing markdown: %v", err)
	}

	// Build expected tree structure
	expected := &TreeNode{
		Kind: "document",
		Children: []*TreeNode{
			{
				Kind: "section",
				Children: []*TreeNode{
					{
						Kind: "paragraph",
						Children: []*TreeNode{
							{
								Kind: "inline",
								Text: "**<u>text</u>**",
								Children: []*TreeNode{
									{
										Kind: "strong_emphasis",
										Children: []*TreeNode{
											{Kind: "emphasis_delimiter", Text: "*"},
											{Kind: "emphasis_delimiter", Text: "*"},
											{
												Kind: "underline",
												Children: []*TreeNode{
													{Kind: "underline_open", Text: "<u>"},
													{Kind: "underline_content", Text: "text"},
													{Kind: "underline_close", Text: "</u>"},
												},
											},
											{Kind: "emphasis_delimiter", Text: "*"},
											{Kind: "emphasis_delimiter", Text: "*"},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	// Convert parsed tree to our comparison format
	actual := convertToTreeNode(tree.RootNode(), content, parser)

	// Compare trees
	if !compareTreeNodes(expected, actual) {
		t.Errorf("Tree structure doesn't match.\nExpected:\n%s\nActual:\n%s",
			printTree(expected, 0), printTree(actual, 0))
	}
}
