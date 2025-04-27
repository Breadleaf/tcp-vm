package assembler

import (
	"testing"
)

func Test_tableCreation(t *testing.T) {
	g, err := newGrammar()
	if err != nil {
		t.Fatalf("newGrammar() failed: %v", err)
	}

	llpt, err := newLLParseTable(*g)
	if err != nil {
		t.Fatalf("newLLParseTable() failed: %v", err)
	}
	llpt.prettyPrint()
}

func Test_llTabularParse(t *testing.T) {
	tokens, err := lex("./add_data.asm")
	if err != nil {
		t.Fatalf("lex() failed: %v", err)
	}

	g, err := newGrammar()
	if err != nil {
		t.Fatalf("newGrammar() failed: %v", err)
	}

	llpt, err := newLLParseTable(*g)
	if err != nil {
		t.Fatalf("newLLParseTable() failed: %v", err)
	}

	start := grammarItem{
		Value: "asm",
		Type:  NonTerminal,
	}

	st, err := llpt.llTabularParse(tokens, start)
	if err != nil {
		t.Fatalf("llTabularParse() failed: %v", err)
	}

	st.prettyPrint()

	if simp := st.applySDT(); simp != nil {
		simp.prettyPrint()
		data, text, err := simp.compile()
		if err != nil {
			t.Fatalf("compile() filed: %v", err)
		}
		t.Logf("data section:\n")
		for _, b := range data {
			t.Logf("%+v", b)
		}
		t.Logf("text section:\n")
		for _, b := range text {
			t.Logf("%+v", b)
		}
	}
}
