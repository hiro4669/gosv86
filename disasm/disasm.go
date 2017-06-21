package disasm

import (
	"fmt"
	"os"
)

type Disasm struct {
	pc   uint16
	text []byte
}

func (dis *Disasm) fetch(opcode *OpCode) byte {
	v := dis.text[dis.pc]
	opcode.Add(v)
	dis.pc++
	return v
}

func (dis *Disasm) fetch2(opcode *OpCode) uint16 {
	var dh byte = dis.text[dis.pc+1]
	var dl byte = dis.text[dis.pc]
	opcode.Add(dl)
	opcode.Add(dh)
	var data uint16 = uint16(dh)<<8 | uint16(dl)
	//	var data uint16 = uint16(dis.text[dis.pc+1])<<8 | uint16(dis.text[dis.pc])
	dis.pc += 2
	return data
}

func (dis *Disasm) Init(text []byte) {
	fmt.Println("init called")
	dis.pc = 0
	dis.text = make([]byte, len(text))
	copy(dis.text, text)
}

func (dis *Disasm) setData(opcode *OpCode) {
	switch opcode.W {
	case 0:
		opcode.Data = uint16(dis.fetch(opcode))
	case 1:
		opcode.Data = dis.fetch2(opcode)
	default:
		fmt.Printf("invalid W = %d", opcode.W)
		os.Exit(1)
	}
}

func test(opcode *OpCode) {
	opcode.W = 1
	opcode.Reg = 3
}

func (dis *Disasm) Run() {
	var opcode OpCode
	var op byte
	for {
		if int(dis.pc) == len(dis.text) {
			break
		}

		switch op = dis.fetch(&opcode); op {
		case 0xb0, 0xb1, 0xb2, 0xb3, 0xb4, 0xb5, 0xb6, 0xb7,
			0xb8, 0xb9, 0xba, 0xbb, 0xbc, 0xbd, 0xbe, 0xbf:
			{
				opcode.setW((op >> 3) & 1)
				opcode.setReg(op & 7)
				dis.setData(&opcode)
				opcode.ShowOpCode()
			}
		default:
			{
				fmt.Printf("%02x is not implemented yet\n", op)
				os.Exit(1)
			}
		}

		opcode.Reset()
	}
	fmt.Printf("\npc = %d\n", dis.pc)
}
