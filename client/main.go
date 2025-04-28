package main

import (
	"fmt"
	"os"
	"tcp-vm/shared/ofstp"
	"tcp-vm/shared/util"
	"time"
)

func main() {
	logTag := "tcp-vm/client - main.go - main()"
	util.LogStart(logTag)
	defer util.LogEnd(logTag)

	routerId := os.Getenv("ROUTER_ID")
	if routerId == "" {
		fmt.Println("ROUTER_ID is not set, unrecoverable...")
		os.Exit(1)
	}

	client, err := ofstp.NewClient(routerId+":11555", 5*time.Second)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// register with router
	_, err = client.Do(&ofstp.ReturnPacket{
		ExitCode: ofstp.RegisterClientCode,
		Output:   nil,
	})
	if err != nil {
		panic(err)
	}

	// wait for "AskBusy"
	resp, err := client.Do(&ofstp.ReturnPacket{
		ExitCode: 0,
		Output:   nil,
	})
	if err != nil {
		panic(err)
	}
	rp := resp.(*ofstp.ReturnPacket)
	if rp.ExitCode != ofstp.AskBusyCode {
		panic("expected AskBusy")
	}

	// tell router we are free
	_, err = client.Do(&ofstp.ReturnPacket{
		ExitCode: ofstp.NotBusyCode,
		Output:   nil,
	})
	if err != nil {
		panic(err)
	}

	// wait for "AskStateless"
	resp, err = client.Do(&ofstp.ReturnPacket{
		ExitCode: 0,
		Output:   nil,
	})
	rp = resp.(*ofstp.ReturnPacket)
	if rp.ExitCode != ofstp.AskStatelessCode {
		panic("expected AskStateless")
	}

	// send real Stateless packet

}
