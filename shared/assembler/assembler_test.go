package assembler

import (
	"testing"
	"fmt"
)

func Test_lex(t *testing.T) {
	files := []string{
		"./add_text.asm",
		"./add_data.asm",
	}

	for _, file := range files {
		fmt.Printf("lexing file: %s\n", file)
		tokens, err := lex(file)
		if err != nil {
			t.Fatalf("lexing error: %v", err)
		}
		fmt.Printf("lexed tokens: %+v\n", tokens)
	}
}
