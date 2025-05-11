package assembler

import (
	"fmt"
	"strconv"
	"strings"
	g "tcp-vm/shared/globals"
	"tcp-vm/shared/vm"
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

func (st *syntaxTree) applySDT() *syntaxTree {
	// simplify the children
	var newChildren []*syntaxTree
	for _, child := range st.Children {
		if simplified := child.applySDT(); simplified != nil {
			// ignore punctuation terminals
			if simplified.Symbol.Type == Terminal {
				sym := simplified.Data
				if sym == "=" || sym == ":" || sym == "," {
					continue
				}
			}

			simplified.Parent = st
			newChildren = append(newChildren, simplified)
		}
	}
	st.Children = newChildren

	// remove lambda/epsilon nodes
	if st.Symbol.Type == Lambda {
		return nil
	}

	// flatten recursive list nodes
	if st.Symbol.Value == "textList" || st.Symbol.Value == "dataList" {
		var flat []*syntaxTree
		for _, child := range st.Children {
			if child.Symbol.Value == st.Symbol.Value {
				// append grand children
				for _, gc := range child.Children {
					gc.Parent = st
					flat = append(flat, gc)
				}
			} else {
				child.Parent = st
				flat = append(flat, child)
			}
		}
		st.Children = flat
	}

	// collapse unary productions
	if len(st.Children) == 1 {
		only := st.Children[0]
		only.Parent = st.Parent
		return only
	}

	return st
}

func (st *syntaxTree) compile() ([g.DataSectionLength]byte, [g.TextSectionLength]byte, error) {
	dataLabels := map[string]uint8{}
	textLabels := map[string]uint8{}
	var dataSection []uint8
	var textSection []uint8

	// Data pass: collect dataItem under dataList
	for _, sec := range st.Children {
		if sec.Symbol.Value != "data" {
			continue
		}
		// find dataList under data
		var items []*syntaxTree
		for _, c := range sec.Children {
			if c.Symbol.Value == "dataList" {
				items = c.Children
				break
			}
		}
		for _, item := range items {
			if item.Symbol.Value != "dataItem" {
				continue
			}
			if len(dataSection) >= g.DataSectionLength {
				return ErrorData, ErrorText, fmt.Errorf("data section overflow: exceeds %d words", g.DataSectionLength)
			}
			label := item.Children[0].Data
			if prev, dup := dataLabels[label]; dup {
				return ErrorData, ErrorText, fmt.Errorf("duplicate data label '%s' at address %d", label, prev)
			}
			dataLabels[label] = uint8(len(dataSection)) + vm.DataStart

			// identifier at [0], immediate at [1]
			lit := item.Children[1].Data
			val, err := parseImmediate(lit)
			if err != nil {
				return ErrorData, ErrorText, err
			}
			dataSection = append(dataSection, val)
		}
	}

	// Text pass: collect instrs under the single textList
	var instrs []*syntaxTree
	for _, sec := range st.Children {
		if sec.Symbol.Value != "text" {
			continue
		}
		for _, c := range sec.Children {
			if c.Symbol.Value == "textList" {
				instrs = c.Children
				break
			}
		}
	}
	// assign label addresses
	var addr uint8
	for _, node := range instrs {
		switch node.Symbol.Value {
		case "identifier":
			lbl := node.Data
			if _, dup := textLabels[lbl]; dup {
				return ErrorData, ErrorText, fmt.Errorf("duplicate text label '%s'", lbl)
			}
			textLabels[lbl] = addr + vm.TextStart
		case "xInstruction", "yInstruction":
			if addr >= g.TextSectionLength {
				return ErrorData, ErrorText, fmt.Errorf("text section overflow: exceeds %d words", g.TextSectionLength)
			}
			addr++
		case "zInstruction":
			if addr+1 >= g.TextSectionLength {
				return ErrorData, ErrorText, fmt.Errorf("text section overflow: exceeds %d words", g.TextSectionLength)
			}
			addr += 2
		}
	}
	// this is technically not needed, I will enforce it for good code practice
	if _, ok := textLabels["main"]; !ok {
		return ErrorData, ErrorText, fmt.Errorf("missing 'main' label in text section")
	}

	// Emit code over instrs
	for _, node := range instrs {
		switch node.Symbol.Value {
		case "xInstruction":
			op := node.Children[0].Data
			b, err := compileX(op, node.Children[1:])
			if err != nil {
				return ErrorData, ErrorText, err
			}
			textSection = append(textSection, b)

		case "yInstruction":
			op := node.Children[0].Data
			b, err := compileY(op, node.Children[1:])
			if err != nil {
				return ErrorData, ErrorText, err
			}
			textSection = append(textSection, b)

		case "zInstruction":
			op := node.Children[0].Data
			args := node.Children[1:]
			if op == "JMP" {
				b, imm, err := compileZJ(op, args, textLabels)
				if err != nil {
					return ErrorData, ErrorText, err
				}
				textSection = append(textSection, b, imm)
			} else {
				b, imm, err := compileZ(op, args, dataLabels)
				if err != nil {
					return ErrorData, ErrorText, err
				}
				textSection = append(textSection, b, imm)
			}
		}
	}

	dataOut := [g.DataSectionLength]byte{}
	for i, v := range dataSection {
		dataOut[i] = byte(v)
	}
	textOut := [g.TextSectionLength]byte{}
	for i, v := range textSection {
		textOut[i] = byte(v)
	}

	fmt.Println("data labels:")
	for k, v := range dataLabels {
		fmt.Printf("%v, %v\n", k, v)
	}
	fmt.Println("text labels:")
	for k, v := range textLabels {
		fmt.Printf("%v, %v\n", k, v)
	}

	return dataOut, textOut, nil
}

func parseImmediate(lit string) (uint8, error) {
	var v uint64
	var err error
	if strings.HasPrefix(lit, "0x") {
		v, err = strconv.ParseUint(lit[2:], 16, 8)
	}
	if err != nil {
		return 0, fmt.Errorf("invalid immediate '%s': %v", lit, err)
	}
	return uint8(v), nil
}

func parseRegister(reg string) (byte, error) {
	switch reg {
	case "R0":
		return 0, nil
	case "R1":
		return 1, nil
	case "SP":
		return 2, nil
	case "PC":
		return 3, nil
	default:
		return 0, fmt.Errorf("invalid register '%s'", reg)
	}
}

func compileX(op string, args []*syntaxTree) (byte, error) {
	// o4 o3 o2 o1 ra1 ra0 rb1 rb0

	// prevent future errors if changes to the parser are made
	if len(args) != 2 {
		return 0, fmt.Errorf(
			"compileX %s: expected 2 operands, got %d",
			op,
			len(args),
		)
	}

	var code byte
	switch op {
	case "MOV":
		code = 0
	case "CMP":
		code = 1
	case "SHL":
		code = 2
	case "SHR":
		code = 3
	case "ADD":
		code = 4
	case "SUB":
		code = 5
	case "AND":
		code = 6
	case "ORR":
		code = 7
	default:
		return 0, fmt.Errorf("invalid X-type opcode: '%s'", op)
	}

	// top nibble
	inst := code << 4

	// registers
	ra, err := parseRegister(args[0].Data)
	if err != nil {
		return 0, err
	}
	rb, err := parseRegister(args[1].Data)
	if err != nil {
		return 0, err
	}
	inst |= (ra<<2 | rb)

	return inst, nil
}

func compileY(op string, args []*syntaxTree) (byte, error) {
	// o4 o3 o2 o1 0 0 ra1 ra0

	// prevent future errors if changes to the parser are made
	if len(args) != 1 {
		return 0, fmt.Errorf(
			"compileY %s: expected 1 operands, got %d",
			op,
			len(args),
		)
	}

	var code byte
	switch op {
	case "NOT":
		code = 8
	case "PSH":
		code = 9
	case "POP":
		code = 0xA
	case "SYS":
		code = 0xB
	default:
		return 0, fmt.Errorf("invalid Y opcode '%s'", op)
	}

	inst := code << 4
	ra, err := parseRegister(args[0].Data)
	if err != nil {
		return 0, err
	}
	inst |= ra
	return inst, nil
}

func compileZ(op string, args []*syntaxTree, labels map[string]uint8) (byte, byte, error) {
	// o4 o3 o2 o1 0 0 ra1 ra0, immediate

	// prevent future errors if changes to the parser are made
	if len(args) != 2 {
		return 0, 0, fmt.Errorf(
			"compileZ %s: expected 2 operands, got %d",
			op,
			len(args),
		)
	}

	var code byte
	switch op {
	case "LDI":
		code = 0xD
	case "LDA":
		code = 0xE
	case "STA":
		code = 0xF
	default:
		return 0, 0, fmt.Errorf("invalid Z opcode '%s'", op)
	}
	inst := code << 4
	// register
	ra, err := parseRegister(args[0].Data)
	if err != nil {
		return 0, 0, err
	}
	inst |= ra << 2
	// immediate or label
	valStr := args[1].Data
	var imm byte
	if addr, ok := labels[valStr]; ok {
		imm = addr
	} else {
		tmp, err := parseImmediate(valStr)
		if err != nil {
			return 0, 0, err
		}
		imm = tmp
	}
	return inst, imm, nil
}

func compileZJ(op string, args []*syntaxTree, labels map[string]uint8) (byte, byte, error) {
	// o4 o3 o2 o1 0 m2 m1 m0, immediate

	// prevent future errors if changes to the parser are made
	if len(args) != 2 {
		return 0, 0, fmt.Errorf(
			"compileZJ %s: expected 2 operands, got %d",
			op,
			len(args),
		)
	}

	// only JMP is allowed
	if op != "JMP" {
		return 0, 0, fmt.Errorf("invalid ZJ opcode '%s'", op)
	}
	// opcode nibble
	inst := byte(0xC) << 4
	// mask bits
	maskStr := args[0].Data
	m, err := strconv.ParseUint(maskStr, 2, 3)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid mask '%s': %v", maskStr, err)
	}
	inst |= byte(m)
	// label
	lbl := args[1].Data
	addr, ok := labels[lbl]
	if !ok {
		return 0, 0, fmt.Errorf("undefined jump label '%s'", lbl)
	}
	return inst, addr, nil
}
