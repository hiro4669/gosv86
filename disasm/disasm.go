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

func (dis *Disasm) setMrr(opcode *OpCode) {
	v := dis.fetch(opcode)
	opcode.setMod((v >> 6) & 3)
	opcode.setReg((v >> 3) & 7)
	opcode.setRm(v & 7)
	dis.resolveDisp(opcode)
}

func (dis *Disasm) resolveDisp(opcode *OpCode) {
	switch opcode.Mod {
	case 0:
		{
			if opcode.Rm == 6 {
				opcode.setDisp(int16(dis.fetch2(opcode)))
			}
		}
	case 1:
		opcode.setDisp(int16(int8(dis.fetch(opcode))))
	case 2:
		opcode.setDisp(int16(dis.fetch2(opcode)))
	case 3:
	default:
		{
			fmt.Printf("invalid Mod = %d\n", opcode.Mod)
			os.Exit(1)
		}
	}
}

func (dis *Disasm) disaRMftR(op byte, opcode *OpCode, opName string, pc uint16) {
	opcode.setW(op & 1)
	opcode.setD((op >> 1) & 1)
	dis.setMrr(opcode)
	dumpRMftR(opcode, pc, opName)
}

func (dis *Disasm) Run() {
	var opcode OpCode
	var op byte
	for {
		if int(dis.pc) == len(dis.text) {
			break
		}
		prevPc := dis.pc
		switch op = dis.fetch(&opcode); op {
		case 0x00, 0x01, 0x02, 0x03:
			{
				dis.disaRMftR(op, &opcode, "add", prevPc)
			}
		case 0x30, 0x31, 0x32, 0x33:
			{
				dis.disaRMftR(op, &opcode, "xor", prevPc)
			}
		case 0x88, 0x89, 0x8a, 0x8b:
			{
				dis.disaRMftR(op, &opcode, "mov", prevPc)
			}
		case 0x8d:
			{
				dis.disaRMftR(op, &opcode, "lea", prevPc)
			}
		case 0xb0, 0xb1, 0xb2, 0xb3, 0xb4, 0xb5, 0xb6, 0xb7,
			0xb8, 0xb9, 0xba, 0xbb, 0xbc, 0xbd, 0xbe, 0xbf:
			{

				opcode.setW((op >> 3) & 1)
				opcode.setReg(op & 7)
				dis.setData(&opcode)
				dumpMov(&opcode, prevPc)
			}
		case 0xcd:
			{

				dis.setData(&opcode)
				dumpInt(&opcode, prevPc)

			}
		default:
			{
				fmt.Printf("%02x is not implemented yet\n", op)
				os.Exit(1)
			}
		}

		opcode.Reset()
	}
}
