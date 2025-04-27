package vm

import (
	_"fmt"
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
	vmDataEnd = vmDataStart + vmDataCount - 1
)

const (
	vmStackStart = vmDataEnd + 1
	vmStackCount = 64
	vmStackEnd = vmStackStart + vmStackCount - 1
)

const (
	vmFlagStart = vmStackEnd + 1
	vmFlagCount = 1
	vmFlagEnd = vmFlagStart + vmFlagCount - 1
)

const (
	vmTextStart = vmFlagEnd + 1
	vmTextCount = 175
	vmTextEnd = vmTextStart + vmTextCount - 1
)

const (
	vmMemStart = vmDataStart
	vmMemEnd = vmTextEnd
)

// program counter

const (
	vmPCMask = vmMemSizeBytes - 1 // 0xFF for 256 bytes
)

type ProgramCounter struct {
	position uint16
}

func NewProgramCounter() *ProgramCounter {
	return &ProgramCounter{
		position: uint16(vmTextStart),
	}
}

func (pc *ProgramCounter) SetPosition(position uint16) {
	util.Assert(
		position <= vmMemEnd,
		"ProgramCounter_SetPosition() - position not in bounds - high",
	)

	// this should be impossible since position is a uint, this it will wrap
	util.Assert(
		position >= vmMemStart,
		"ProgramCounter_SetPosition() - position not in bounds - low",
	)

	pc.position = position
}

func (pc *ProgramCounter) IncrementPositionBy(count uint16) {
	pc.position = (pc.position + count) & vmPCMask
}

func (pc *ProgramCounter) IncrementPosition(count uint16) {
	pc.IncrementPositionBy(1)
}

// memory

type Memory struct {
	data [vmMemSizeBytes]uint8 // 256 bytes
}

func NewMemory() *Memory {
	memory := &Memory{}
	return memory
}

func (m *Memory) WriteU8(addr, val uint8) error {

}
