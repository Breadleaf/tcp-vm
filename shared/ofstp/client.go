package ofstp

import (
	"fmt"
	"io"
	"net"
	"time"
)

type Client struct {
	conn    net.Conn
	timeout time.Duration
}

func NewClient(addr string, timeout time.Duration) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Client{
		conn: conn,
		timeout: timeout,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Do(packet Packet) (Packet, error) {
	raw, err := packet.Marshal()
	if err != nil {
		return nil, err
	}

	// read with timeout
	if c.timeout > 0 {
		c.conn.SetWriteDeadline(time.Now().Add(c.timeout))
	}
	if _, err := c.conn.Write(raw); err != nil {
		return nil, err
	}

	// read response header
	header := make([]byte, 1)
	if c.timeout > 0 {
		c.conn.SetReadDeadline(time.Now().Add(c.timeout))
	}
	if _, err := io.ReadFull(c.conn, header); err != nil {
		return nil, err
	}

	pt := PacketType(header[0])

	var restLength int
	switch pt {
	case Stateless:
		restLength = 16 + 175
	case Stateful:
		restLength = (1*4) + 16 + 64 + 1 + 175
	case Return:
		// read the exit code
		exit := make([]byte, 1)
		if _, err := io.ReadFull(c.conn, exit); err != nil {
			return nil, err
		}

		// read rest of file
		payload := make([]byte, 1498)
		n, _ := c.conn.Read(payload) // get length of remainder
		raw = append(header, exit...)
		raw = append(raw, payload[:n]...) // read up to remainder
		return ParsePacket(raw)
	default:
		return nil, fmt.Errorf("unknown response type: %v", pt)
	}

	// read the rest of fixed length packet
	rest := make([]byte, restLength)
	if _, err := io.ReadFull(c.conn, rest); err != nil {
		return nil, err
	}
	raw = append(header, rest...)

	return ParsePacket(raw)
}
