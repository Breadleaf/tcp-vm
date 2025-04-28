package router

import (
	"tcp-vm/shared/ofstp"
)

func main() {
	server := ofstp.NewServer()

	server.Register(ofstp.Return, func(req *ofstp.Request) {
	})
}
