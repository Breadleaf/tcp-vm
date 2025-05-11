.text
main:
	# setup stack

	LDI R0, 0x01
	PSH R0

	LDI R0, 0x02
	PSH R0

	LDI R0, 0x03
	PSH R0

	LDI R0, 0x03
	PSH R0

	# stack:
	# top -> [3, 3, 2, 1] <- bottom

	# unconditional jump to adder
	CMP R0 R0
	JMP 010, adder

adderreturn:
	# get results from adder
	POP R1

	# exit with the result from adder as the exit status
	PSH R1
	# known bug: if not sys_exit or sys_sleep something crazy happens with
	# the flag
	LDI R0, 0x00 # sys_exit
	SYS R0

# adder(argc, argv)

adder:
	# use R1 as a counter
	POP R1
	# use R0 as the sum (init to 0)
	LDI R0, 0x00

adderloopstart:
	# store sum
	PSH R0
	# see if loop is done
	LDI R0, 0x00
	CMP R1 R0
	JMP 010, adderloopend
	# restore sum
	POP R0

	# store counter
	PSH R1
	# get next argument
	POP R1
	# sum += arg
	ADD R0 R1
	# restore counter
	POP R1

	# store sum
	PSH R0
	# decrement counter
	LDI R0, 0x01
	SUB R1 R0
	# restore sum
	POP R0

	# unconditional jump to adderloopstart
	CMP R0 R0
	JMP 010, adderloopstart

adderloopend:
	# return value
	PSH R0
	# unconditional jump to adderreturn
	CMP R0 R0
	JMP 010, adderreturn
