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
with just 2-bits. 2 Of these registers are general purpose `R0-R1`, the next
register is the program counter `PC`, and the last one is the pointer to the
top of the stack `SP`.

## ISA

Instructions are also 8-bit or one word in the machine.

Each instruction is in the shape of:
| MSB | Opcode | Ra | Rb | LSB |
| :-: | :-: | :-: | :-: | :-: |
| | NNNN | NN | NN | |

Since this means the opcode is 4-bits wide, we get 2^4=16 possible
instructions. This is limiting but possible. One such limitation is that the
first register Ra will also be used as the register result for all
instructions. For example `ADD R0 R1` would be the same as `R0 = R0 + R1`.

Instructions:
| Opcode (Binary) | Instruction | Effect |
| :-: | :-: | :-: |
| 0000 | MOV Ra Rb | Ra = Rb |
| 0001 | ADD Ra Rb | Ra = Ra + Rb |
| 0010 | SUB Ra Rb | Ra = Ra - Rb |
| 0011 | AND Ra Rb | Ra = Ra & Rb |
| 0100 | ORR Ra Rb | Ra = Ra \| Rb |
| 0101 | NOT Ra | Ra = ~Ra |
| 0110 | SHL Ra Rb | Ra = Ra << Rb |
| 0111 | SHR Ra Rb | Ra = Ra >> Rb |

| 1000 | PSH Ra |
| 1001 | POP

## Syscalls

The `SYS` opcode will expect `R0` to hold the syscall number. The stack will
contain all the arguments to the function.

Syscalls:
| Hex Code | Name | Stack Args | Effect |
| :-: | :-: | :-: |
| 0 | SYS_EXIT | \[SP\]: exit code | exit(\[SP\]) |
| 1 | SYS_OUT | \[SP\]: number of char <br> \[SP-1..N\]: remainder of arguments (WARNING: stack overflow) | ... |
| 2 | SYS_RAND | NONE | rand(0, 255) |
| 3 | SYS_WAIT | \[SP\]: number of cycles to wait |

## Memory

The entire system lives in the 4 registers as well as 256 words of ram. This
equates to 256-Bytes of ram. The first 32 words will be the `.data` section,
where all global variables will live. The next 32 words will be the stack,
which will start at `0x32` and grow down when you `PSH` and grow up when you
`POP`. Finally the last 192 words will be allocated to the `.text` section.
Here is where all instructions will live.

## Limitations
