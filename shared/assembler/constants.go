package assembler

import (
	"strings"
)

var grammarText string 

func init() {
	grammarText = strings.TrimSpace(`
lambda $ asm

asm data text $

data ".data dataList
data lambda

dataList dataItem dataList
dataList lambda

dataItem ident "= immediate

text ".text textList
text lambda

textList identifier ": textList
textList xInstruction textList
textList yInstruction textList
textList zInstruction textList
textList lambda

xInstruction xCommand register register

yInstruction yCommand register

zInstruction zCommand register ", zItem
zInstruction zjCommand mask ", zItem

zItem immediate
zItem ident
	`)
}
