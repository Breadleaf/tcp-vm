package assembler

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"tcp-vm/shared/util"
)

type ttype int

const (
	Section ttype = iota
	CommandX
	CommandY
	CommandZ
	CommandZJ
	Register
	Mask
	Immediate
	Identifier
	Comma
	Equals
	Colon
	Unknown
)

func (tt ttype) String() string {
	switch tt {
	case Section:
		return "ttype.Section"
	case CommandX:
		return "ttype.CommandX"
	case CommandY:
		return "ttype.CommandY"
	case CommandZ:
		return "ttype.CommandZ"
	case CommandZJ:
		return "ttype.CommandZJ"
	case Register:
		return "ttype.Register"
	case Mask:
		return "ttype.Mask"
	case Immediate:
		return "ttype.Immediate"
	case Identifier:
		return "ttype.Identifier"
	case Comma:
		return "ttype.Comma"
	case Equals:
		return "ttype.Equals"
	case Colon:
		return "ttype.Colon"
	default:
		return "ttype.Unknown"
	}
}

type token struct {
	val string
	typ ttype
	lin int
}

func (t token) String() string {
	return fmt.Sprintf("token('%s'/%v/%d)", t.val, t.typ, t.lin)
}

func Assemble(sourcePath string) ([]byte, []byte, error) {
	logTag := "tcp-vm/shared/assembler - assembler.go - Assemble()"
	util.LogStart(logTag)
	defer util.LogEnd(logTag)

	e := []byte{}

	tokens, err := lex("./add_data.asm")
	if err != nil {
		return e, e, fmt.Errorf("lex() failed: %v", err)
	}

	g, err := newGrammar()
	if err != nil {
		return e, e, fmt.Errorf("newGrammar() failed: %v", err)
	}

	llpt, err := newLLParseTable(*g)
	if err != nil {
		return e, e, fmt.Errorf("newLLParseTable() failed: %v", err)
	}

	start := grammarItem{
		Value: "asm",
		Type:  NonTerminal,
	}

	st, err := llpt.llTabularParse(tokens, start)
	if err != nil {
		return e, e, fmt.Errorf("llTabularParse() failed: %v", err)
	}

	util.LogMessage(func() {
		fmt.Println("CST:")
		st.prettyPrint()
	})

	if simp := st.applySDT(); simp != nil {
		util.LogMessage(func() {
			fmt.Println("AST:")
			simp.prettyPrint()
		})

		data, text, err := simp.compile()
		if err != nil {
			return e, e, fmt.Errorf("compile() filed: %v", err)
		}

		return data, text, nil
	}

	return e, e, fmt.Errorf("appltSDT() returned nil")
}

func lex(sourcePath string) ([]token, error) {
	tokenSpecs := map[ttype]string{
		Section:    `\..*`,
		CommandX:   `(MOV|CMP|SHL|SHR|ADD|SUB|AND|ORR)`,
		CommandY:   `(NOT|PSH|POP|SYS)`,
		CommandZ:   `(LDI|LDA|STA)`,
		CommandZJ:  `(JMP)`,
		Register:   `(R\d|PC|SP)`,
		Mask:       `[0-1]{3}`,
		Immediate:  `0x[A-F0-9]{2}`,
		Identifier: `[a-z]*`,
		Comma:      `,`,
		Equals:     `=`,
		Colon:      `:`,
	}

	type spec struct {
		typ ttype
		re  *regexp.Regexp
	}

	var specs []spec
	for tt, raw := range tokenSpecs {
		re, err := regexp.Compile("^(" + raw + ")")
		if err != nil {
			return nil, fmt.Errorf(
				"invalid regex for %v: %v",
				tt,
				raw,
			)
		}
		specs = append(specs, spec{
			typ: tt,
			re:  re,
		})
	}

	f, err := os.Open(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("opening file: %v", err)
	}
	defer f.Close()

	var tokens []token
	scanner := bufio.NewScanner(f)
	line_number := 0

	// TODO: make a more robust loop
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		line_number += 1

		// convert comment lines to empty and remove inline comments
		if idx := strings.Index(line, "#"); idx != -1 {
			line = line[:idx]
		}

		// remove empty lines
		if line == "" {
			continue
		}

		pos := 0
		for pos < len(line) {
			// skip spaces or tabs
			if c := line[pos]; c == ' ' || c == '\t' {
				pos += 1
				continue
			}

			sub := line[pos:]
			best := ""
			bestType := Unknown

			for _, s := range specs {
				if m := s.re.FindString(sub); len(m) > len(best) {
					best = m
					bestType = s.typ
				}
			}

			// fall back to singe char unknown
			if best == "" {
				best = sub[:1]
				bestType = Unknown
			}

			tokens = append(tokens, token{
				val: best,
				typ: bestType,
				lin: line_number,
			})

			pos += len(best)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading file: %v", err)
	}

	return tokens, nil
}
