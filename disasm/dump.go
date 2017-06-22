package disasm

import "fmt"

//var reg16 [2]string = [2]string{"a", "b"}
var reg16 = [...]string{"ax", "cx", "dx", "bx", "sp", "bp", "si", "di"}
var reg8 = [...]string{"al", "cl", "dl", "bl", "ah", "ch", "bh", "dh"}

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

func format(prefix string, opr string, op1 string, op2 string) string {
	var fmtStr string
	fmtStr = prefix + opr + " " + op1
	if op2 != "" {
		fmtStr += " "
	}
	fmtStr += op2
	return fmtStr
}

func dumpMov(opcode *OpCode, pc uint16) {
	fmt.Println("dump")
	/*
		fmt.Println(formatPrefix(dumpAddress(pc), dumpRawData(opcode)))
		fmt.Println(formatData(0, 0x1234))
		dstr := formatPrefix(dumpAddress(pc), dumpRawData(opcode)) +
			"mov " + dumpReg(opcode.W, opcode.Reg) + " " + formatData(opcode.W, opcode.Data)
		fmt.Println(dstr)
		fmt.Println("----")
	*/
	fmt.Println(format(formatPrefix(dumpAddress(pc), dumpRawData(opcode)), "mov",
		dumpReg(opcode.W, opcode.Reg), formatData(opcode.W, opcode.Data)))
}
