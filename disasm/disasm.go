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

func (dis *Disasm) disaTwoMrrW(op byte, opcode *OpCode, opName string, prevPc uint16) {
	opcode.setW(op & 1)
	dis.setMrr(opcode)
	dumpTwoMrr(opcode, prevPc, opName)
}

func (dis *Disasm) disaOneReg(op byte, opcode *OpCode, opName string, prevPc uint16) {
	opcode.setReg(op & 7)
	dumpOneReg(opcode, prevPc, opName)
}
func (dis *Disasm) disaOneRegAc(op byte, opcode *OpCode, opName string, prevPc uint16) {
	opcode.setReg(op & 7)
	dumpOneRegAc(opcode, prevPc, opName)
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

func (dis *Disasm) disaImtoAc(op byte, opcode *OpCode, opName string, prevPc uint16) {
	opcode.setW(op & 1)
	dis.setData(opcode)
	dumpImtoAc(opcode, prevPc, opName)
}

func (dis *Disasm) disaStringMan(op byte, opcode *OpCode, opName string, prevPc uint16) {
	opcode.setW(op & 1)
	dumpStringMan(opcode, prevPc, opName)
}

func (dis *Disasm) Run() {
	var opcode OpCode
	var op byte
	for {
		if int(dis.pc) == len(dis.text) {
			break
		}
		prevPc := dis.pc
	REP:
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
		case 0x10, 0x11, 0x12, 0x13:
			{
				dis.disaRMftR(op, &opcode, "adc", prevPc)
			}
		case 0x18, 0x19, 0x1a, 0x1b:
			{
				dis.disaRMftR(op, &opcode, "sbb", prevPc)
			}
		case 0x20, 0x21, 0x22, 0x23:
			{
				dis.disaRMftR(op, &opcode, "and", prevPc)
			}
		case 0x28, 0x29, 0x2a, 0x2b:
			{
				dis.disaRMftR(op, &opcode, "sub", prevPc)
			}
		case 0x2c, 0x2d:
			{
				dis.disaImtoAc(op, &opcode, "sub", prevPc)
			}
		case 0x30, 0x31, 0x32, 0x33:
			{
				dis.disaRMftR(op, &opcode, "xor", prevPc)
			}
		case 0x38, 0x39, 0x3a, 0x3b:
			{
				dis.disaRMftR(op, &opcode, "cmp", prevPc)
			}
		case 0x3c, 0x3d:
			{
				dis.disaImtoAc(op, &opcode, "cmp", prevPc)
			}
		case 0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47:
			{
				dis.disaOneReg(op, &opcode, "inc", prevPc)
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
		case 0x72:
			{
				dis.disaJump(op, &opcode, "jb", prevPc)
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
		case 0x76:
			{
				dis.disaJump(op, &opcode, "jbe", prevPc)
			}
		case 0x77:
			{
				dis.disaJump(op, &opcode, "jnbe", prevPc)
			}
		case 0x7c:
			{
				dis.disaJump(op, &opcode, "jl", prevPc)
			}
		case 0x7d:
			{
				dis.disaJump(op, &opcode, "jnl", prevPc)
			}
		case 0x7e:
			{
				dis.disaJump(op, &opcode, "jle", prevPc)
			}
		case 0x7f:
			{
				dis.disaJump(op, &opcode, "jnle", prevPc)
			}
		case 0x80, 0x81, 0x82, 0x83:
			{
				nv := dis.lookahead()
				switch (nv >> 3) & 7 {
				case 0:
					{ // add
						dis.disaIfRM(op, &opcode, "add", prevPc)
					}
				case 1: // or
					{
						dis.disaItRM(op, &opcode, "or", prevPc)
					}
				case 3:
					{ // sbb
						dis.disaIfRM(op, &opcode, "sbb", prevPc)
					}
				case 4:
					{ // and
						dis.disaItRM(op, &opcode, "and", prevPc)
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
		case 0x84, 0x85:
			{
				dis.disaRMftR(op, &opcode, "test", prevPc)
			}
		case 0x86, 0x87:
			{
				dis.disaTwoMrrW(op, &opcode, "xchg", prevPc)
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
		case 0x90, 0x91, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97:
			{
				dis.disaOneRegAc(op, &opcode, "xchg", prevPc)
			}
		case 0x98:
			{
				dumpSingleOp(&opcode, prevPc, "cbw")
			}
		case 0x99:
			{
				dumpSingleOp(&opcode, prevPc, "cwd")
			}
		case 0xa4, 0xa5:
			{
				dis.disaStringMan(op, &opcode, "movs", prevPc)
			}
		case 0xa8, 0xa9:
			{
				dis.disaImtoAc(op, &opcode, "test", prevPc)
			}
		case 0xb0, 0xb1, 0xb2, 0xb3, 0xb4, 0xb5, 0xb6, 0xb7,
			0xb8, 0xb9, 0xba, 0xbb, 0xbc, 0xbd, 0xbe, 0xbf:
			{

				opcode.setW((op >> 3) & 1)
				opcode.setReg(op & 7)
				dis.setData(&opcode)
				dumpMov(&opcode, prevPc)
			}
		case 0xc2:
			{
				opcode.setJDisp(dis.fetch2(&opcode))
				dumpJump(&opcode, prevPc, "ret")
			}
		case 0xc3:
			{
				dumpSingleOp(&opcode, prevPc, "ret")
			}
		case 0xc6, 0xc7:
			{
				nv := dis.lookahead()
				switch (nv >> 3) & 7 {
				case 0:
					{
						dis.disaItRM(op, &opcode, "mov", prevPc)
					}
				default:
					{
						fmt.Println("not implemented for next byte in 0xc6~0xc7")
						os.Exit(1)
					}
				}

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
				case 2:
					{ // rcl
						dis.disaLogic(op, &opcode, "rcl", prevPc)
					}
				case 4:
					{ // shl
						dis.disaLogic(op, &opcode, "shl", prevPc)
					}
				case 5:
					{ // shr
						dis.disaLogic(op, &opcode, "shr", prevPc)
					}
				case 7:
					{ // sar
						dis.disaLogic(op, &opcode, "sar", prevPc)
					}
				default:
					{
						fmt.Printf("%02x is not implemented yet\n", op)
						os.Exit(1)
					}
				}
			}
		case 0xe2:
			{
				dis.disaJump(op, &opcode, "loop", prevPc)
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
		case 0xf2, 0xf3:
			{
				opcode.Rep = true
				goto REP
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
						dis.disaOneMrrW(op, &opcode, "neg", prevPc)
					}
				case 4:
					{ // mul
						dis.disaOneMrrW(op, &opcode, "mul", prevPc)
					}
				case 6:
					{ // div
						dis.disaOneMrrW(op, &opcode, "div", prevPc)
					}
				default:
					{
						fmt.Println("not implemented for next byte in 0xf6~0xf7")
						os.Exit(1)
					}
				}
			}
		case 0xfc:
			{
				dumpSingleOp(&opcode, prevPc, "cld")
			}
		case 0xfd:
			{
				dumpSingleOp(&opcode, prevPc, "std")
			}
		case 0xfe, 0xff:
			{
				nv := dis.lookahead()
				switch (nv >> 3) & 7 {
				case 0:
					{ // inc
						dis.disaOneMrrW(op, &opcode, "inc", prevPc)
					}
				case 1:
					{ // dec
						dis.disaOneMrrW(op, &opcode, "dec", prevPc)
					}
				case 2: // call
					{
						dis.disaOneMrrW(op, &opcode, "call", prevPc)
					}
				case 4: // jmp
					{
						dis.disaOneMrrW(op, &opcode, "jmp", prevPc)
					}
				case 6: // push
					{
						dis.disaOneMrrW(op, &opcode, "push", prevPc)
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
