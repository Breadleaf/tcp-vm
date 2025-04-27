.text
main:
	# R0 = 5
	LDI R0, 0x05
	# R1 = 200
	LDI R1, 0xC8
	# argv[0] = R1 + R0
	ADD R1 R0
	PSH R1
	# SYS_EXIT
	LDI R0, 0x00
	SYS R0
