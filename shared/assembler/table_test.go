package assembler

import (
	"testing"
)

func Test_tableCreation(t *testing.T) {
	g, err := newGrammar()
	if err != nil {
		t.Errorf("newGrammar() failed: %v", err)
	}

	llpt := newLLParseTable(*g)
	llpt.prettyPrint()
}
