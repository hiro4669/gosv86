package disasm

import "fmt"

type Disasm struct {
	pc   uint16
	text []byte
}

func (dis *Disasm) fetch() byte {
	v := dis.text[dis.pc]
	dis.pc++
	return v

}

func (dis *Disasm) Init(text []byte) {
	fmt.Println("init called")
	dis.pc = 0
	dis.text = make([]byte, len(text))
	copy(dis.text, text)
}

func (dis *Disasm) Run() {
	fmt.Println("start disasm")

	for i := 0; i < len(dis.text); i++ {
		fmt.Printf("%02x ", dis.fetch())
	}
	fmt.Printf("\npc = %d\n", dis.pc)

	//	fmt.Println(dis.pc)
	//	fmt.Println(dis.text)
	/*
		for _, v := range dis.text {
			fmt.Printf("%02x ", v)
		}
	*/
}
