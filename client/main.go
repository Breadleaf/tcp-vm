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

	// TODO: make a ROUTER_PORT field in the docker file too

	routerIp := fmt.Sprintf("%s:11555", routerId)

	client, err := ofstp.NewClient(routerIp, 5 * time.Second)
	if err != nil {
		fmt.Printf("Failed to init client: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	pkt, err := ofstp.NewReturnPacket(0x00, []byte("Hello from client"))
	if err != nil {
		fmt.Printf("Failed to create return packet: %v\n", err)
		os.Exit(1)
	}

	resp, err := client.Do(pkt)
	if err != nil {
		fmt.Printf("Error while sending packet: %v\n", err)
		os.Exit(1)
	}

	retPacket, ok := resp.(*ofstp.ReturnPacket)
	if !ok {
		fmt.Printf("Expected return packet, got: %T\n", resp)
	}

	fmt.Printf("got response: %v\n", retPacket)
}
