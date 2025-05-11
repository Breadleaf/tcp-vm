package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"tcp-vm/shared/assembler"
	o "tcp-vm/shared/ofstp"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: client <program.asm>")
		os.Exit(1)
	}
	asm := os.Args[1]

	data, text, err := assembler.Assemble(asm)
	if err != nil {
		log.Fatal(err)
	}

	routerID := os.Getenv("ROUTER_ID")
	if routerID == "" {
		log.Fatal("ROUTER_ID not set")
	}
	addr := routerID + ":11555"

	cli, err := o.NewClient(addr, 5*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	// Register and immediately receive AskBusy
	pkt, err := cli.Do(&o.ReturnPacket{ExitCode: o.RegisterClientCode})
	if err != nil {
		log.Fatal(err)
	}
	rp := pkt.(*o.ReturnPacket)
	if rp.ExitCode != o.AskBusyCode {
		log.Fatalf("expected AskBusy; got %d", rp.ExitCode)
	}

	// Tell the router we are not busy
	_, err = cli.Do(&o.ReturnPacket{ExitCode: o.NotBusyCode})
	if err != nil {
		log.Fatal(err)
	}

	// Receive AskStateless
	pkt, err = cli.Do(&o.ReturnPacket{})
	if err != nil {
		log.Fatal(err)
	}
	rp = pkt.(*o.ReturnPacket)
	if rp.ExitCode != o.AskStatelessCode {
		log.Fatalf("expected AskStateless; got %d", rp.ExitCode)
	}

	// Send the real Stateless packet
	stateless, _ := o.NewStatelessPacket(data[:], text[:])
	_, err = cli.Do(stateless)
	if err != nil {
		log.Fatal(err)
	}

	// Finally, read back exit/output
	for {
		pkt, err := cli.Do(&o.ReturnPacket{ExitCode: 0})
		if err != nil {
			log.Fatal(err)
		}
		if rp, ok := pkt.(*o.ReturnPacket); ok {
			if len(rp.Output) > 0 {
				fmt.Printf("OUTPUT: %s\n", string(rp.Output))
				return
			}
			fmt.Printf("Exit code: %d\n", rp.ExitCode)
			return
		}
	}
}
