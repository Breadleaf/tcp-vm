package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	g "tcp-vm/shared/globals"
	o "tcp-vm/shared/ofstp"
	vm "tcp-vm/shared/vm"
)

func recvPacket(r io.Reader) (o.Packet, error) {
	header := make([]byte, 1)
	if _, err := io.ReadFull(r, header); err != nil {
		return nil, err
	}
	pt := o.PacketType(header[0])
	switch pt {
	case o.Stateless:
		rest := make([]byte, g.DataSectionLength+g.TextSectionLength)
		if _, err := io.ReadFull(r, rest); err != nil {
			return nil, err
		}
		return o.ParsePacket(append(header, rest...))
	case o.Return:
		exit := make([]byte, 1)
		if _, err := io.ReadFull(r, exit); err != nil {
			return nil, err
		}
		payload := make([]byte, 1498)
		n, _ := r.Read(payload)
		return o.ParsePacket(append(append(header, exit...), payload[:n]...))
	case o.Stateful:
		size := 4 + g.DataSectionLength + g.StackLength + g.FlagLength + g.TextSectionLength
		rest := make([]byte, size)
		if _, err := io.ReadFull(r, rest); err != nil {
			return nil, err
		}
		return o.ParsePacket(append(header, rest...))
	default:
		return nil, fmt.Errorf("unknown packet type: %v", pt)
	}
}

func main() {
	routerID := os.Getenv("ROUTER_ID")
	if routerID == "" {
		log.Fatal("ROUTER_ID not set")
	}
	conn, err := net.Dial("tcp", routerID+":11555")
	if err != nil {
		log.Fatalf("dial router: %v", err)
	}
	defer conn.Close()

	// register as a server
	regPkt, _ := o.NewReturnPacket(o.RegisterServerCode, nil)
	conn.Write(o.MustMarshal(regPkt))

	reader := bufio.NewReader(conn)
	for {
		pkt, err := recvPacket(reader)
		if err != nil {
			log.Fatalf("recv: %v", err)
		}
		switch p := pkt.(type) {
		case *o.StatelessPacket:
			var dataArr [g.DataSectionLength]byte
			var textArr [g.TextSectionLength]byte
			copy(dataArr[:], p.Data[:])
			copy(textArr[:], p.Text[:])

			machine := new(vm.VirtualMachine)
			machine.ResetFromStateless(dataArr, textArr)
			if err := machine.RunUntilStop(); err != nil {
				errPkt, _ := o.NewReturnPacket(1, []byte(err.Error()))
				conn.Write(o.MustMarshal(errPkt))
				continue
			}

			flag := machine.Memory[vm.FlagStart]
			if flag&g.HaltFlag != 0 {
				exitPkt, _ := o.NewReturnPacket(byte(machine.R0), nil)
				conn.Write(o.MustMarshal(exitPkt))
				outPkt, _ := o.NewReturnPacket(0, []byte(machine.Output))
				conn.Write(o.MustMarshal(outPkt))

			} else if flag&g.SleepFlag != 0 {
				wake := time.Now().Unix() + int64(machine.R0)
				buf := make([]byte, 8)
				binary.BigEndian.PutUint64(buf, uint64(wake))
				sleepPkt, _ := o.NewReturnPacket(0, buf)
				conn.Write(o.MustMarshal(sleepPkt))

				statePkt, _ := o.NewStatefulPacket(
					byte(machine.R0), byte(machine.R1),
					byte(machine.SP), byte(machine.PC),
					machine.Memory[vm.DataStart:vm.DataStart+g.DataSectionLength],
					machine.Memory[vm.StackStart:vm.StackStart+g.StackLength],
					machine.Memory[vm.FlagStart],
					machine.Memory[vm.TextStart:vm.TextStart+g.TextSectionLength],
				)
				conn.Write(o.MustMarshal(statePkt))
			}

		case *o.ReturnPacket:
			// ignore AskOutput
		}
	}
}
