package mld

import (
	"bytes"
	"debug/elf"
	"encoding/binary"
	"fmt"
	"os"
	"testing"
)

// llvm-project/llvm/include/llvm/BinaryFormat/ELF.h
var testFile = "init_syscall"

type rewriteSectionTestData struct {
	expectByte  int
	addedSybTab []Elf64_Sym
}

func TestRewriteSection(t *testing.T) {
	testData := rewriteSectionTestData{
		expectByte: 192, // 単純に3つsynbolをタス
		addedSybTab: []Elf64_Sym{
			// __bss_start
			{
				st_name:  3, // TODO!!
				st_info:  0,
				st_other: 0,
				st_shndx: 2,
				st_value: 0x00000000006000e4,
				st_size:  0,
			},
			// _edata
			{
				st_name:  3, // TODO!!
				st_info:  0,
				st_other: 0,
				st_shndx: 2,
				st_value: 0x00000000004000d4,
				st_size:  0,
			},
			// _end
			{
				st_name:  3, // TODO!!
				st_info:  0,
				st_other: 0,
				st_shndx: 2,
				st_value: 0x00000000006000e8,
				st_size:  0,
			},
		},
	}
	// symtab
	// dがどういうdataになってる？
	f, err := elf.Open(testFile)
	if err != nil {
		t.Errorf(err.Error())
	}
	secData, err := f.Section(".symtab").Data()
	if err != nil {
		t.Errorf(err.Error())
	}

	secData2 := rewriteSection(".symtab", secData)
	if len(secData2) != testData.expectByte {
		t.Errorf("expect %d byte,  but got %d bytes\n", testData.expectByte, len(secData2))
	}

	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, secData2)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	dst, err := os.Create("test_tmp")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	_, err = dst.Write(buf.Bytes())
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	fmt.Println("Write done")
}
