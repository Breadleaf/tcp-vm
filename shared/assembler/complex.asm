.data
i = 0x01
target = 0x00

.text
main:
	# R0=i, R1=target
	LDA R0, i
	LDA R1, target

	# if R0 = R1 then goto end
	CMP R0 R1
	JMP 010, end

	# else R0+=1
	LDI R1, 0x01
	ADD R0 R1
	STA R0, i

	# uncond jump to main
	LDI R0, 0x00
	LDI R1, 0x00
	CMP R0 R1
	JMP 010, main

end:
	# load sys 255
	LDI R0, 0xFF
	SYS R0
