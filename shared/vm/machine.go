package vm

import (
	"fmt"
	"tcp-vm/shared/util"
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

type Memory []uint8

// virtual machine

type VirtualMachine struct {
	R0     Register
	R1     Register
	SP     Register
	PC     Register
	Memory Memory
}

func NewVirtualMachine()
