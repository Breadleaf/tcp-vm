package assembler

import (
	"fmt"
)

type llParseTable struct {
	Data map[grammarItem]map[grammarItem][]grammarItem
}

func newLLParseTable(grammar grammar) *llParseTable {
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

	return &llpt
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
