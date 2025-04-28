package assembler

import (
	"os"
	"strings"
	g "tcp-vm/shared/globals"
)

// used across the assembler to fulfill the return types of functions that are
// returning errors but need to send a []byte of a size
var (
	ErrorData = [g.DataSectionLength]byte{}
	ErrorText = [g.TextSectionLength]byte{}
)

var LOG_PARSED_GRAMMAR_OBJECT bool

var grammarText string

func init() {
	LOG_PARSED_GRAMMAR_OBJECT = os.Getenv("LOG_PARSED") != ""

	grammarText = strings.TrimSpace(`
asm -> data text $

data -> .data dataList
data -> lambda

dataList -> dataItem dataList
dataList -> lambda

dataItem -> identifier = immediate

text -> .text textList
text -> lambda

textList -> identifier : textList
textList -> xInstruction textList
textList -> yInstruction textList
textList -> zInstruction textList
textList -> lambda

xInstruction -> CommandX register register

yInstruction -> CommandY register

zInstruction -> CommandZ register , zItem
zInstruction -> CommandZJ mask , zItem

zItem -> immediate
zItem -> identifier
	`)
}
