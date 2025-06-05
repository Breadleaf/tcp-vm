.data
first = 0x12
second = 0x34

.text
main:
    LDA   R0, first      # R0 ← Memory[first]
    PSH   R0             # push that byte

    LDI   R0, 0x12       # R0 ← 0x12 (expected high‐byte)
    POP   R1             # R1 ← (original Memory[first])

    CMP   R1 R0          # flag = EQ if (R1 == 0x12)
    JMP   010, big       # if EQ, jump to big

little:
    LDI   R0, 0x02       # R0 = 0x02 (return code for little‐endian)
    PSH   R0             # push return code
    LDI   R0, 0x00       # R0 = 0 (sys_exit)
    SYS   R0             # pop arg → R0 = 0x02, halt

big:
    LDI   R0, 0x01       # R0 = 0x01 (return code for big‐endian)
    PSH   R0             # push return code
    LDI   R0, 0x00       # R0 = 0 (sys_exit)
    SYS   R0             # pop arg → R0 = 0x01, halt
