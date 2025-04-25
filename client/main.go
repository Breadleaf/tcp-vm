package main

import (
	"fmt"
	"tcp-vm/shared/util"
)

func main() {
	logTag := "tcp-vm/client - main.go - main()"
	util.LogStart(logTag)
	defer util.LogEnd(logTag)

	util.LogMessage(func() {
		fmt.Println("Hello from client")
	})
}
