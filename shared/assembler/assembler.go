package assembler

import (
	"strings"
	"os"
	"bufio"
	"fmt"
	"regexp"
)

type ttype int

const (
	Section ttype = iota
	Label
	Instruction
	Register
	Immediate
	Comma
	Equals
	Unknown
)

func (tt ttype) String() string {
	switch tt {
	case Section:
		return "ttype.Section"
	case Label:
		return "ttype.Label"
	case Instruction:
		return "ttype.Instruction"
	case Register:
		return "ttype.Register"
	case Immediate:
		return "ttype.Immediate"
	case Comma:
		return "ttype.Comma"
	case Equals:
		return "ttype.Equals"
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

func Assemble(sourcePath string) error {
	return nil
}

func lex(sourcePath string) ([]token, error) {
	tokenSpecs := map[ttype]string{
		Section: `\..*`,
		Label: `.*:`,
		Instruction: `[A-Z]*`,
		Register: `(R\d|PC|SP)`,
		Immediate: `0x[A-F0-9]{2}`,
		Comma: `,`,
		Equals: `=`,
	}

	type spec struct {
		typ ttype
		re *regexp.Regexp
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
			re: re,
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
