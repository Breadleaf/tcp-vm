package globals

import (
	"fmt"
	"os"
)

var DEBUG bool

func init() {
	DEBUG = os.Getenv("DEBUG") != ""
	if DEBUG {
		fmt.Println("DEBUG MODE IS ENABLED")
	}
}
