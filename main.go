package main

import (
	"fmt"
	"time"
	"os"
	"encoding/binary"
)



func main() {
	start := time.Now()

	fmt.Println(" [ Start ] ")

	file, err := os.Open("Main.class")
	if err != nil { fmt.Println(err); return }


	defer func() {
		err := file.Close()
		fmt.Printf(" [ End | took: %dms ] \n", time.Now().UnixMilli() - start.UnixMilli())
		if err != nil { fmt.Println(err); return }
	}()


	err = ParseBytecode(file)
	if err != nil { fmt.Println(err) }

}

type Class struct {
	Magic string
	Minor int
	Major int
}


func readu1(file *os.File) []byte {
	data := make([]byte, 1)
	file.Read(data)
	return data
}
func readu2(file *os.File) []byte {
	data := make([]byte, 2)
	file.Read(data)
	return data
}
func readu4(file *os.File) []byte {
	data := make([]byte, 4)
	file.Read(data)
	return data
}


func ParseBytecode(file *os.File) error {

	class := Class{
		Magic: fmt.Sprintf("%x", int(binary.BigEndian.Uint16(readu4(file)))),
		Minor: int(binary.BigEndian.Uint16(readu2(file))),
		Major: int(binary.BigEndian.Uint16(readu2(file))),
	}
	fmt.Printf("class = %v\n", class)

	constantPoolCount := int(binary.BigEndian.Uint16(readu2(file)))
	fmt.Printf("Constant Pool Count = %d\n", constantPoolCount)

	for i := 0; i < constantPoolCount-1; i++ {

	}


	return nil
}
