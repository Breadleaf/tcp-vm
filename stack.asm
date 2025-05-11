.text
main:
	LDI R0, 0x05
	PSH R0
	LDI R0, 0x06
	PSH R0
	POP R1
	SYS R1
