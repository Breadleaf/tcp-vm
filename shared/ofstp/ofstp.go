package ofstp

import (
	"errors"
	"fmt"
)

const (
	dataSize  = 16
	stackSize = 64
	flagSize  = 1
	textSize  = 1
)

type PacketType byte

const (
	Stateless PacketType = 0x01
	Stateful  PacketType = 0x02
	Return    PacketType = 0x03
)

func (pt PacketType) String() string {
	switch pt {
	case Stateless:
		return "PacketType.Stateless"
	case Stateful:
		return "PacketType.Stateful"
	case Return:
		return "PacketType.Return"
	default:
		return fmt.Sprintf("PacketType.Unknown(0x%02X)", byte(pt))
	}
}

// common interface for all OFSTP packets

type Packet interface {
	Type() PacketType
	Marshal() ([]byte, error)
}

// stateless packet (1 + 16 + 175)

type StatelessPacket struct {
	Data [dataSize]byte
	Text [textSize]byte
}

func NewStatelessPacket(data, text []byte) (*StatelessPacket, error) {
	if len(data) != dataSize || len(text) != textSize {
		return nil, fmt.Errorf(
			"StatelessPacket got: len(data): %d, len(text): %d",
			len(data),
			len(text),
		)
	}

	var p StatelessPacket
	copy(p.Data[:], data)
	copy(p.Text[:], text)
	return &p, nil
}

func (p *StatelessPacket) Type() PacketType {
	return Stateless
}

func (p *StatelessPacket) Marshal() ([]byte, error) {
	buf := make([]byte, 1+16+175)
	buf[0] = byte(Stateless)
	copy(buf[1:17], p.Data[:])
	copy(buf[17:], p.Text[:])
	return buf, nil
}

// stateful packet (1 + 1*4 + 16 + 64 + 1 + 175)

type StatefulPacket struct {
	R0, R1, SP, PC byte
	Data           [dataSize]byte
	Stack          [stackSize]byte
	Flag           byte
	Text           [textSize]byte
}

func NewStatefulPacket(
	r0, r1, sp, pc byte,
	data []byte, stack []byte,
	flag byte, text []byte,
) (*StatefulPacket, error) {
	if len(data) != 16 || len(stack) != 64 || len(text) != 175 {
		return nil, fmt.Errorf(
			"StatefulPacket got: len(data): %d, len(stack): %d, len(text): %d",
			len(data), len(stack), len(text),
		)
	}

	var p StatefulPacket
	p.R0, p.R1, p.SP, p.PC = r0, r1, sp, pc
	copy(p.Data[:], data)
	copy(p.Stack[:], stack)
	p.Flag = flag
	copy(p.Text[:], text)
	return &p, nil
}

func (p *StatefulPacket) Type() PacketType {
	return Stateful
}

func (p *StatefulPacket) Marshal() ([]byte, error) {
	buf := make([]byte, 1+4+16+64+1+175)
	i := 0
	buf[i] = byte(Stateful)
	i++
	buf[i], buf[i+1], buf[i+2], buf[i+3] = p.R0, p.R1, p.SP, p.PC
	i += 4
	copy(buf[i:i+16], p.Data[:])
	i += 16
	copy(buf[i:i+64], p.Stack[:])
	i += 64
	buf[i] = p.Flag
	i++
	copy(buf[i:], p.Text[:])
	return buf, nil
}

// return packet (1 + 1 + up to 1498)

type ReturnPacket struct {
	ExitCode byte
	Output   []byte
}

func NewReturnPacket(exit byte, out []byte) (*ReturnPacket, error) {
	if len(out) > 1498 {
		return nil, fmt.Errorf(
			"ReturnPacket: payload too large: len(out): %d",
			len(out),
		)
	}

	return &ReturnPacket{
		ExitCode: exit,
		Output:   out,
	}, nil
}

func (p *ReturnPacket) Type() PacketType {
	return Return
}

func (p *ReturnPacket) Marshal() ([]byte, error) {
	buf := make([]byte, 2+len(p.Output))
	buf[0] = byte(Return)
	buf[1] = p.ExitCode
	copy(buf[2:], p.Output)
	return buf, nil
}

// unified "constructor"

func ParsePacket(raw []byte) (Packet, error) {
	if len(raw) == 0 {
		return nil, errors.New("empty packet")
	}

	switch PacketType(raw[0]) {
	case Stateless:
		if len(raw) != 192 {
			return nil, fmt.Errorf(
				"invalid Stateless length: %d",
				len(raw),
			)
		}
		data := raw[1:17]
		text := raw[17:]
		return NewStatelessPacket(data, text)
	case Stateful:
		expect := 1 + 4 + 16 + 64 + 1 + 175
		if len(raw) != expect {
			return nil, fmt.Errorf(
				"invalid Stateful length: %d != %d",
				len(raw),
				expect,
			)
		}
		r0, r1, sp, pc := raw[1], raw[2], raw[3], raw[4]
		data := raw[5:21]
		stack := raw[21:85]
		flag := raw[85]
		text := raw[86:]
		return NewStatefulPacket(
			r0, r1, sp, pc, data, stack, flag, text,
		)
	case Return:
		if len(raw) < 2 {
			return nil, fmt.Errorf(
				"invalid Return length: %d",
				len(raw),
			)
		}
		exit := raw[1]
		out := raw[2:]
		return NewReturnPacket(exit, out)
	default:
		return nil, fmt.Errorf("unknown packet type 0x%02X", raw[0])
	}
}
