package day16

import (
	"fmt"
	"math/bits"
	"encoding/hex"
	"io/ioutil"
	"os"
	"math"
)

func byteMask(offset, len int) byte {
	var b byte
	for i := 0; i < len; i++ {
		b = (b << 1) + 1
	}
	return b << (8 - len - offset)
}

func extractBits(b byte, offset, len int) int {
	mask := byteMask(offset, len)
	return int((b & mask) >> (8 - len - offset))
}

type PacketType int

const (
	Sum PacketType = iota
	Product
	Minimum
	Maximum
	Literal
	GreaterThan
	LessThan
	Equal
)

type LengthTypeId int

const (
	TotalLength LengthTypeId = 0
	SubPacketsCount = 1
)

type BasePacket struct {
	version int
	packetType PacketType
}

type LiteralPacket struct {
	BasePacket
	literal int
}

type OperatorPacket struct {
	BasePacket
	lengthTypeId LengthTypeId
	subPackets []Packet
}

type Packet interface {
	Version() int
	Type() PacketType
	SumVersions() int
	Eval() int
}

func (lp LiteralPacket) Version() int {
	return lp.version
}

func (lp LiteralPacket) Type() PacketType {
	return lp.packetType
}

func (lp LiteralPacket) SumVersions() int {
	return lp.Version()
}

func (lp LiteralPacket) Eval() int {
	return lp.literal
}

func (bp BasePacket) Version() int {
	return bp.version
}

func (bp BasePacket) Type() PacketType {
	return bp.packetType
}

func (bp BasePacket) SumVersions() int {
	return bp.Version()
}

func (bp BasePacket) Eval() int {
	return 0
}

func (op OperatorPacket) Version() int {
	return op.version
}

func (op OperatorPacket) Type() PacketType {
	return op.packetType
}

func (op OperatorPacket) SumVersions() int {
	sum := op.Version()
	for _, p := range op.subPackets {
		sum += p.SumVersions()
	}
	return sum
}

func (op OperatorPacket) Eval() int {
	switch op.packetType {
	case Sum:
		sum := 0
		for _, p := range op.subPackets {
			sum += p.Eval()
		}
		return sum
	case Product:
		product := 1
		for _, p := range op.subPackets {
			product *= p.Eval()
		}
		return product
	case Minimum:
		min := math.MaxInt
		for _, p := range op.subPackets {
			val := p.Eval()
			if val < min {
				min = val
			}
		}
		return min
	case Maximum:
		max := math.MinInt
		for _, p := range op.subPackets {
			val := p.Eval()
			if val > max {
				max = val
			}
		}
		return max
	case GreaterThan:
		if op.subPackets[0].Eval() > op.subPackets[1].Eval() {
			return 1
		} else {
			return 0
		}
	case LessThan:
		if op.subPackets[0].Eval() < op.subPackets[1].Eval() {
			return 1
		} else {
			return 0
		}
	case Equal:
		if op.subPackets[0].Eval() == op.subPackets[1].Eval() {
			return 1
		} else {
			return 0
		}
	}
	return 0
}

type bytes []byte

func (bs bytes) subBits(offset, len int) int {
	var sub int
	extractEnd := offset + len
	for i, b := range bs {
		byteEnd := (i + 1) * 8
		prevByteEnd := i * 8
		currentByteOffset := offset - prevByteEnd
		// completely contained within current byte
		if offset >= prevByteEnd && extractEnd <= byteEnd {
			return extractBits(b, currentByteOffset, len)
		} else if offset < byteEnd {
			var current int
			var extractLen int
			if offset < prevByteEnd && extractEnd > byteEnd {
				extractLen = 8
				current = int(b)
			} else if offset < prevByteEnd && extractEnd <= byteEnd {
				extractLen = extractEnd - prevByteEnd
				current = extractBits(b, 0, extractLen)
			} else if offset > prevByteEnd {
				extractLen = byteEnd - offset
				current = extractBits(b, offset - prevByteEnd, extractLen)
			}
			// combine with bits from previous bytes
			sub = (sub << extractLen) + current
			if extractEnd <= byteEnd {
				return sub
			}
		}
	}
	return 0
}

func checkBit(i int, mask int) int {
	return bits.OnesCount64(uint64(i & (1 << mask)))
}

type ParseState struct {
	decoded bytes
	offset int
}

func makeParseState(input string) *ParseState {
	bs, _ := hex.DecodeString(input)
	decoded := bytes(bs)
	return &ParseState{ decoded: decoded, offset: 0 }
}

func (state *ParseState) nextBits(len int) int {
	bits := state.decoded.subBits(state.offset, len)
	state.offset += len
	// msgFmt := fmt.Sprintf("parsed %d bits: %%0%db new offset is %d\n", len, len, state.offset)
	// fmt.Printf(msgFmt, bits)
	return bits
}

func (state *ParseState) parseNextPacket() Packet {
	// fmt.Printf("parsing packet starting at offset %d\n", state.offset)
	version := state.nextBits(3)
	packetType := PacketType(state.nextBits(3))
	base := BasePacket{ version: version, packetType: packetType}

	if packetType == Literal {
		// fmt.Println("parsing literal packet")
		val := 0
		groupLen := 5
		lastBit := len(state.decoded) * 8
		for state.offset < lastBit {
			nextBits := state.nextBits(groupLen)
			// fmt.Printf("group %b starting at offset %d\n", nextBits, state.offset)
			val = (val << 4) + (0b1111 & nextBits)
			if checkBit(nextBits, 4) == 0 {
				break
			}
		}
		return LiteralPacket{BasePacket: base, literal: val }
	} else {
		// fmt.Println("parsing operator packet")
		lengthTypeId := LengthTypeId(state.nextBits(1))
		switch lengthTypeId {
		case TotalLength:
			packetsLength := state.nextBits(15)
			// fmt.Printf("packetsLength %d\n", packetsLength)
			packetsEnd := state.offset + packetsLength
			packets := make([]Packet, 0)
			for state.offset < packetsEnd {
				packets = append(packets, state.parseNextPacket())
			}
			return OperatorPacket{BasePacket: base, lengthTypeId: lengthTypeId, subPackets: packets }
		case SubPacketsCount:
			packetsCount := state.nextBits(11)
			// fmt.Printf("packetsCount %d\n", packetsCount)
			packets := make([]Packet, 0)
			for i := 0; i < packetsCount; i++ {
				packets = append(packets, state.parseNextPacket())
			}
			return OperatorPacket{BasePacket: base, lengthTypeId: lengthTypeId, subPackets: packets }
		}
	}
	
	return base
}

func parsePacket(input string) Packet {
	return makeParseState(input).parseNextPacket()
}

func testEval(input string, expected int) {
	val := parsePacket(input).Eval()
	if val == expected {
		fmt.Printf("SUCCESS! value of packet %s %d == %d\n", input, val, expected)
	} else {
		fmt.Printf("FAILURE! value of packet %s %d != %d\n", input, val, expected)
	}
}

func runTests() {
	testEval("C200B40A82", 3)
	testEval("04005AC33890", 54)
	testEval("880086C3E88112", 7)
	testEval("CE00C43D881120", 9)
	testEval("D8005AC2A8F0", 1)
	testEval("F600BC2D8F", 0)
	testEval("9C005AC2F8F0", 0)
	testEval("9C0141080250320F1802104A08", 1)
}

func Day16() {
	// inputs := []string{ "38006F45291200", "EE00D40C823060" }
	// for _, input := range inputs {
	// 	fmt.Printf("parsing input %s\n", input)
	// 	packet := parsePacket(input).(OperatorPacket)
	// 	fmt.Printf("parsed packet %v\n", packet)
	// 	fmt.Printf("sum of versions = %d\n", packet.SumVersions())
	// }
	input, _ := ioutil.ReadFile(os.Args[1])

	fmt.Printf("packet value = %d\n", parsePacket(string(input)).Eval())
	// fmt.Printf("parsed packet %v\n", packet)
	// fmt.Printf("sum of versions = %d\n", packet.SumVersions())
	// runTests()
}
