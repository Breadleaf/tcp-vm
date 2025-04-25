# Server - VM

This package contains the virtual machine. This file serves the goal of
documenting the specifications of the virtual machine.

## Table of Contents

1. [General Overview](#general-overview)
2. [Registers](#registers)
3. [ISA](#isa)
4. [Limitations](#limitations)

## General Overview

## Registers

There are 4 8-bit registers. Since we know `2^2=4` we can address each register
with just 2-bits.

## ISA

Instructions are also 8-bit.

Each instruction is in the shape of:
| MSB | Opcode | Ra | Rb | LSB |
| :-: | :-: | :-: | :-: | :-: |
| | NNNN | NN | NN | |

Since this means the opcode is 4-bits wide, we get `2^4=16` possible
instructions. This is limiting but possible. One such limitation is that the
first register `Ra` will also be used as the register result for all
instructions. For example `ADD R0 R1` would be the same as `R0 = R0 + R1`.

| Opcode (Binary) | Opcode (Hex) | Instruction | Effect |
| :-: | :-: | :-: | :-: |
| 0000 | 0 | MOV Ra Rb | Ra = Rb |
| 0001 | 1 | ADD Ra Rb | Ra = Ra + Rb |
| 0010 | 2 | SUB Ra Rb | Ra = Ra - Rb |
| 0011 | 3 | AND Ra Rb | Ra = Ra & Rb |
| 0100 | 4 | OR Ra Rb | Ra = Ra \| Rb |
| 0101 | 5 | Not Ra | Ra = ~Ra |
| 0110 | 6 | SHL Ra Rb | Ra = Ra << Rb |
| 0111 | 7 | SHR Ra Rb | Ra = Ra >> Rb |
| 1000 | 8 | LDI Ra Imm | Ra = Imm |
| 1001 | 9 | CMP Ra Rb | sets flags based on: Ra - Rb |
| 1010 | A | JEQ Addr | if eq flag: PC = Addr |
| 1011 | B | JNE Addr | if not eq flag: PC = Addr |
| 1100 | C | JMP Addr | PC = Addr |
| 1101 | D | LOAD Ra Rb | Ra = Memory\[Rb\] |
| 1110 | E | STOR Ra Rb | Memory\[Rb\] = Ra |
| 1111 | F | HALT | stop execution |

## Limitations
