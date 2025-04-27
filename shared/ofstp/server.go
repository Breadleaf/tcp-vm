package ofstp

import (
	"io"
	"log"
	"net"
	"sync"
)

const (
	maxPacketSize = 1500
)

type HandlerFunc func(req *Request)

type Server struct {
	mutex    sync.RWMutex
	handlers map[PacketType]HandlerFunc
}

func NewServer() *Server {
	return &Server{
		handlers: make(map[PacketType]HandlerFunc),
	}
}

func (s *Server) Register(pt PacketType, h HandlerFunc) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.handlers[pt] = h
}

func (s *Server) Listen(addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, maxPacketSize)
	for {
		// TODO: read first line of header to find type then read into
		//       buf of correct size.
		// for now, always read into max size buffer
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Printf("")
			}
			return
		}

		pkt, err := ParsePacket(buf[:n])
		if err != nil {
			// TODO: send ReturnPacket with error
			continue
		}

		s.mutex.RLock()
		h, exists := s.handlers[pkt.Type()]
		s.mutex.RUnlock()
		if !exists {
			// no handler registered, ignore or send back default
			continue
		}

		req := &Request{
			conn:   conn,
			Packet: pkt,
		}
		go h(req)
	}
}
