package main

import (
	"fmt"
	"./disasm"
)

func main() {
	var dis disasm.Disasm
	fmt.Println("Hello Disasm")
	dis.Run()
}

