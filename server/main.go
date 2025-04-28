package main

import (
	"fmt"
	"os"
	"tcp-vm/shared/assembler"
	g "tcp-vm/shared/globals"
	"tcp-vm/shared/vm"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run server/main.go [fpath]")
		os.Exit(1)
	}

	fmt.Println(os.Args[1])

	data, text, err := assembler.Assemble(os.Args[1])
	if err != nil {
		fmt.Printf("assembler does like that one: %v", err)
		os.Exit(1)
	}

	if len(data) != g.DataSectionLength || len(text) != g.TextSectionLength {
		fmt.Printf("len(data)=%d\nlen(text)=%d\n", len(data), len(text))
		fmt.Printf("assembler needs to check this and return correct type, this is to prevent future foot guns: %v", err)
		os.Exit(1)
	}

	vm := new(vm.VirtualMachine)
	vm.ResetFromStateless(data, text)

	fmt.Printf("vm:\n%v", vm)
}

/*
import (
	"fmt"
	"tcp-vm/shared/util"
	"tcp-vm/shared/ofstp"
)

func main() {
	logTag := "tcp-vm/server - main.go - main()"
	util.LogStart(logTag)
	defer util.LogEnd(logTag)

	server := ofstp.NewServer()

	server.Register(ofstp.Return, func(req *ofstp.Request) {
		b, err := req.Packet.Marshal()
		if err != nil {
			req.Respond(&ofstp.ReturnPacket{
				ExitCode: 0x01,
				Output: fmt.Appendf(
					[]byte{},
					"error: bad packet:\n%v\nmessage:\n%v\n",
					req.Packet,
					err,
				),
			})
		}
		fmt.Printf("server got:\n%+v\n", b)
	})
}
*/
