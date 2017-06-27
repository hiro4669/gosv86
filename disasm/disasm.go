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

func (dis *Disasm) lookahead() byte {
	return dis.text[dis.pc]
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

func (dis *Disasm) setSData(opcode *OpCode) {
	switch {
	case opcode.S == 0 && opcode.W == 1:
		{
			opcode.Data = dis.fetch2(opcode)
		}
	default:
		{
			opcode.Data = uint16(dis.fetch(opcode))
		}
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

func (dis *Disasm) disaIfRM(op byte, opcode *OpCode, opName string, pc uint16) {
	opcode.setW(op & 1)
	opcode.setS((op >> 1) & 1)
	dis.setMrr(opcode)
	dis.setSData(opcode)
	dumpIfRM(opcode, pc, opName)
}

func (dis *Disasm) disaItRM(op byte, opcode *OpCode, opName string, pc uint16) {
	opcode.setW(op & 1)
	dis.setMrr(opcode)
	dis.setData(opcode)
	dumpItRM(opcode, pc, opName)
}

func (dis *Disasm) disaLogic(op byte, opcode *OpCode, opName string, pc uint16) {
	opcode.setW(op & 1)
	opcode.setV((op >> 1) & 1)
	dis.setMrr(opcode)
	dumpLogic(opcode, pc, opName)
}

func (dis *Disasm) disaJump(op byte, opcode *OpCode, opName string, prevPc uint16) {
	off := dis.fetch(opcode)
	opcode.setJDisp(uint16((int32(dis.pc) + int32(int8(off))) & 0xffff))
	dumpJump(opcode, prevPc, opName)
}

func (dis *Disasm) disa2Jump(op byte, opcode *OpCode, opName string, prevPc uint16) {
	off := dis.fetch2(opcode)
	opcode.setJDisp(uint16((int32(dis.pc) + int32(int16(off))) & 0xffff))
	dumpJump(opcode, prevPc, opName)
}

func (dis *Disasm) disaOneMrr(op byte, opcode *OpCode, opName string, prevPc uint16) {
	opcode.setW(1)
	dis.setMrr(opcode)
	dumpOneMrr(opcode, prevPc, opName)
}

func (dis *Disasm) disaOneMrrW(op byte, opcode *OpCode, opName string, prevPc uint16) {
	opcode.setW(op & 1)
	dis.setMrr(opcode)
	dumpOneMrr(opcode, prevPc, opName)
}

func (dis *Disasm) disaOneReg(op byte, opcode *OpCode, opName string, prevPc uint16) {
	opcode.setReg(op & 7)
	dumpOneReg(opcode, prevPc, opName)
}

func (dis *Disasm) disaInOutPort(op byte, opcode *OpCode, opName string, prevPc uint16) {
	opcode.setW(op & 1)
	opcode.setPort(dis.fetch(opcode))
	dumpInOutPort(opcode, prevPc, opName)
}

func (dis *Disasm) disaInOutVar(op byte, opcode *OpCode, opName string, prevPc uint16) {
	opcode.setW(op & 1)
	dumpInOutVar(opcode, prevPc, opName)
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
				if int(dis.pc) == len(dis.text) {
					dumpUndefined(&opcode, prevPc)
					break
				}
				dis.disaRMftR(op, &opcode, "add", prevPc)
			}
		case 0x08, 0x09, 0x0a, 0x0b:
			{
				dis.disaRMftR(op, &opcode, "or", prevPc)
			}
		case 0x18, 0x19, 0x1a, 0x1b:
			{
				dis.disaRMftR(op, &opcode, "sbb", prevPc)
			}
		case 0x20, 0x21, 0x22, 0x23:
			{
				dis.disaRMftR(op, &opcode, "and", prevPc)
			}
		case 0x30, 0x31, 0x32, 0x33:
			{
				dis.disaRMftR(op, &opcode, "xor", prevPc)
			}
		case 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4e, 0x4f:
			{
				dis.disaOneReg(op, &opcode, "dec", prevPc)
			}
		case 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57:
			{
				dis.disaOneReg(op, &opcode, "push", prevPc)
			}
		case 0x58, 0x59, 0x5a, 0x5b, 0x5c, 0x5d, 0x5e, 0x5f:
			{
				dis.disaOneReg(op, &opcode, "pop", prevPc)
			}
		case 0x73:
			{
				dis.disaJump(op, &opcode, "jnb", prevPc)
			}
		case 0x74:
			{
				dis.disaJump(op, &opcode, "je", prevPc)
			}
		case 0x75:
			{
				dis.disaJump(op, &opcode, "jne", prevPc)
			}
		case 0x7c:
			{
				dis.disaJump(op, &opcode, "jl", prevPc)
			}
		case 0x7d:
			{
				dis.disaJump(op, &opcode, "jnl", prevPc)
			}
		case 0x80, 0x81, 0x82, 0x83:
			{
				nv := dis.lookahead()
				switch (nv >> 3) & 7 {
				case 0:
					{ // add
						dis.disaIfRM(op, &opcode, "add", prevPc)
					}
				case 5:
					{ // sub
						dis.disaIfRM(op, &opcode, "sub", prevPc)

					}
				case 7:
					{ // cmp
						dis.disaIfRM(op, &opcode, "cmp", prevPc)
					}
				default:
					{
						fmt.Println("not implemented for next byte in 0x80~0x83")
						os.Exit(1)
					}
				}
			}
		case 0x88, 0x89, 0x8a, 0x8b:
			{
				dis.disaRMftR(op, &opcode, "mov", prevPc)
			}
		case 0x8d:
			{
				opcode.setD(op & 1)
				dis.setMrr(&opcode)
				dumpRMftR(&opcode, prevPc, "lea")
			}
		case 0xb0, 0xb1, 0xb2, 0xb3, 0xb4, 0xb5, 0xb6, 0xb7,
			0xb8, 0xb9, 0xba, 0xbb, 0xbc, 0xbd, 0xbe, 0xbf:
			{

				opcode.setW((op >> 3) & 1)
				opcode.setReg(op & 7)
				dis.setData(&opcode)
				dumpMov(&opcode, prevPc)
			}
		case 0xc3:
			{
				dumpSingleOp(&opcode, prevPc, "ret")
			}
		case 0xcd:
			{

				dis.setData(&opcode)
				dumpInt(&opcode, prevPc)

			}
		case 0xd0, 0xd1, 0xd2, 0xd3:
			{
				nv := dis.lookahead()
				switch (nv >> 3) & 7 {
				case 4:
					{ // shl
						dis.disaLogic(op, &opcode, "shl", prevPc)
					}
				default:
					{
						fmt.Printf("%02x is not implemented yet\n", op)
						os.Exit(1)
					}
				}
			}
		case 0xe4, 0xe5:
			{
				dis.disaInOutPort(op, &opcode, "in", prevPc)
			}
		case 0xe8:
			{
				dis.disa2Jump(op, &opcode, "call", prevPc)
			}
		case 0xe9:
			{
				dis.disa2Jump(op, &opcode, "jmp", prevPc)
			}
		case 0xeb:
			{
				dis.disaJump(op, &opcode, "jmp short", prevPc)
			}
		case 0xec, 0xed:
			{
				dis.disaInOutVar(op, &opcode, "in", prevPc)
			}
		case 0xf4:
			{
				dumpSingleOp(&opcode, prevPc, "hlt")
			}
		case 0xf6, 0xf7:
			{
				nv := dis.lookahead()
				switch (nv >> 3) & 7 {
				case 0:
					{ // cmp
						dis.disaItRM(op, &opcode, "test", prevPc)
					}
				case 3: // neg
					{
						dis.disaOneMrr(op, &opcode, "neg", prevPc)
					}
				default:
					{
						fmt.Println("not implemented for next byte in 0xf6~0xf7")
						os.Exit(1)
					}
				}
			}
		case 0xff:
			{
				nv := dis.lookahead()
				switch (nv >> 3) & 7 {
				case 2: // call
					{
						dis.disaOneMrr(op, &opcode, "call", prevPc)
					}
				case 6: // push
					{
						dis.disaOneMrr(op, &opcode, "push", prevPc)
					}
				default:
					{
						fmt.Println("not implemented for next byte in 0xff")
						os.Exit(1)
					}
				}
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
