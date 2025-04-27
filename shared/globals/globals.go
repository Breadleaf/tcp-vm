package globals

import (
	"os"
	"fmt"
)

var DEBUG bool

func init() {
	DEBUG = os.Getenv("DEBUG") != ""
	if DEBUG {
		fmt.Println("DEBUG MODE IS ENABLED")
	}
}
