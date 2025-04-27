package assembler

import (
	"fmt"
	"strings"
)

type syntaxTree struct {
	Parent   *syntaxTree
	Children []*syntaxTree
	Data     string
	Symbol   grammarItem
}

var marker = grammarItem{
	Type:  Terminal,
	Value: "<MARKER>",
}

func newSyntaxTree(sym grammarItem, data string) *syntaxTree {
	st := &syntaxTree{
		Parent:   nil,
		Children: []*syntaxTree{},
		Data:     data,
		Symbol:   sym,
	}
	return st
}

func (st *syntaxTree) prettyPrint() {
	getIndent := func(indent int) string {
		return strings.Repeat("  ", indent)
	}

	var helper func(*syntaxTree, int)
	helper = func(node *syntaxTree, indent int) {
		fmt.Printf(
			"%sSymbol: %s, Data: %s\n",
			getIndent(indent),
			node.Symbol.Value,
			node.Data,
		)

		fmt.Printf("%schildren:\n", getIndent(indent))
		for _, child := range node.Children {
			helper(child, indent+1)
		}
	}

	helper(st, 0)
}

func (st *syntaxTree) addChild(sym grammarItem, data string) *syntaxTree {
	child := newSyntaxTree(sym, data)
	child.Parent = st
	st.Children = append(st.Children, child)
	return child
}
