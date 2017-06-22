package disasm

import (
	"fmt"
	"os"
)

var reg16 = [...]string{"ax", "cx", "dx", "bx", "sp", "bp", "si", "di"}
var reg8 = [...]string{"al", "cl", "dl", "bl", "ah", "ch", "bh", "dh"}
var eaPrefix = [...]string{"bx+si", "bx+di", "bp+si", "bp+di", "si", "di", "bp", "bx"}

//	reg1 := [8]string{"ax", "cx", "dx", "bx", "sp", "bp", "si", "di"} // こう書くとError 関数の外だかららしい
func formatPrefix(pcs string, rbs string) string {
	var r = pcs + rbs
	var sp string
	for i := 20 - len(r); i > 0; i-- {
		sp += " "
	}
	return r + sp
}

func formatData(w uint8, data uint16) string {
	pat := "%04x"
	if w == 0 {
		pat = "%02x"
		data &= 0xff
	}
	return fmt.Sprintf(pat, data)
}

func dumpAddress(pc uint16) string {
	return fmt.Sprintf("%04x: ", pc)
}

func dumpRawData(opcode *OpCode) string {
	var rbs string
	for i := 0; i < opcode.rawlen; i++ {
		rbs += fmt.Sprintf("%02x", opcode.rawdata[i])
	}
	return rbs
}

func dumpReg(w uint8, r uint8) string {
	reg := reg16[:]
	if w == 0 {
		reg = reg8[:]
	}
	return reg[r]
}

func resolveMrr(w uint8, mod uint8, reg uint8, rm uint8, disp int16) (string, string) {
	regStr := dumpReg(w, reg)
	var eaStr string
	switch mod {
	case 0:
		if rm == 6 {
			eaStr = fmt.Sprintf("%04x", disp)
		} else {
			eaStr = fmt.Sprintf("[%s]", eaPrefix[rm])
		}
	case 1:
		fallthrough
	case 2:
		fallthrough
	case 3:
		fallthrough
	default:
		{
			fmt.Println("not implement or illeagal mod")
			os.Exit(1)
		}
	}

	return regStr, eaStr
}

func format(prefix string, opr string, op1 string, op2 string) string {
	var fmtStr string
	fmtStr = prefix + opr + " " + op1
	if op2 != "" {
		fmtStr += ", "
	}
	fmtStr += op2
	return fmtStr
}

func makePrefix(opcode *OpCode, pc uint16) string {
	return formatPrefix(dumpAddress(pc), dumpRawData(opcode))
}

func dumpMov(opcode *OpCode, pc uint16) {
	/*
		fmt.Println(formatPrefix(dumpAddress(pc), dumpRawData(opcode)))
		fmt.Println(formatData(0, 0x1234))
		dstr := formatPrefix(dumpAddress(pc), dumpRawData(opcode)) +
			"mov " + dumpReg(opcode.W, opcode.Reg) + " " + formatData(opcode.W, opcode.Data)
		fmt.Println(dstr)
		fmt.Println("----")
	*/
	fmt.Println(format(makePrefix(opcode, pc), "mov",
		dumpReg(opcode.W, opcode.Reg), formatData(opcode.W, opcode.Data)))
}

func dumpInt(opcode *OpCode, pc uint16) {
	fmt.Println(format(makePrefix(opcode, pc), "int", formatData(opcode.W, opcode.Data), ""))
}

func dumpAdd(opcode *OpCode, pc uint16) {
	reg, ea := resolveMrr(opcode.W, opcode.Mod, opcode.Reg, opcode.Rm, opcode.Disp)
	if opcode.D == 0 {
		fmt.Println(format(makePrefix(opcode, pc), "add", ea, reg))
	} else {
		fmt.Println(format(makePrefix(opcode, pc), "add", reg, ea))
	}
}
