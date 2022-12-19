package main

import (
	"fmt"
	"time"
	"os"
	"encoding/binary"
	"encoding/json"
)




const (
	CONSTANT_Class				= 7		// implemented
	CONSTANT_Fieldref			= 9		// implemented
	CONSTANT_Methodref			= 10	// implemented
	CONSTANT_InterfaceMethodref	= 11
	CONSTANT_String				= 8		// implemented
	CONSTANT_Integer			= 3
	CONSTANT_Float				= 4
	CONSTANT_Long				= 5
	CONSTANT_Double				= 6
	CONSTANT_NameAndType		= 12	// implemented
	CONSTANT_Utf8				= 1		// implemented
	CONSTANT_MethodHandle		= 15
	CONSTANT_MethodType			= 16
	CONSTANT_InvokeDynamic		= 18

	ACC_PUBLIC 		 = 0x0001
	ACC_FINAL 		 = 0x0010
	ACC_SUPER 		 = 0x0020
	ACC_INTERFACE   = 0x0200
	ACC_ABSTRACT   = 0x0400
	ACC_SYNTHETIC  = 0x1000
	ACC_ANNOTATION = 0x2000
	ACC_ENUM 	 = 0x4000
)

var AccesFlags = map[string]int{}


type AccessFlag struct {
	Num  int
	Names []string
}

type Class struct {
	Magic        string
	Minor        int
	Major        int
	ConstantPool []ConstantPool
	AccessFlags  AccessFlag
}

type ConstantPool struct {
	Tag      int    `json:"tag,omitempty"`
	TagName  string `json:"tag_name,omitempty"`
	ClassIdx int    `json:"class_idx,omitempty"`

	NameAndTypeIdx int    `json:"name_and_type_idx,omitempty"`
	NameIdx        int    `json:"name_idx,omitempty"`
	DescriptorIdx  int    `json:"descriptor_idx,omitempty"`
	Bytes          string `json:"bytes,omitempty"`
	StringIdx      int    `json:"string_idx,omitempty"`
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
func readLen(file *os.File, length int) []byte {
	data := make([]byte, length)
	file.Read(data)
	return data
}
func readu1i(file *os.File) int {
	data := make([]byte, 1)
	file.Read(data)
	return int(data[0])
}
func readu2i(file *os.File) int {
	data := make([]byte, 2)
	file.Read(data)
	return int(binary.BigEndian.Uint16(data))
}
func readu4i(file *os.File) int {
	data := make([]byte, 4)
	file.Read(data)
	return int(binary.BigEndian.Uint16(data))
}
func readLens(file *os.File, length int) string {
	data := make([]byte, length)
	file.Read(data)
	return string(data)
}
func readu1s(file *os.File) string {
	data := make([]byte, 1)
	file.Read(data)
	return string(data)
}
func readu2s(file *os.File) string {
	data := make([]byte, 2)
	file.Read(data)
	return string(data)
}
func readu4s(file *os.File) string {
	data := make([]byte, 4)
	file.Read(data)
	return string(data)
}



func _init() {
	AccesFlags["ACC_PUBLIC"] 	 = 0x0001
	AccesFlags["ACC_FINAL"] 	 = 0x0010
	AccesFlags["ACC_SUPER"] 	 = 0x0020
	AccesFlags["ACC_INTERFACE"]  = 0x0200
	AccesFlags["ACC_ABSTRACT"]   = 0x0400
	AccesFlags["ACC_SYNTHETIC"]  = 0x1000
	AccesFlags["ACC_ANNOTATION"] = 0x2000
	AccesFlags["ACC_ENUM"] 	 	 = 0x4000
}



func main() {
	start := time.Now()

	_init()
	fmt.Println(" [ Start ] ")

	file, err := os.Open("Main.class")
	if err != nil { fmt.Println(err); return }


	defer func() {
		err := file.Close()
		fmt.Printf(" [ End | took: %dms ] \n", time.Now().UnixMilli() - start.UnixMilli())
		if err != nil { fmt.Println(err); return }
	}()


	class, err := ParseBytecode(file)
	if err != nil { fmt.Println(err) }

	_, err = json.MarshalIndent(class, "", "  ")
	if err != nil { fmt.Println(err); return }
	// fmt.Printf("ConstantPool = %s\n", string(classJson))


	fmt.Println()
	fmt.Println("---------------")
	fmt.Println()

	for _, cp_info := range class.ConstantPool {
		if cp_info.Tag == CONSTANT_Class {
			fmt.Println(class.ConstantPool[cp_info.NameIdx - 1].Bytes)
		}
	}

}

func parseAccessFlag(flags int) AccessFlag {
	var result []string
	for key, val := range AccesFlags {
		if (flags & val) != 0 {
			result = append(result, key)
		}
	}
	return AccessFlag{
		Num: flags,
		Names: result,
	}
}


func ParseBytecode(file *os.File) (Class, error) {

	class := Class{
		Magic: fmt.Sprintf("%x", readu4i(file)),
		Minor: readu2i(file),
		Major: readu2i(file),
	}
	fmt.Printf("class = %v\n", class)

	var constantPool []ConstantPool
	constantPoolCount := readu2i(file)
	for i := 0; i < constantPoolCount-1; i++ {
		tag := readu1i(file)
		cp_info := ConstantPool{ Tag: tag }

		switch cp_info.Tag {
		case CONSTANT_Methodref:
			cp_info.TagName 	   = "CONSTANT_Methodref"
			cp_info.ClassIdx 	   = readu2i(file)
			cp_info.NameAndTypeIdx = readu2i(file)
			break
		case CONSTANT_Class:
			cp_info.TagName = "CONSTANT_Class"
			cp_info.NameIdx = readu2i(file)
			break
		case CONSTANT_NameAndType:
			cp_info.TagName = "CONSTANT_NameAndType"
			cp_info.NameIdx = readu2i(file)
			cp_info.DescriptorIdx = readu2i(file)
			break
		case CONSTANT_Utf8:
			cp_info.TagName = "CONSTANT_Utf8"
			length := readu2i(file)
			cp_info.Bytes = readLens(file, length)
			break
		case CONSTANT_Fieldref:
			cp_info.TagName 	   = "CONSTANT_Fieldref"
			cp_info.ClassIdx 	   = readu2i(file)
			cp_info.NameAndTypeIdx = readu2i(file)
			break
		case CONSTANT_String:
			cp_info.TagName   = "CONSTANT_String"
			cp_info.StringIdx = readu2i(file)
			break

		case 0: break // EOF

		default:
			return Class{}, fmt.Errorf(fmt.Sprintf("Unknown Tag: %d", tag))
		}

		constantPool = append(constantPool, cp_info)
	}

	class.ConstantPool = constantPool


	class.AccessFlags = parseAccessFlag(readu2i(file))
	fmt.Println("access flags = ", class.AccessFlags)


	return class, nil
}
