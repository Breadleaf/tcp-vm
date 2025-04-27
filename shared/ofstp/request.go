package ofstp

import (
	"net"
)

type Request struct {
	conn   net.Conn
	Packet Packet
}

func (r *Request) Respond(p Packet) error {
	data, err := p.Marshal()
	if err != nil {
		return err
	}

	_, err = r.conn.Write(data)
	return err
}
