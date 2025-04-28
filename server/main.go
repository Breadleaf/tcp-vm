package main

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
