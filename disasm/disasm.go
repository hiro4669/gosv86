package disasm

import "fmt"


type Disasm struct {
	pc uint16
	text []byte
}

func (dis Disasm) Run() {

	fmt.Println("start disasm")
}
