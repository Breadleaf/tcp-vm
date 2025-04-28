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
	vmMemStart = vmDataStart
	vmMemEnd   = vmTextEnd
)

// hardware

type Register uint8

type Memory [vmMemSizeWords]uint8

// virtual machine

type VirtualMachine struct {
	R0     Register
	R1     Register
	SP     Register
	PC     Register
	Memory Memory
}

func NewVirtualMachine(
	data [vmDataCount]byte,
	text [vmTextCount]byte,
) *VirtualMachine {
	vm := &VirtualMachine{
		R0: Register(0),
		R1: Register(0),
		SP: Register(0),
		PC: Register(0),
	}

	copy(vm.Memory[vmDataStart:vmDataEnd+1], data[:]) // end is exclusive

	for i := vmDataEnd + 1; i <= vmStackEnd; i++ {
		vm.Memory[i] = 0
	}

	vm.Memory[vmFlagStart] = 0

	copy(vm.Memory[vmTextStart:vmTextEnd+1], text[:]) // end is exclusive

	return vm
}

func (vm *VirtualMachine) String() string {
	out := ""
	out += fmt.Sprintf(
		"R0: %d, R1: %d, SP: %d, PC: %d\n",
		vm.R0, vm.R1, vm.SP, vm.PC,
	)
	for idx, byt := range vm.Memory {
		out += fmt.Sprintf("%d: %08b", idx, byt)
	}
	return out
}
