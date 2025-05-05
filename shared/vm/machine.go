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
	vm.SP = Register(vmStackStart)
	vm.PC = Register(vmTextStart)

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

func (vm *VirtualMachine) RunUntilStop() error {
	util.LogStart(FILE_LOG_TAG)
	defer util.LogEnd(FILE_LOG_TAG)

	const MaxStepsPerRun = 0xFFFF
	const BottomTwoMask = 0x03 // 0b 0000 0011
	const BottomThreeMask = 0x07 // 0b 0000 0111

	const (
		InstXa = 0x00
		InstXb = 0x01
		InstY = 0x02
		InstZ = 0x03
	)

	const (
		MOV = 0x00
		CMP = 0x01
		SHL = 0x02
		SHR = 0x03

		ADD = 0x00
		SUB = 0x01
		AND = 0x02
		ORR = 0x03

		NOT = 0x00
		PSH = 0x01
		POP = 0x02
		SYS = 0x03

		JMP = 0x00
		LDI = 0x01
		LDA = 0x02
		STA = 0x03
	)

	const (
		LT = 0b100
		EQ = 0b010
		GT = 0b001
	)

	binaryToRegister := map[byte]*Register{
		0x00: &vm.R0,
		0x01: &vm.R1,
		0x02: &vm.SP,
		0x03: &vm.PC,
	}

	for stepCount := 0; stepCount <= MaxStepsPerRun; stepCount++ {
		current := vm.Memory[vm.PC]
		vm.PC++

		top2 := (current >> 6) & BottomTwoMask
		middleTop2 := (current >> 4) & BottomTwoMask

		middleBottom2 := (current >> 2) & BottomTwoMask
		bottom2 := current & BottomTwoMask

		instType := top2
		instSpecifier := middleTop2

		util.LogMessage(func() {
			fmt.Printf("current: %08b\n", current)
			fmt.Printf("chuncked: %02b, %02b, %02b, %02b\n", top2, middleTop2, middleBottom2, bottom2)
		})

		switch instType {
		case InstXa:
			ra := binaryToRegister[middleBottom2]
			rb := binaryToRegister[bottom2]

			switch instSpecifier {
			case MOV:
				*ra = *rb
			case CMP:
				var flag byte
				switch {
				case *ra < *rb:
					flag = LT
				case *ra == *rb:
					flag = EQ
				case *ra > *rb:
					flag = GT
				}

				vm.Memory[vmFlagStart] = flag
			case SHL:
				*ra = *ra << *rb
			case SHR:
				*ra = *ra >> *rb
			}

		case InstXb:
			ra := binaryToRegister[middleBottom2]
			rb := binaryToRegister[bottom2]

			switch instSpecifier {
			case ADD:
				*ra = *ra + *rb
			case SUB:
				*ra = *ra - *rb
			case AND:
				*ra = *ra & *rb
			case ORR:
				*ra = *ra | *rb
			}

		case InstY:
			ra := binaryToRegister[bottom2]

			switch instSpecifier {
			case NOT:
				*ra = ^(*ra)
			case PSH:
				if vm.SP < vmStackStart {
					return fmt.Errorf("segfault: stack underflow")
				}
				if vm.SP > vmStackEnd {
					return fmt.Errorf("segfault: stack overflow")
				}

				vm.Memory[vm.SP] = byte(*ra)
				vm.SP++ // grow stack down
			case POP:
				if vm.SP < vmStackStart {
					return fmt.Errorf("segfault: stack underflow")
				}
				if vm.SP > vmStackEnd {
					return fmt.Errorf("segfault: stack overflow")
				}

				*ra = Register(vm.Memory[vm.SP])
				vm.SP-- // shrink stack up
			case SYS:
				fmt.Printf("SYSTEM CALL: %d\n", *ra)
				fmt.Printf("step count: %d\n", stepCount)
				return nil
			}

		case InstZ:
			ra := binaryToRegister[middleBottom2]
			imm := vm.Memory[vm.PC]
			vm.PC++

			switch instSpecifier {
			case JMP:
				mask := current & BottomThreeMask
				if vm.Memory[vmFlagStart] & mask != 0 {
					vm.PC = Register(imm)
				}
			case LDI:
				*ra = Register(imm)
			case LDA:
				*ra = Register(vm.Memory[imm])
			case STA:
				vm.Memory[imm] = byte(*ra)
			}
		}
	}

	return fmt.Errorf("program exceeded max number of steps: %d", MaxStepsPerRun)
}
