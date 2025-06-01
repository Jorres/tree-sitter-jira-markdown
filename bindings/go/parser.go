package tree_sitter_markdown

import (
	"fmt"

	sitter "github.com/tree-sitter/go-tree-sitter"
)

// AdfMarkdownParser encapsulates the dual-grammar complexity of tree-sitter-markdown
type AdfMarkdownParser struct {
	blockParser    *sitter.Parser
	inlineParser   *sitter.Parser
	blockLanguage  *sitter.Language
	inlineLanguage *sitter.Language

	InlineToTree map[uintptr]*sitter.Tree
}

// NewAdfMarkdownParser creates a new parser that handles both block and inline grammars
func NewAdfMarkdownParser() *AdfMarkdownParser {
	blockLanguage := sitter.NewLanguage(Language())
	inlineLanguage := sitter.NewLanguage(InlineLanguage())

	blockParser := sitter.NewParser()
	err := blockParser.SetLanguage(blockLanguage)
	if err != nil {
		panic(err)
	}

	inlineParser := sitter.NewParser()
	err = inlineParser.SetLanguage(inlineLanguage)
	if err != nil {
		panic(err)
	}

	return &AdfMarkdownParser{
		blockParser:    blockParser,
		inlineParser:   inlineParser,
		blockLanguage:  blockLanguage,
		inlineLanguage: inlineLanguage,

		InlineToTree: make(map[uintptr]*sitter.Tree),
	}
}

// Parse parses markdown content and returns a unified tree with inline content processed
// This hides the complexity of dual-grammar parsing from the application
func (p *AdfMarkdownParser) Parse(content []byte) (*sitter.Tree, error) {
	// Parse the document structure with block grammar
	blockTree := p.blockParser.Parse(content, nil)
	if blockTree == nil {
		return nil, fmt.Errorf("failed to parse with block grammar")
	}

	// Process inline content and embed it into the block tree
	p.processInlineContent(blockTree.RootNode(), content)

	return blockTree, nil
}

// processInlineContent finds inline nodes and processes them with the inline grammar
func (p *AdfMarkdownParser) processInlineContent(node *sitter.Node, content []byte) {
	if node.Kind() == "inline" {
		// Parse this inline content with the inline grammar
		inlineContent := content[node.StartByte():node.EndByte()]
		inlineTree := p.inlineParser.Parse(inlineContent, nil)

		if inlineTree != nil {
			p.attachInlineTree(node, inlineTree, content)
		}
	}

	// Recursively process children
	childCount := int(node.ChildCount())
	for i := range childCount {
		child := node.Child(uint(i))
		if child != nil {
			p.processInlineContent(child, content)
		}
	}
}

// attachInlineTree attaches an inline tree to a block node
// This is where we could implement a more sophisticated mapping strategy
func (p *AdfMarkdownParser) attachInlineTree(blockNode *sitter.Node, inlineTree *sitter.Tree, content []byte) {
	p.InlineToTree[blockNode.Id()] = inlineTree
}

// GetInlineTree retrieves the inline tree for a given block inline node
// This provides access to the processed inline content
func (p *AdfMarkdownParser) GetInlineTree(blockInlineNode *sitter.Node, content []byte) *sitter.Tree {
	if blockInlineNode.Kind() != "inline" {
		return nil
	}

	// Extract and parse the inline content
	inlineContent := content[blockInlineNode.StartByte():blockInlineNode.EndByte()]
	return p.inlineParser.Parse(inlineContent, nil)
}
