.text
main:
	# load 5 into R0
	LDI R0, 0x05

	# push return address
	LDI R1, firstnegatereturn
	PSH R1

	# jump to negate
	CMP R0 R0
	JMP 010, negate

firstnegatereturn:
	PSH R0
	LDI R0, 0x00
	SYS R0

# negate:
#  inputs:
#   R0 = some 8 bit value X
#  outputs:
#   R0 = two's complement of X
#  clobbers:
#   R1
#  notes:
#   uses stack to return via `POP PC`
negate:
	NOT R0
	LDI R1, 0x01
	ADD R0 R1
	POP PC
