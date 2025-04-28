package ofstp

import (
	"sync"
	g "tcp-vm/shared/globals"
	o "tcp-vm/shared/ofstp"
)

type session struct {
	client, server *o.Request
}

type Router struct {
	*o.Server
	mutex          sync.Mutex
	waitingClients []*o.Request // TODO: use util.Stack or util.Queue
	waitingServers []*o.Request
	sessions       map[*o.Request]*session
}

func NewRouter() *Router {
	s := o.NewServer()
	r := &Router{
		Server:   s,
		sessions: make(map[*o.Request]*session),
	}

	s.Register(o.Return, r.handleReturn)
	s.Register(o.Stateless, r.handleStateless)
	s.Register(o.Stateful, r.handleStateful)

	return r
}

// attempt to pair a waiting client to a waiting server
func (r *Router) tryMatch() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if len(r.waitingClients) == 0 || len(r.waitingServers) == 0 {
		return
	}

	// create the session

	cli := r.waitingClients[0]
	r.waitingClients = r.waitingClients[1:]

	srv := r.waitingServers[0]
	r.waitingServers = r.waitingServers[1:]

	sess := &session{
		client: cli,
		server: srv,
	}
	r.sessions[cli] = sess
	r.sessions[srv] = sess

	// ask the client if it is not busy
	askBusy, _ := o.NewReturnPacket(o.AskBusyCode, nil)
	cli.Respond(askBusy)
}

// handle registrations, "not busy", and final output forwarding
func (r *Router) handleReturn(req *o.Request) {
	rp := req.Packet.(*o.ReturnPacket)

	switch rp.ExitCode {
	case o.RegisterClientCode:
		r.mutex.Lock()
		r.waitingClients = append(r.waitingClients, req)
		r.mutex.Unlock()
		r.tryMatch()

	case o.RegisterServerCode:
		r.mutex.Lock()
		r.waitingServers = append(r.waitingServers, req)
		r.mutex.Unlock()
		r.tryMatch()

	case o.NotBusyCode:
		sess := r.sessions[req]
		askStateless, _ := o.NewReturnPacket(o.AskStatelessCode, nil)
		sess.client.Respond(askStateless)

	default:
		// any other ReturnPacket should be forwarded to the client who
		// started this session
		sess := r.sessions[req]
		sess.client.Respond(rp)
	}
}

// handle forwarding client -> router Stateless packets over to the matched server
func (r *Router) handleStateless(req *o.Request) {
	st := req.Packet.(*o.StatelessPacket)
	sess := r.sessions[req]
	sess.server.Respond(st)
}

// handle when server sends Stateful back, inspect flag
func (r *Router) handleStateful(req *o.Request) {
	sp := req.Packet.(*o.StatefulPacket)
	sess := r.sessions[req]

	if sp.Flag&g.HaltFlag != 0 {
		askOut, _ := o.NewReturnPacket(o.AskOutputCode, nil)
		sess.server.Respond(askOut)
	} else {
		// TODO: implement queue for sleeping
	}
}
