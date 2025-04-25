# Server - VM

This package contains the virtual machine. This file serves the goal of
documenting the specifications of the virtual machine.

## Table of Contents

1. [General Overview](#general-overview)
2. [Registers](#registers)
3. [ISA](#isa)
4. [Memory](#memory)
5. [Limitations](#limitations)

## General Overview

## Registers

There are 4 8-bit registers. Since we know 2^2=4 we can address each register
with just 2-bits. 3 Of these registers are general purpose `R0-R3` and the last
one is the pointer to the top of the stack `SP`.

## ISA

Instructions are also 8-bit or one word in the machine.

Each instruction is in the shape of:
| MSB | Opcode | Ra | Rb | LSB |
| :-: | :-: | :-: | :-: | :-: |
| | NNNN | NN | NN | |

Since this means the opcode is 4-bits wide, we get 2^4=16 possible
instructions. This is limiting but possible. One such limitation is that the
first register Ra will also be used as the register result for all
instructions. For example ADD R0 R1 would be the same as R0 = R0 + R1.

| Opcode (Binary) | Opcode (Hex) | Instruction | Effect |
| :-: | :-: | :-: | :-: |
| 0000 | 0 | MOV Ra Rb | Ra = Rb |
| 0001 | 1 | ADD Ra Rb | Ra = Ra + Rb |
| 0010 | 2 | SUB Ra Rb | Ra = Ra - Rb |
| 0011 | 3 | AND Ra Rb | Ra = Ra & Rb |
| 0100 | 4 | ORR Ra Rb | Ra = Ra \| Rb |
| 0101 | 5 | NOT Ra | Ra = ~Ra |
| 0110 | 6 | SHL Ra Rb | Ra = Ra << Rb |
| 0111 | 7 | SHR Ra Rb | Ra = Ra >> Rb |
| 1000 | 8 | LDI Ra, Imm | R0 = Imm |
| 1001 | 9 | PSH Ra | Memory\[SP\] = Ra, SP = SP - 1 |
| 1010 | A | POP Rb | SP = SP + 1, Rb = Memory\[SP\] |
| 1011 | B | CMP Ra Rb | sets flags based on: Ra - Rb |
| 1100 | C | JEQ | if eq flag: PC = R0 |
| 1101 | D | LOD Ra Rb | Ra = Memory\[Rb\] |
| 1110 | E | STR Ra Rb | Memory\[Rb\] = Ra |
| 1111 | F | HLT | stop execution |

Note: Instructions with a comma in their `Instruction` will use 2 words. This
way they can operate using 8-bit immediate values which is also the maximum
value I plan to implement on this project.

## Memory

The entire system lives in the 4 registers as well as 256 words of ram. This
equates to 256-Bytes of ram. The first 32 words will be the `.data` section,
where all global variables will live. The next 32 words will be the stack,
which will start at `0x32` and grow down when you `PSH` and grow up when you
`POP`. Finally the last 192 words will be allocated to the `.text` section.
Here is where all instructions will live.

## Limitations
