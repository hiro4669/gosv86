package disasm

import "fmt"

type OpCode struct {
	W       uint8
	D       uint8
	Reg     uint8
	Mod     uint8
	Rm      uint8
	Data    uint16
	Disp    int16
	rawlen  int
	rawdata [20]byte
}

func (op *OpCode) Reset() {
	op.W = 0
	op.D = 0
	op.Mod = 0
	op.Reg = 0
	op.Rm = 0
	op.Data = 0
	op.rawlen = 0
	op.Disp = 0
}

func (op *OpCode) ShowOpCode() {
	fmt.Printf("d = %d, w = %d, reg = %04x\n", op.D, op.W, op.Reg)
	for i := 0; i < op.rawlen; i++ {
		fmt.Printf("%02x ", op.rawdata[i])
	}
	fmt.Println()
}

func (op *OpCode) Add(v byte) byte {
	op.rawdata[op.rawlen] = v
	op.rawlen++
	return v
}

func (op *OpCode) setW(w uint8) {
	op.W = w
}

func (op *OpCode) setD(d uint8) {
	op.D = d
}

func (op *OpCode) setReg(reg uint8) {
	op.Reg = reg
}

func (op *OpCode) setData(data uint16) {
	op.Data = data
}

func (op *OpCode) setMod(mod uint8) {
	op.Mod = mod
}

func (op *OpCode) setRm(rm uint8) {
	op.Rm = rm
}

func (op *OpCode) setDisp(disp int16) {
	op.Disp = disp
}
