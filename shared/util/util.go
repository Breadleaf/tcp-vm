package util

import (
	"fmt"
	"strings"
)

func GenerateLine(count int) string {
	return strings.Repeat("-", count)
}

func LogStart(logTag string) {
	fmt.Println(GenerateLine(100))
	fmt.Printf("%s - debug log start\n", logTag)
}

func LogMessage(log func()) {
	fmt.Println(GenerateLine(100))
	log()
}

func LogEnd(logTag string) {
	fmt.Println(GenerateLine(100))
	fmt.Printf("%s - debug log end\n", logTag)
}
