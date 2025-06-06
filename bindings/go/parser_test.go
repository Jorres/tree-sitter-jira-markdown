package tree_sitter_markdown_test

import (
	"testing"

	tree_sitter_markdown "github.com/tree-sitter-grammars/tree-sitter-markdown/bindings/go"
)

func TestParseSimpleHeader(t *testing.T) {
	parser := tree_sitter_markdown.NewAdfMarkdownParser()
	content := []byte("# header\n")

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
						Kind: "atx_heading",
						Children: []*TreeNode{
							{Kind: "atx_h1_marker", Text: "#"},
							{Kind: "inline", Text: "header"},
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

func TestParseSimplePeopleMention(t *testing.T) {
	parser := tree_sitter_markdown.NewAdfMarkdownParser()
	content := []byte("@jorres@nebius.com")

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
								Text: "@jorres@nebius.com",
								Children: []*TreeNode{
									{
										Kind: "people_mention",
										Text: "@jorres@nebius.com",
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

func TestParseSimpleAttachment(t *testing.T) {
	parser := tree_sitter_markdown.NewAdfMarkdownParser()
	content := []byte("{attachment:file.txt}")

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
						Kind: "attachment",
						Children: []*TreeNode{
							{
								Kind: "attachment_start_mark",
								Text: "{attachment:",
							},
							{
								Kind: "attachment_path",
								Text: "file.txt",
							},
							{
								Kind: "attachment_end_mark",
								Text: "}",
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

func TestParsePanelWithType(t *testing.T) {
    parser := tree_sitter_markdown.NewAdfMarkdownParser()
    content := []byte(`{panel:type=info}
# Header

possible **text-inside** with ` + "`" + `inline stuff` + "`" + `

` + "```" + `
and code blocks
` + "```" + `
{/panel}`)

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
                        Kind: "panel",
                        Children: []*TreeNode{
                            {
                                Kind: "panel_start",
                                Children: []*TreeNode{
                                    {Kind: "panel_start_mark", Text: "{panel:"},
                                    {
                                        Kind: "panel_type",
                                        Children: []*TreeNode{
                                            {Kind: "=", Text: "="},
                                            {Kind: "type", Text: "info"},
                                        },
                                    },
                                    {Kind: "panel_end_bracket", Text: "}"},
                                },
                            },
                            {
                                Kind: "section",
                                Children: []*TreeNode{
                                    {
                                        Kind: "atx_heading",
                                        Children: []*TreeNode{
                                            {Kind: "atx_h1_marker", Text: "#"},
                                            {Kind: "inline", Text: "Header"},
                                        },
                                    },
                                    {
                                        Kind: "paragraph",
                                        Children: []*TreeNode{
                                            {
                                                Kind: "inline",
                                                Text: "possible **text-inside** with `inline stuff`",
                                                Children: []*TreeNode{
                                                    {
                                                        Kind: "strong_emphasis",
                                                        Children: []*TreeNode{
                                                            {Kind: "emphasis_delimiter", Text: "*"},
                                                            {Kind: "emphasis_delimiter", Text: "*"},
                                                            {Kind: "-", Text: "-"},
                                                            {Kind: "emphasis_delimiter", Text: "*"},
                                                            {Kind: "emphasis_delimiter", Text: "*"},
                                                        },
                                                    },
                                                    {
                                                        Kind: "code_span",
                                                        Children: []*TreeNode{
                                                            {Kind: "code_span_delimiter", Text: "`"},
                                                            {Kind: "code_span_delimiter", Text: "`"},
                                                        },
                                                    },
                                                },
                                            },
                                        },
                                    },
                                    {
                                        Kind: "fenced_code_block",
                                        Children: []*TreeNode{
                                            {Kind: "fenced_code_block_delimiter", Text: "```"},
                                            {Kind: "block_continuation"},
                                            {
                                                Kind: "code_fence_content",
                                                Children: []*TreeNode{
                                                    {Kind: "block_continuation"},
                                                },
                                            },
                                            {Kind: "fenced_code_block_delimiter", Text: "```"},
                                        },
                                    },
                                },
                            },
                            {Kind: "panel_end_mark", Text: "{/panel}"},
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

func TestParseComplexDocument(t *testing.T) {
	parser := tree_sitter_markdown.NewAdfMarkdownParser()
	content := []byte(`# Main Header

1. Item with ` + "`" + `code` + "`" + `span
2. Item with @user@example.com mention

{attachment:file.txt}

` + "```" + `
code block content
` + "```\n")

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
					// Header
					{
						Kind: "atx_heading",
						Children: []*TreeNode{
							{Kind: "atx_h1_marker", Text: "#"},
							{Kind: "inline", Text: "Main Header", Children: []*TreeNode{}},
						},
					},
					// Numbered list
					{
						Kind: "list",
						Children: []*TreeNode{
							// Item 1 with code span
							{
								Kind: "list_item",
								Children: []*TreeNode{
									{Kind: "list_marker_dot", Text: "1."},
									{
										Kind: "paragraph",
										Children: []*TreeNode{
											{
												Kind: "inline",
												Text: "Item with `code`span",
												Children: []*TreeNode{
													{
														Kind: "code_span",
														Children: []*TreeNode{
															{Kind: "code_span_delimiter", Text: "`"},
															{Kind: "code_span_delimiter", Text: "`"},
														},
													},
												},
											},
										},
									},
								},
							},
							// Item 2 with people mention
							{
								Kind: "list_item",
								Children: []*TreeNode{
									{Kind: "list_marker_dot", Text: "2."},
									{
										Kind: "paragraph",
										Children: []*TreeNode{
											{
												Kind: "inline",
												Text: "Item with @user@example.com mention",
												Children: []*TreeNode{
													{Kind: "people_mention", Text: "@user@example.com"},
												},
											},
											{Kind: "block_continuation"},
										},
									},
								},
							},
						},
					},

					{
						Kind: "attachment",
						Children: []*TreeNode{
							{
								Kind: "attachment_start_mark",
								Text: "{attachment:",
							},
							{
								Kind: "attachment_path",
								Text: "file.txt",
							},
							{
								Kind: "attachment_end_mark",
								Text: "}",
							},
						},
					},
					// Code block
					{
						Kind: "fenced_code_block",
						Children: []*TreeNode{
							{Kind: "fenced_code_block_delimiter", Text: "```"},
							{Kind: "block_continuation"},
							{
								Kind: "code_fence_content",
								Children: []*TreeNode{
									{Kind: "block_continuation"},
								},
							},
							{Kind: "fenced_code_block_delimiter", Text: "```"},
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
