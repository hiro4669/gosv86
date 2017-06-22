package main

import (
	"fmt"
	"os"

	"./disasm"
)

func fetch2(data []byte, idx int) uint16 {
	return uint16(data[idx+1])<<8 | uint16(data[idx])
}

func partialCopy(src []byte, off int, len int) []byte {
	return append(make([]byte, 0), src[off:off+len]...)
}

func getText(data []byte) []byte {
	var len uint16 = fetch2(data, 8)
	text := partialCopy(data, 0x20, int(len))
	return text
}

func readFile(name string) []byte {
	file, err := os.Open(name)
	if err != nil {
		fmt.Println("cannot open file : " + name)
	}
	defer file.Close()

	var fSize int
	if fi, err := file.Stat(); err == nil {
		if fSize = int(fi.Size()); fSize == 0 {
			fmt.Println("size of the file is zero: exit")
			os.Exit(1)
		}
	}
	var buf []byte = make([]byte, fSize)
	for {
		n, err := file.Read(buf)
		if n == 0 {
			break
		}
		if err != nil {
			break
		}
	}
	return buf

}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("引数の数が違います")
		os.Exit(1)
	}
	fileName := os.Args[1]

	buf := readFile(fileName)
	/*
		for i := 0; i < len(buf); i++ {
			if i%16 == 0 {
				fmt.Println()
			}
			fmt.Printf("%02x ", buf[i])
		}
		fmt.Println()
	*/
	text := getText(buf)
	/*
		for i := 0; i < len(text); i++ {
			if i%16 == 0 {
				fmt.Println()
			}
			fmt.Printf("%02x ", text[i])
		}
		fmt.Println()
	*/
	//	var dis *disasm.Disasm = new(disasm.Disasm)
	var dis disasm.Disasm
	dis.Init(text)
	dis.Run()
}
