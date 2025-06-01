package tree_sitter_markdown

import (
	"fmt"

	sitter "github.com/tree-sitter/go-tree-sitter"
)

// ExpMarkdownParser encapsulates the dual-grammar complexity of tree-sitter-markdown
type ExpMarkdownParser struct {
	parser         *sitter.Parser
	blockLanguage  *sitter.Language
	inlineLanguage *sitter.Language
}

// MarkdownTree holds a combined markdown tree with both block and inline content
type MarkdownTree struct {
	blockTree     *sitter.Tree
	inlineTrees   []*sitter.Tree
	inlineIndices map[uintptr]*sitter.Tree
}

// MarkdownCursor provides a unified interface for walking the combined tree
type MarkdownCursor struct {
	markdownTree *MarkdownTree
	blockCursor  *sitter.TreeCursor
	inlineCursor *sitter.TreeCursor
}

// NewMarkdownParser creates a new parser that handles both block and inline grammars
func NewMarkdownParser() (*ExpMarkdownParser, error) {
	blockLanguage := sitter.NewLanguage(Language())
	inlineLanguage := sitter.NewLanguage(InlineLanguage())

	parser := sitter.NewParser()
	// err := parser.SetLanguage(blockLanguage)
	// if err != nil {
	// 	return nil, err
	// }

	return &ExpMarkdownParser{
		parser:         parser,
		blockLanguage:  blockLanguage,
		inlineLanguage: inlineLanguage,
	}, nil
}

// Parse parses markdown content and returns a MarkdownTree with inline content processed
func (p *ExpMarkdownParser) Parse(content []byte) (*MarkdownTree, error) {
	// Parse the document structure with block grammar
	p.parser.Reset()
	err := p.parser.SetLanguage(p.blockLanguage)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(content))
	blockTree := p.parser.Parse(content, nil)
	if blockTree == nil {
		return nil, fmt.Errorf("failed to parse with block grammar")
	}

	// return &MarkdownTree{
	// 	blockTree:     blockTree,
	// 	inlineTrees:   nil,
	// 	inlineIndices: nil,
	// }, nil

	// Process inline content
	inlineTrees := make([]*sitter.Tree, 0)
	inlineIndices := make(map[uintptr]*sitter.Tree)

	err = p.parser.SetLanguage(p.inlineLanguage)
	if err != nil {
		return nil, err
	}

	cursor := blockTree.Walk()
	p.processInlineNodes(cursor, content, &inlineTrees, inlineIndices)

	return &MarkdownTree{
		blockTree:     blockTree,
		inlineTrees:   inlineTrees,
		inlineIndices: inlineIndices,
	}, nil
}

// processInlineNodes processes all inline nodes in the tree
func (p *ExpMarkdownParser) processInlineNodes(cursor *sitter.TreeCursor, content []byte, inlineTrees *[]*sitter.Tree, inlineIndices map[uintptr]*sitter.Tree) {
	for {
		node := cursor.Node()
		kind := node.Kind()

		if kind == "inline" || kind == "pipe_table_cell" {
			// Create ranges for inline parsing, excluding child nodes
			ranges := p.createInlineRanges(node, content)
			if len(ranges) > 0 {
				err := p.parser.SetIncludedRanges(ranges)
				if err == nil {
					inlineTree := p.parser.Parse(content, nil)
					if inlineTree != nil {
						*inlineTrees = append(*inlineTrees, inlineTree)
						inlineIndices[node.Id()] = inlineTree
					}
				}
			}
		}

		if cursor.GotoFirstChild() {
			p.processInlineNodes(cursor, content, inlineTrees, inlineIndices)
			cursor.GotoParent()
		}

		if !cursor.GotoNextSibling() {
			break
		}
	}
}

// createInlineRanges creates ranges for inline parsing, excluding named child nodes
func (p *ExpMarkdownParser) createInlineRanges(node *sitter.Node, content []byte) []sitter.Range {
	ranges := make([]sitter.Range, 0)

	startByte := node.StartByte()
	endByte := node.EndByte()
	startPoint := node.StartPosition()
	endPoint := node.EndPosition()

	// Start with the full range
	currentRange := sitter.Range{
		StartByte:  startByte,
		EndByte:    endByte,
		StartPoint: startPoint,
		EndPoint:   endPoint,
	}

	// Find named children and exclude their ranges
	hasNamedChildren := false
	for i := uint(0); i < node.ChildCount(); i++ {
		child := node.Child(i)
		if child != nil && child.IsNamed() {
			hasNamedChildren = true
			// Add range before this child
			if currentRange.StartByte < child.StartByte() {
				ranges = append(ranges, sitter.Range{
					StartByte:  currentRange.StartByte,
					EndByte:    child.StartByte(),
					StartPoint: currentRange.StartPoint,
					EndPoint:   child.StartPosition(),
				})
			}
			// Move range start past this child
			currentRange.StartByte = child.EndByte()
			currentRange.StartPoint = child.EndPosition()
		}
	}

	// Add the final range (or the entire range if no named children)
	if currentRange.StartByte < currentRange.EndByte {
		ranges = append(ranges, currentRange)
	}

	// If no named children were found, include the entire range
	if !hasNamedChildren {
		ranges = []sitter.Range{{
			StartByte:  startByte,
			EndByte:    endByte,
			StartPoint: startPoint,
			EndPoint:   endPoint,
		}}
	}

	return ranges
}

// BlockTree returns the block tree
func (mt *MarkdownTree) BlockTree() *sitter.Tree {
	return mt.blockTree
}

// InlineTree returns the inline tree for a given node
func (mt *MarkdownTree) InlineTree(node *sitter.Node) *sitter.Tree {
	return mt.inlineIndices[node.Id()]
}

// InlineTrees returns all inline trees
func (mt *MarkdownTree) InlineTrees() []*sitter.Tree {
	return mt.inlineTrees
}

// Walk creates a new cursor for walking the tree
func (mt *MarkdownTree) Walk() *MarkdownCursor {
	return &MarkdownCursor{
		markdownTree: mt,
		blockCursor:  mt.blockTree.Walk(),
		inlineCursor: nil,
	}
}

// Node returns the current node
func (mc *MarkdownCursor) Node() *sitter.Node {
	if mc.inlineCursor != nil {
		return mc.inlineCursor.Node()
	}
	return mc.blockCursor.Node()
}

// IsInline returns true if the cursor is currently in an inline tree
func (mc *MarkdownCursor) IsInline() bool {
	return mc.inlineCursor != nil
}

// GoToFirstChild moves to the first child
func (mc *MarkdownCursor) GoToFirstChild() bool {
	if mc.inlineCursor != nil {
		return mc.inlineCursor.GotoFirstChild()
	}

	// Check if we can move into an inline tree
	node := mc.blockCursor.Node()
	if node.Kind() == "inline" || node.Kind() == "pipe_table_cell" {
		if inlineTree := mc.markdownTree.InlineTree(node); inlineTree != nil {
			mc.inlineCursor = inlineTree.Walk()
			return mc.inlineCursor.GotoFirstChild()
		}
	}

	return mc.blockCursor.GotoFirstChild()
}

// GoToParent moves to the parent node
func (mc *MarkdownCursor) GoToParent() bool {
	if mc.inlineCursor != nil {
		if mc.inlineCursor.GotoParent() {
			// Check if we're at the root of the inline tree
			if mc.inlineCursor.Node().Parent() == nil {
				mc.inlineCursor = nil
			}
			return true
		}
		mc.inlineCursor = nil
		return true
	}
	return mc.blockCursor.GotoParent()
}

// GoToNextSibling moves to the next sibling
func (mc *MarkdownCursor) GoToNextSibling() bool {
	if mc.inlineCursor != nil {
		return mc.inlineCursor.GotoNextSibling()
	}
	return mc.blockCursor.GotoNextSibling()
}
