.data
x = 0x05
y = 0x02

.text
main:
	# R0 = x
	LDA R0, x

	# negate R0
	LDI R1, afternegate
	PSH R1

	# goto negate
	CMP R0 R0
	JMP 010, negate

afternegate:
	# R1 = R0
	MOV R1 R0

	# R0 = y
	LDA R0, y

	# R0 -= R1
	SUB R0 R1

	# exit sys_exit(R0)
	PSH R0
	LDI R0, 0x00
	SYS R0

# R0 = negate(R0): (~R0)+1
# clobbers: R1
# notes:
# - returns using `POP PC`
negate:
	# (~R0)+1
	NOT R0
	LDI R1, 0x01
	ADD R0 R1

	# return to return address
	POP PC
