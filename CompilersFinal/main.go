package main

import (
	"fmt"
	"os"
	"tcp-vm/shared/assembler"
	"tcp-vm/shared/vm"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("usage: %s [path to `.asm` file]\n", os.Args[0])
		os.Exit(1)
	}

	path := os.Args[1]

	data, text, err := assembler.Assemble(path)
	if err != nil {
		fmt.Printf("assembler error: %v\n", err)
		os.Exit(1)
	}

	v := new(vm.VirtualMachine)
	v.ResetFromStateless(data, text)

	err = v.RunUntilStop()
	if err != nil {
		fmt.Printf("error in vm: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("final machine state:\n%s", v)

	arg := v.Memory[v.SP - 1] // get item at the top of the stack
	sys := v.Memory[vm.FlagStart]
	fmt.Printf("sys: %d, arg: %d\n", sys, arg)
}
