package main

import (
	"encoding/binary"
	"flag"
	"log"
	"net"
	"sync"
	"time"

	g "tcp-vm/shared/globals"
	o "tcp-vm/shared/ofstp"
)

type session struct {
	client net.Conn
	server net.Conn
}

type programJob struct {
	wake  time.Time
	conn  net.Conn
	state *o.StatefulPacket
}

type Router struct {
	*o.Server
	mutex          sync.Mutex
	waitingClients []net.Conn
	waitingServers []net.Conn
	sessions       map[net.Conn]*session
	programQueue   []*programJob
}

func NewRouter() *Router {
	s := o.NewServer()
	r := &Router{
		Server:       s,
		sessions:     make(map[net.Conn]*session),
		programQueue: make([]*programJob, 0),
	}
	s.Register(o.Return, r.handleReturn)
	s.Register(o.Stateless, r.handleStateless)
	s.Register(o.Stateful, r.handleStateful)
	return r
}

func (r *Router) tryMatch() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	now := time.Now()
	// wake sleeping jobs
	i := 0
	for i < len(r.programQueue) {
		job := r.programQueue[i]
		if !job.wake.After(now) {
			r.waitingServers = append(r.waitingServers, job.conn)
			r.programQueue = append(r.programQueue[:i], r.programQueue[i+1:]...)
		} else {
			i++
		}
	}

	if len(r.waitingClients) > 0 && len(r.waitingServers) > 0 {
		cli := r.waitingClients[0]
		srv := r.waitingServers[0]
		r.waitingClients = r.waitingClients[1:]
		r.waitingServers = r.waitingServers[1:]

		r.sessions[cli] = &session{client: cli, server: srv}
		r.sessions[srv] = &session{client: cli, server: srv}

		askBusy, _ := o.NewReturnPacket(o.AskBusyCode, nil)
		cli.Write(o.MustMarshal(askBusy))
	}
}

func (r *Router) handleReturn(req *o.Request) {
	rp := req.Packet.(*o.ReturnPacket)
	conn := req.Conn()

	switch rp.ExitCode {
	case o.RegisterClientCode:
		r.mutex.Lock()
		r.waitingClients = append(r.waitingClients, conn)
		r.mutex.Unlock()
		r.tryMatch()

	case o.RegisterServerCode:
		r.mutex.Lock()
		r.waitingServers = append(r.waitingServers, conn)
		r.mutex.Unlock()
		r.tryMatch()

	case o.NotBusyCode:
		sess := r.sessions[conn]
		askStateless, _ := o.NewReturnPacket(o.AskStatelessCode, nil)
		sess.client.Write(o.MustMarshal(askStateless))

	default:
		sess := r.sessions[conn]
		sess.client.Write(o.MustMarshal(rp))
	}
}

func (r *Router) handleStateless(req *o.Request) {
	st := req.Packet.(*o.StatelessPacket)
	sess := r.sessions[req.Conn()]
	sess.server.Write(o.MustMarshal(st))
}

func (r *Router) handleStateful(req *o.Request) {
	sp := req.Packet.(*o.StatefulPacket)
	conn := req.Conn()
	sess := r.sessions[conn]

	flag := sp.Flag
	if flag&g.HaltFlag != 0 {
		askOut, _ := o.NewReturnPacket(o.AskOutputCode, nil)
		sess.server.Write(o.MustMarshal(askOut))
	} else if flag&g.SleepFlag != 0 {
		wakeUnix := int64(binary.BigEndian.Uint64(sp.Data[0:8]))
		r.mutex.Lock()
		r.programQueue = append(r.programQueue, &programJob{
			wake:  time.Unix(wakeUnix, 0),
			conn:  conn,
			state: sp,
		})
		r.mutex.Unlock()
		r.tryMatch()
	}
}

func main() {
	port := flag.String("port", "11555", "router listen port")
	flag.Parse()

	addr := ":" + *port
	log.Printf("router listening on %s\n", addr)
	r := NewRouter()
	if err := r.Listen(addr); err != nil {
		log.Fatal(err)
	}
}
