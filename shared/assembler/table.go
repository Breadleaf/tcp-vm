package assembler

import (
	"fmt"
	"strings"
	"tcp-vm/shared/util"
)

type llParseTable struct {
	Data map[grammarItem]map[grammarItem][]grammarItem
}

func newLLParseTable(grammar grammar) (*llParseTable, error) {
	// verify that the grammar is LL(1)
	for nt, _ := range grammar.NonTerminals {
		if !grammar.verifyPredictPairwiseDisjoint(nt) {
			return nil, fmt.Errorf(
				"grammar is not LL(1): predict sets for"+
					" non-terminal %s are not disjoint",
				nt.Value,
			)
		}
	}

	llpt := llParseTable{
		Data: make(map[grammarItem]map[grammarItem][]grammarItem),
	}

	for lhs, rhs := range grammar.Rules {
		for _, rule := range rhs {
			follow := grammar.predictSet(lhs, rule)
			for t, _ := range follow {
				// create the map if it does not exist
				if _, ok := llpt.Data[lhs]; !ok {
					llpt.Data[lhs] = make(
						map[grammarItem][]grammarItem,
					)
				}
				// TODO: make sure not overwrite, voilate LL(1)
				llpt.Data[lhs][t] = rule
			}
		}
	}

	return &llpt, nil
}

func (llpt *llParseTable) prettyPrint() {
	counter := 0
	for nt, column := range llpt.Data {
		fmt.Printf("NonTerminal: %+v:\n", nt)
		for t, rule := range column {
			fmt.Printf("Terminal: %+v - Rule: %+v\n", t, rule)
		}
		if counter < len(llpt.Data)-1 {
			fmt.Println()
		}
		counter += 1
	}
}

func tokenMatches(tok token, gi grammarItem) bool {
	if gi.Value == tok.val {
		return true
	}

	typeName := strings.TrimPrefix(tok.typ.String(), "ttype.")
	if strings.EqualFold(typeName, gi.Value) {
		return true
	}

	return false
}

func (table *llParseTable) llTabularParse(
	tokens []token,
	start grammarItem,
) (*syntaxTree, error) {
	var ts util.Queue[token]
	for _, t := range tokens {
		ts.Push(t)
	}
	eof := token{val: "$", typ: Unknown, lin: 0}
	ts.Push(eof)

	root := newSyntaxTree(start, start.Value)
	current := root

	var stack util.Stack[grammarItem]
	stack.Push(start)

	for len(stack) > 0 {
		x, _ := stack.Pop()

		if x.Type == Terminal && x.Value == "$" {
			if _, err := ts.Pop(); err != nil {
				return nil, fmt.Errorf("unexpected end of input")
			}
			continue
		}

		// ascend marker
		if x == marker {
			if current.Parent != nil {
				current = current.Parent
			}
			continue
		}

		switch x.Type {
		case NonTerminal:
			peekTok, err := ts.Peek()
			if err != nil {
				return nil, fmt.Errorf("unexpected end of input")
			}

			col, exists := table.Data[x]
			if !exists {
				return nil, fmt.Errorf(
					"no LL(1) table entry for non-terminal %s",
					x.Value,
				)
			}

			var sel grammarItem
			matched := false
			for term := range col {
				if tokenMatches(peekTok, term) {
					sel = term
					matched = true
					break
				}
			}
			if !matched {
				return nil, fmt.Errorf(
					"no rule for %s when next token is %v",
					x.Value,
					peekTok,
				)
			}
			rule := col[sel]

			stack.Push(marker)
			for i := len(rule) - 1; i >= 0; i-- {
				stack.Push(rule[i])
			}

			node := current.addChild(x, x.Value)
			current = node

		case Terminal:
			if x.Type == Terminal {
				peekTok, err := ts.Peek()
				if err != nil {
					return nil, fmt.Errorf(
						"unexpected end of input",
					)
				}

				if !tokenMatches(peekTok, x) {
					return nil, fmt.Errorf(
						"expected '%s' but got '%s'",
						x.Value,
						peekTok.val,
					)
				}

				tok, _ := ts.Pop()
				current.addChild(x, tok.val)
			} else {
				// lambda (empty string)
				current.addChild(x, "")
			}

		default:
			// marker: ascend a node
			if x == marker && current.Parent != nil {
				current = current.Parent
			}
		}
	}

	return root, nil
}
