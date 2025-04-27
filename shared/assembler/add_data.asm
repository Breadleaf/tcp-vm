.data
x = 0x05
y = 0xC8

.text
main:
	# R0 = x
	LDI R0, x
	# R1 = y
	LDI R1, y
	# argv[0] = R1 + R0
	ADD R1 R0
	PSH R1
	# SYS_EXIT
	LDI R0, 0x00
	SYS R0
