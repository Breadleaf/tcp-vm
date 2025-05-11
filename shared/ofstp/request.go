package ofstp

import (
	"net"
)

// Request wraps an incoming connection + parsed Packet.
type Request struct {
	conn   net.Conn
	Packet Packet
}

// Respond writes p.Marshal() back on the original conn.
func (r *Request) Respond(p Packet) error {
	data, err := p.Marshal()
	if err != nil {
		return err
	}
	_, err = r.conn.Write(data)
	return err
}

// Conn exposes the underlying net.Conn for custom writes.
func (r *Request) Conn() net.Conn {
	return r.conn
}
