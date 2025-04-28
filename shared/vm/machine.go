package vm

import (
	"fmt"
)

const FILE_LOG_TAG = "tcp-vm/shared/vm/vm.go"

// memory size calculations

const (
	// see spec sheet
	vmMemSizeWords = 16 + 64 + 1 + 175
	vmMemSizeBytes = vmMemSizeWords // 1 word : 1 byte
)

// memory partitioning

const (
	vmDataStart = 0
	vmDataCount = 16
	vmDataEnd   = vmDataStart + vmDataCount - 1
)

const (
	vmStackStart = vmDataEnd + 1
	vmStackCount = 64
	vmStackEnd   = vmStackStart + vmStackCount - 1
)

const (
	vmFlagStart = vmStackEnd + 1
	vmFlagCount = 1
	vmFlagEnd   = vmFlagStart + vmFlagCount - 1
)

const (
	vmTextStart = vmFlagEnd + 1
	vmTextCount = 175
	vmTextEnd   = vmTextStart + vmTextCount - 1
)

const (
	MemoryStart = vmDataStart
	MemoryEnd   = vmTextEnd
)

const (
	DataStart  = vmDataStart
	StackStart = vmStackStart
	FlagStart  = vmFlagStart
	TextStart  = vmTextStart
)

// hardware

type Register byte

type Memory [vmMemSizeWords]byte

// virtual machine

type VirtualMachine struct {
	R0     Register
	R1     Register
	SP     Register
	PC     Register
	Memory Memory
	Output string
}

func (vm *VirtualMachine) ResetFromStateless(
	data [vmDataCount]byte,
	text [vmTextCount]byte,
) {
	vm.R0 = Register(0)
	vm.R1 = Register(0)
	vm.SP = Register(0)
	vm.PC = Register(0)

	copy(vm.Memory[vmDataStart:vmDataEnd+1], data[:]) // end is exclusive

	for i := vmDataEnd + 1; i <= vmStackEnd; i++ {
		vm.Memory[i] = 0
	}

	vm.Memory[vmFlagStart] = 0

	copy(vm.Memory[vmTextStart:vmTextEnd+1], text[:]) // end is exclusive

	vm.Output = ""
}

func (vm *VirtualMachine) ResetFromStateful(
	r0 byte,
	r1 byte,
	sp byte,
	pc byte,
	data [vmDataCount]byte,
	stack [vmStackCount]byte,
	flag [vmFlagCount]byte,
	text [vmTextCount]byte,
) {
	vm.R0 = Register(r0)
	vm.R1 = Register(r1)
	vm.SP = Register(sp)
	vm.PC = Register(pc)

	copy(vm.Memory[vmDataStart:vmDataEnd+1], data[:]) // end is exclusive
	copy(vm.Memory[vmStackStart:vmStackEnd+1], stack[:]) // end is exclusive
	copy(vm.Memory[vmFlagStart:vmFlagEnd+1], flag[:]) // end is exclusive
	copy(vm.Memory[vmTextStart:vmTextEnd+1], text[:]) // end is exclusive

	vm.Output = ""
}

func (vm *VirtualMachine) String() string {
	out := fmt.Sprintf(
		"R0: %d, R1: %d, SP: %d, PC: %d\n",
		vm.R0, vm.R1, vm.SP, vm.PC,
	)

	for idx, byt := range vm.Memory {
		// print section labels
		switch idx {
		case vmDataStart:
			out += "Data Section:\n"
		case vmStackStart:
			out += "Stack Section:\n"
		case vmFlagStart:
			out += "Flag Section:\n"
		case vmTextStart:
			out += "Text Section:\n"
		}

		out += fmt.Sprintf("%3d: %08b\n", idx, byt)
	}

	out += fmt.Sprintf("Output:\n%s\n------\n", vm.Output)

	return out
}
