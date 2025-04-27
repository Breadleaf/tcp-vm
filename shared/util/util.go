package util

import (
	"fmt"
	"os"
	"strings"
	"tcp-vm/shared/globals"
)

func GenerateLine(count int) string {
	return strings.Repeat("-", count)
}

func LogStart(logTag string) {
	if globals.DEBUG {
		fmt.Println(GenerateLine(100))
		fmt.Printf("%s - debug log start\n", logTag)
	}
}

func LogMessage(log func()) {
	if globals.DEBUG {
		fmt.Println(GenerateLine(100))
		log()
	}
}

func LogEnd(logTag string) {
	if globals.DEBUG {
		fmt.Println(GenerateLine(100))
		fmt.Printf("%s - debug log end\n", logTag)
	}
}

func Assert(condition bool, message string) {
	if condition {
		return
	}
	panic(message)
}
